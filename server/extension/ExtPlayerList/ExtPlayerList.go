package ExtPlayerList

import "obsidian/net/packet"

type ExtAddPlayerName struct {
	NameID                          int16
	PlayerName, ListName, GroupName string
	GroupRank                       byte
}

func (ExtAddPlayerName) ID() byte {
	return 0x16
}

func (e *ExtAddPlayerName) Decode(r packet.Reader) {
	r.Short(&e.NameID)
	r.String(&e.PlayerName)
	r.String(&e.ListName)
	r.String(&e.GroupName)
	r.Byte(&e.GroupRank)
}

func (e ExtAddPlayerName) Encode(w packet.Writer) {
	w.Short(e.NameID)
	w.String(e.PlayerName)
	w.String(e.ListName)
	w.String(e.GroupName)
	w.Byte(e.GroupRank)
}

type ExtAddEntity2 struct {
	EntityID               byte
	InGameName             string
	SkinName               string
	SpawnX, SpawnY, SpawnZ int16
	SpawnYaw, SpawnPitch   byte
}

func (ExtAddEntity2) ID() byte {
	return 0x21
}

func (e *ExtAddEntity2) Decode(r packet.Reader) {
	r.Byte(&e.EntityID)
	r.String(&e.InGameName)
	r.String(&e.SkinName)
	r.Short(&e.SpawnX)
	r.Short(&e.SpawnY)
	r.Short(&e.SpawnZ)
	r.Byte(&e.SpawnYaw)
	r.Byte(&e.SpawnPitch)
}

func (e ExtAddEntity2) Encode(w packet.Writer) {
	w.Byte(e.EntityID)
	w.String(e.InGameName)
	w.String(e.SkinName)
	w.Short(e.SpawnX)
	w.Short(e.SpawnY)
	w.Short(e.SpawnZ)
	w.Byte(e.SpawnYaw)
	w.Byte(e.SpawnPitch)
}

type ExtRemovePlayerName struct {
	NameID int16
}

func (ExtRemovePlayerName) ID() byte {
	return 0x18
}

func (e *ExtRemovePlayerName) Decode(r packet.Reader) {
	r.Short(&e.NameID)
}

func (e ExtRemovePlayerName) Encode(w packet.Writer) {
	w.Short(e.NameID)
}
