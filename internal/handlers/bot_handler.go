package handlers

import (
	"fmt"
	"gambling-bot/configs"
	"gambling-bot/internal/usecase"

	"github.com/bwmarrin/discordgo"
)

func Bot_Handler(s *discordgo.Session, m *discordgo.MessageCreate){
	if m.Author.Bot{
		return
	}

	if m.Content == "hi" {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("HI %s!", m.Author.Username))
	}

	if m.Content == "ld" {
		configs.GetAllUsers(s,m)
	}

	if m.Content == "bj" {
		usecase.StartBlackjack(s, m)
	}else if m.Content == "hit"{
		usecase.AddCard(s, m)
	}else if m.Content == "stand"{
		usecase.StayCard(s, m)
	}

	if m.Content == "mymoney"{
		configs.CheckUserData(s, m, m.Author.Username)
	}

	if m.Content == "spin"{
		usecase.Spin(s, m)
	}

}