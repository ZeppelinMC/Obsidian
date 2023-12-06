package server

import (
	"fmt"
	"net"
	"obsidian/log"
	net2 "obsidian/net"
	"obsidian/net/packet"
	"obsidian/server/broadcast"
	"obsidian/server/command/core"
	"obsidian/server/player"
	"obsidian/server/world"
	"time"
)

type Server struct {
	players  *broadcast.Broadcaster[*player.Player]
	config   Config
	world    *world.World
	listener *net.TCPListener
}

func (srv *Server) Start(startTime time.Time) {
	log.Infof("Done! (%s) Listening for connections on %s", time.Since(startTime), srv.listener.Addr())
	for {
		c, err := srv.listener.Accept()
		if err != nil {
			break
		}

		go srv.handleConnection(c)
	}
}

func (srv *Server) Stop() {
	srv.world.Save()
	srv.players.Range(func(t *player.Player) bool {
		t.Disconnect("Server closed")
		return true
	})
	srv.listener.Close()
}

func (srv *Server) handleConnection(c net.Conn) {
	conn := net2.Conn{Conn: c}
	p := packet.ReadPacket(c)
	if p == nil {
		return
	}
	if pk, ok := p.(*packet.PlayerIdentification); !ok {
		return
	} else {
		if p := srv.players.Get(pk.Username); p != nil {
			conn.WritePacket(&packet.DisconnectPlayer{Reason: "You are already connected to the server on a different client"})
			conn.Close()
		}

		conn.WritePacket(&packet.ServerIdentification{
			ProtocolVersion: 0x07,
			ServerName:      srv.config.ServerName,
			ServerMOTD:      srv.config.ServerMOTD,
			UserType:        0x64,
		})
		p := player.New(pk.Username, conn, srv.world, srv.players, core.Manager)
		srv.players.Set(pk.Username, p)

		msg := fmt.Sprintf("%s has joined the game", p.Name())

		srv.players.Range(func(t *player.Player) bool {
			t.SendMessage(msg)
			return true
		})

		p.Join()

		for {
			pac := packet.ReadPacket(c)
			if pac == nil {
				srv.players.Remove(pk.Username)

				msg := fmt.Sprintf("%s has left the game", p.Name())

				srv.players.Range(func(t *player.Player) bool {
					t.SendMessage(msg)
					if t.IsSpawned(p) {
						t.DespawnPlayer(p)
					}
					return true
				})

				return
			}
			switch pk := pac.(type) {
			case *packet.Message:
				p.Chat(pk.Message)
			case *packet.PlayerPositionOrientation:
				p.Move(pk.X, pk.Y, pk.Z, pk.Yaw, pk.Pitch)
			case *packet.SetBlockServer:
				if pk.Mode == 0 {
					pk.BlockType = 0
				}
				p.SetBlock(pk.X, pk.Y, pk.Z, pk.BlockType)
			}
		}
	}
}
