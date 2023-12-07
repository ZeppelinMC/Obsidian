package player

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"math"
	"obsidian/atomic"
	"obsidian/net"
	"obsidian/net/packet"
	"obsidian/server/broadcast"
	"obsidian/server/command"
	"obsidian/server/world"
	"slices"
	"strings"
	"sync"
	"unsafe"
)

type counter int8

func (c *counter) Add(i int8) int8 {
	*c++
	return int8(*c)
}

var idCounter counter = -128

type Player struct {
	conn net.Conn
	name string
	id   int8

	Extensions atomic.Value[[]string]
	AppName    atomic.Value[string]

	OP         atomic.Value[bool]
	X, Y, Z    atomic.Value[int16]
	Yaw, Pitch atomic.Value[byte]
	world      *world.World
	players    *broadcast.Broadcaster[*Player]

	commandMgr *command.CommandManager

	mu             sync.RWMutex
	spawnedPlayers []int8
}

func New(name string, conn net.Conn, w *world.World, players *broadcast.Broadcaster[*Player], mgr *command.CommandManager) *Player {
	return &Player{
		name:       name,
		conn:       conn,
		world:      w,
		players:    players,
		id:         idCounter.Add(1),
		commandMgr: mgr,
		OP:         atomic.New(Operators.Has(name)),
	}
}

func (p *Player) Join() {
	p.conn.WritePacket(packet.LevelInitialize{})
	p.sendWorldData()
	p.conn.WritePacket(&packet.LevelFinalize{XSize: p.world.Data.X, YSize: p.world.Data.Y, ZSize: p.world.Data.Z})

	p.conn.WritePacket(&packet.PlayerPositionOrientation{
		PlayerID: -1,
		X:        p.world.Data.Spawn.X,
		Y:        p.world.Data.Spawn.Y,
		Z:        p.world.Data.Spawn.Z,
		Yaw:      byte(p.world.Data.Spawn.H),
		Pitch:    byte(p.world.Data.Spawn.P),
	})
}

func (p *Player) SetBlock(x, y, z int16, blockType byte) {
	p.world.SetBlock(x, y, z, blockType)
	p.players.Range(func(t *Player) bool {
		t.conn.WritePacket(&packet.SetBlock{
			X: x, Y: y, Z: z,
			BlockType: blockType,
		})
		return true
	})
}

func (p *Player) Move(x, y, z int16, yaw, pitch byte) {
	p.X.Set(x)
	p.Y.Set(y)
	p.Z.Set(z)
	p.Yaw.Set(yaw)
	p.Pitch.Set(pitch)

	p.players.Range(func(t *Player) bool {
		if t.Name() == p.Name() {
			return true
		}
		if t.IsSpawned(p) {
			t.conn.WritePacket(&packet.PlayerPositionOrientation{
				PlayerID: int8(p.id),
				X:        x,
				Y:        y,
				Z:        z,
				Yaw:      yaw,
				Pitch:    pitch,
			})
		} else {
			t.SpawnPlayer(p)
		}
		return true
	})
}

func (p *Player) SpawnPlayer(pl *Player) {
	p.mu.Lock()
	p.spawnedPlayers = append(p.spawnedPlayers, pl.id)
	p.mu.Unlock()
	p.conn.WritePacket(&packet.SpawnPlayer{
		PlayerID:   int8(pl.id),
		PlayerName: pl.Name(),
		X:          pl.X.Get(),
		Y:          pl.Y.Get(),
		Z:          pl.Z.Get(),
		Yaw:        pl.Yaw.Get(),
		Pitch:      p.Pitch.Get(),
	})
}

func (p *Player) DespawnPlayer(pl *Player) {
	p.mu.Lock()
	slices.DeleteFunc(p.spawnedPlayers, func(i int8) bool {
		return i == pl.id
	})
	p.mu.Unlock()

	p.conn.WritePacket(&packet.DespawnPlayer{
		PlayerID: pl.id,
	})
}

func (p *Player) IsSpawned(pl *Player) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, s := range p.spawnedPlayers {
		if s == pl.id {
			return true
		}
	}
	return false
}

func (p *Player) sendWorldData() {
	var buf bytes.Buffer
	gun := gzip.NewWriter(&buf)

	l := int32(len(p.world.Data.BlockArray))
	gun.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)})
	gun.Write(*(*[]byte)(unsafe.Pointer(&p.world.Data.BlockArray)))
	gun.Close()

	bytes := buf.Bytes()

	for i := 0; i < len(bytes); i += 1024 {
		x := bytes[i:int(math.Min(float64(i+1024), float64(len(bytes))))]
		complete := byte(0)
		if i != 0 {
			complete = byte(math.Ceil(float64(i) / float64(len(bytes)) * 100))
		}
		p.conn.WritePacket(&packet.LevelDataChunk{
			ChunkData:       x,
			ChunkLength:     int16(len(x)),
			PercentComplete: complete,
		})
	}
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) Disconnect(reason string) {
	p.conn.WritePacket(&packet.DisconnectPlayer{Reason: reason})
	p.conn.Close()
}

func (p *Player) Chat(message string) {
	if strings.HasPrefix(message, "/") {
		if len(message) <= 1 {
			goto chat
		}
		p.command(strings.TrimPrefix(message, "/"))
		return
	}

chat:
	msg := fmt.Sprintf("&f<%s> %s", p.name, message)
	p.players.Range(func(t *Player) bool {
		t.SendMessage(msg)
		return true
	})
}

func (p *Player) command(cmd string) {
	args := strings.Split(cmd, " ")
	cmd = args[0]
	args = args[1:]

	c, ok := p.commandMgr.Search(cmd)
	if !ok {
		p.SendMessage("&cUnknown command. Use \"/help\" for a list of commands.")
		return
	}
	if c.OperatorOnly && !p.OP.Get() {
		p.SendMessage("&cThis command can only be used by operators.")
		return
	}

	c.Execute(command.CommandContext{
		Arguments: args,
		Executor:  p,
		Manager:   p.commandMgr,
	})
}

func (p *Player) SendMessage(msg string) {
	p.conn.WritePacket(&packet.Message{PlayerID: -1, Message: msg})
}
