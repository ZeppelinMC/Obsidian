package net

import (
	"net"
	"obsidian/net/packet"
)

type Conn struct {
	net.Conn
}

func (c Conn) WritePacket(pk packet.Packet) {
	packet.WritePacket(c, pk)
}

func (c Conn) ReadPacket() packet.Packet {
	return packet.ReadPacket(c)
}
