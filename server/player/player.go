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
	"obsidian/server/world"
	a "sync/atomic"
	"unsafe"
)

var idCounter a.Int32

type Player struct {
	conn net.Conn
	name string
	id   int32

	OP         atomic.Value[bool]
	X, Y, Z    atomic.Value[float32]
	Yaw, Pitch atomic.Value[byte]
	world      *world.World
	players    *broadcast.Broadcaster[*Player]
}

func New(name string, conn net.Conn, w *world.World, players *broadcast.Broadcaster[*Player]) *Player {
	return &Player{
		name:    name,
		conn:    conn,
		world:   w,
		players: players,
		id:      idCounter.Add(1),
	}
}

func (p *Player) Join() {
	p.conn.WritePacket(&packet.PlayerPositionOrientation{
		PlayerID: -1,
		X:        float32(p.world.Data.Spawn.X),
		Y:        float32(p.world.Data.Spawn.Y),
		Z:        float32(p.world.Data.Spawn.Z),
		Yaw:      byte(p.world.Data.Spawn.H),
		Pitch:    byte(p.world.Data.Spawn.P),
	})

	p.spawn()

	p.conn.WritePacket(packet.LevelInitialize{})
	p.sendWorldData()
	p.conn.WritePacket(&packet.LevelFinalize{XSize: p.world.Data.X, YSize: p.world.Data.Y, ZSize: p.world.Data.Z})
}

func (p *Player) spawn() {
	p.players.Range(func(t *Player) bool {
		if t.Name() == p.Name() {
			return true
		}

		t.SpawnPlayer(p)
		return true
	})
}

func (p *Player) Move(x, y, z float32, yaw, pitch byte) {
	p.X.Set(x)
	p.Y.Set(y)
	p.Z.Set(z)
	p.Yaw.Set(yaw)
	p.Pitch.Set(pitch)

	p.players.Range(func(t *Player) bool {
		if t.Name() == p.Name() {
			return true
		}

		t.conn.WritePacket(&packet.SpawnPlayer{
			PlayerID:   int8(p.id),
			PlayerName: p.Name(),
			X:          x,
			Y:          y,
			Z:          z,
			Yaw:        yaw,
			Pitch:      pitch,
		})
		return true
	})
}

func (p *Player) SpawnPlayer(pl *Player) {
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
	msg := fmt.Sprintf("&f<%s> %s", p.name, message)
	p.players.Range(func(t *Player) bool {
		t.SendMessage(msg)
		return true
	})
}

func (p *Player) SendMessage(msg string) {
	p.conn.WritePacket(&packet.Message{PlayerID: -1, Message: msg})
}
