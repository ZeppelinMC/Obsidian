package server

import (
	"fmt"
	"net"
	"obsidian/log"
	net2 "obsidian/net"
	"obsidian/net/packet"
	"obsidian/server/auth"
	"obsidian/server/broadcast"
	"obsidian/server/command/core"
	"obsidian/server/extension"
	"obsidian/server/player"
	"obsidian/server/world"
	"time"
)

type Server struct {
	players       *broadcast.Broadcaster[*player.Player]
	config        Config
	world         *world.World
	listener      *net.TCPListener
	authenticator *auth.Authenticator
}

func (srv *Server) Start(startTime time.Time) {
	log.Infof("Done! (%s) Listening for connections on %s", time.Since(startTime), srv.listener.Addr())
	if srv.config.Listing.Enabled {
		s, err := srv.authenticator.Heartbeat(srv.players.Count())
		if err != nil {
			log.Errorf("Heartbeat error: %s", err)
			return
		}
		log.Infof("Server available at %s", s)
		go func() {
			for range time.Tick(time.Millisecond * time.Duration(srv.config.Listing.HeartbeatFrequency)) {
				srv.authenticator.Heartbeat(srv.players.Count())
			}
		}()
	}
	for {
		c, err := srv.listener.Accept()
		if err != nil {
			break
		}

		go srv.handleConnection(c)
	}
}

// Server.Stop saves the world, disconnects each player then stops the server
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
		if srv.config.Listing.Enforced && !srv.authenticator.Validate(pk.VerificationKey, pk.Username) {
			conn.WritePacket(&packet.DisconnectPlayer{Reason: "Failed to authenticate"})
			conn.Close()
			return
		}
		if srv.config.Whitelist && !player.Whitelist.Has(pk.Username) {
			conn.WritePacket(&packet.DisconnectPlayer{Reason: "You are not white-listed in this server"})
			conn.Close()
			return
		}
		if player.BannedPlayers.Has(pk.Username) {
			conn.WritePacket(&packet.DisconnectPlayer{Reason: "You are banned from this server"})
			conn.Close()
			return
		}
		if p := srv.players.Get(pk.Username); p != nil {
			conn.WritePacket(&packet.DisconnectPlayer{Reason: "You are already connected to the server on a different client"})
			conn.Close()
			return
		}
		if srv.players.Count() >= srv.config.MaxPlayers {
			conn.WritePacket(&packet.DisconnectPlayer{Reason: "The server is full"})
			conn.Close()
			return
		}

		p := player.New(pk.Username, conn, srv.world, srv.players, core.Manager)
		srv.players.Set(pk.Username, p)

		if pk.CPE {
			extension.EncodeExtensions(conn)
			app, exts := extension.DecodeExtensions(conn)
			p.AppName.Set(app)
			p.SetExtensions(exts)
		}

		op := byte(0x00)
		if p.OP.Get() {
			op = 0x64
		}
		conn.WritePacket(&packet.ServerIdentification{
			ProtocolVersion: 0x07,
			ServerName:      srv.config.ServerName,
			ServerMOTD:      srv.config.ServerMOTD,
			UserType:        op,
		})

		log.Infof("[%s] Player %s has joined the server using %s", c.RemoteAddr(), p.Name(), p.AppName.Get())

		msg := fmt.Sprintf("&e%s has joined the game", p.Name())

		srv.players.Range(func(t *player.Player) bool {
			t.SendMessage(msg, 0)
			t.AddPlayer(p)

			return true
		})

		p.Join()

		for {
			pac := packet.ReadPacket(c)
			if pac == nil {
				srv.players.Remove(pk.Username)

				log.Infof("[%s] Player %s has left the server", c.RemoteAddr(), p.Name())

				msg := fmt.Sprintf("&e%s has left the game", p.Name())

				srv.players.Range(func(t *player.Player) bool {
					t.SendMessage(msg, 0)
					t.RemovePlayer(p)
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
				p.Move(pk.X, pk.Y, pk.Z, pk.Yaw, pk.Pitch, pk.PlayerID)
			case *packet.SetBlockServer:
				if pk.Mode == 0 {
					pk.BlockType = 0
				}
				p.SetBlock(pk.X, pk.Y, pk.Z, pk.BlockType)
			}
		}
	}
}
