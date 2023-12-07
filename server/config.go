package server

import (
	"net"
	"obsidian/log"
	"obsidian/server/broadcast"
	"obsidian/server/player"
	"obsidian/server/world"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Address    string
	ServerName string
	ServerMOTD string
	Whitelist  bool
}

func (cfg Config) New() *Server {
	i, err := net.ResolveTCPAddr("tcp", cfg.Address)
	if err != nil {
		panic(err)
	}
	l, err := net.ListenTCP("tcp", i)
	if err != nil {
		panic(err)
	}
	log.Info("Loading world")
	w := world.LoadWorld()

	player.LoadPlayerData()
	return &Server{
		listener: l,
		players:  broadcast.New[*player.Player](),
		world:    w,
		config:   cfg,
	}
}

func LoadConfig() (Config, error) {
	var cfg Config
	f, err := os.ReadFile("config.toml")
	if err != nil {
		cfg = Default
		createConfig(cfg)
		return cfg, err
	}
	if err = toml.Unmarshal(f, &cfg); err != nil {
		cfg = Default
		createConfig(cfg)
		return cfg, err
	}

	return cfg, nil
}

func createConfig(cfg Config) {
	f, _ := os.Create("config.toml")

	toml.NewEncoder(f).Encode(cfg)

	f.Close()
}

var Default = Config{
	Address:    "localhost:25565",
	ServerName: "SomeServer",
	ServerMOTD: "This is a Minecraft server powered by Obsidian!",
}
