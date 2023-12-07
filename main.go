package main

import (
	"bufio"
	"obsidian/log"
	"obsidian/server"
	"obsidian/server/command"
	"obsidian/server/command/core"
	"os"
	"os/signal"
	"strings"
	"time"
)

var startTime = time.Now()

func main() {
	log.Info("Starting Obsidian 0.30 Minecraft Server")
	cfg, _ := server.LoadConfig()
	srv := cfg.New()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Info("Stopping server")
		srv.Stop()
	}()
	go scanConsole(srv)
	srv.Start(startTime)
}

func scanConsole(srv *server.Server) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		cmd := scanner.Text()

		if cmd == "" {
			continue
		}
		args := strings.Split(cmd, " ")
		cmd = args[0]
		args = args[1:]

		c, ok := core.Manager.Search(cmd)
		if !ok {
			log.Print("&cUnknown command. Use \"/help\" for a list of commands.")
			continue
		}
		if c.PlayerOnly {
			log.Print("&cThis command can only be used by players.")
			continue
		}
		c.Execute(command.CommandContext{Arguments: args, Executor: srv, Manager: core.Manager})
	}
}
