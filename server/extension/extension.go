package extension

import (
	"obsidian/net"
	"obsidian/net/packet"
)

var Extensions = map[string]int32{
	"ExtPlayerList": 2,
	"MessageTypes":  1,
	"FullCP437":     1,
}

func EncodeExtensions(c net.Conn) {
	c.WritePacket(&packet.ExtInfo{
		AppName:        "Obsidian",
		ExtensionCount: int16(len(Extensions)),
	})

	for n, v := range Extensions {
		c.WritePacket(&packet.ExtEntry{
			ExtName: n,
			Version: v,
		})
	}
}

func DecodeExtensions(c net.Conn) (appName string, extensions map[string]int32) {
	i := c.ReadPacket()
	inf, ok := i.(*packet.ExtInfo)
	if !ok {
		return
	}
	appName = inf.AppName

	extensions = make(map[string]int32)

	for i := 0; i < int(inf.ExtensionCount); i++ {
		e := c.ReadPacket()
		ext, ok := e.(*packet.ExtEntry)
		if !ok {
			return
		}
		extensions[ext.ExtName] = ext.Version
	}
	return
}
