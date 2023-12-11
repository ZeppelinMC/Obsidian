package player

import (
	"encoding/json"
	"errors"
	"os"
)

type playerList []string

func (p *playerList) Has(s string) bool {
	for _, p := range *p {
		if p == s {
			return true
		}
	}
	return false
}

func (p *playerList) Add(s string) {
	for _, p := range *p {
		if p == s {
			return
		}
	}
	*p = append(*p, s)
}

func (p *playerList) Remove(s string) {
	for i, pl := range *p {
		if pl == s {
			*p = append((*p)[:i], (*p)[i+1:]...)
			return
		}
	}
}

var (
	Operators     playerList
	Whitelist     playerList
	BannedPlayers playerList
)

var lists = map[string]*playerList{
	"ops.json":            &Operators,
	"whitelist.json":      &Whitelist,
	"banned_players.json": &BannedPlayers,
}

func LoadPlayerData() {
	for f, l := range lists {
		loadList(f, l)
	}
}

func loadList(file string, dst *playerList) {
	f, e := os.ReadFile(file)
	if e != nil {
		if errors.Is(e, os.ErrNotExist) {
			os.WriteFile(file, []byte("[]"), 0755)
		}
		*dst = []string{}
		return
	}
	if json.Unmarshal(f, dst) != nil {
		*dst = []string{}
		return
	}
}
