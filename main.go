package main

import (
	"obsidian/log"
	"obsidian/server"
	"os"
	"os/signal"
	"time"
)

var startTime = time.Now()

func main() {
	log.Info("Starting Obsidian 0.30 Minecraft Server")
	srv := server.New("localhost:25565")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Info("Stopping server")
		srv.Stop()
	}()
	srv.Start(startTime)
}
