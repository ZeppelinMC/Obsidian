package server

import (
	"net"
	"obsidian/log"
	"obsidian/server/auth"
	"obsidian/server/broadcast"
	"obsidian/server/player"
	"obsidian/server/world"
	"obsidian/server/world/generator"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type listing struct {
	HeartbeatURL       string
	HeartbeatFrequency int
	Public             bool
	Enforced           bool
	Enabled            bool
}

type Config struct {
	Address    string
	ServerName string
	ServerMOTD string
	Whitelist  bool
	MaxPlayers int

	TexturePackURL string

	Listing listing
}

func (cfg Config) New() *Server {
	i, err := net.ResolveTCPAddr("tcp", cfg.Address)
	if err != nil {
		panic(err)
	}

	if len(cfg.TexturePackURL) > 64 {
		log.Error("Please get a texture pack url that's 64 characters or shorter. This is due to Minecraft Classic protocol limitations.")
		cfg.TexturePackURL = ""
	}

	l, err := net.ListenTCP("tcp", i)
	if err != nil {
		panic(err)
	}
	log.Info("Loading world")
	w := world.LoadWorld()

	if len(w.Data.BlockArray) != (int(w.Data.X)*int(w.Data.Y)*int(w.Data.Z))-1 {
		log.Infon("Generating world... 0%")
		w.Data.BlockArray = (&generator.DefaultGenerator{}).GenerateWorld(w.Data.X, w.Data.Y, w.Data.Z)
	}

	player.LoadPlayerData()
	return &Server{
		listener:      l,
		players:       broadcast.New[*player.Player](),
		world:         w,
		config:        cfg,
		authenticator: auth.NewAuthenticator("http://www.classicube.net/heartbeat.jsp", cfg.ServerName, 32, i.Port, false),
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
	MaxPlayers: 32,
	Listing: listing{
		HeartbeatURL:       "http://www.classicube.net/heartbeat.jsp",
		Enabled:            true,
		HeartbeatFrequency: 45_000,
	},
}
