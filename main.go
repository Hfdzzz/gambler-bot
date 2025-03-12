package main

import (
	"fmt"
	"gambling-bot/configs"
	"gambling-bot/internal/handlers"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)



func main() {
	configs.LoadEnv()
	configs.ConnectDB()

	var tokenBot = configs.GetEnv("TOKEN_DISCORD")
	dg, err := discordgo.New("Bot " + tokenBot)
	if err != nil {
		log.Fatalf("Failed to connect bot")
		return
	}

	err = dg.Open()
	if err != nil {
		log.Fatalf("Failed to open session bot")
		
		return
	}

	defer dg.Close()
	fmt.Println("Bot is running")

	dg.AddHandler(handlers.Bot_Handler)
	

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

}