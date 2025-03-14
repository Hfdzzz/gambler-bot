package usecase

import (
	"gambling-bot/configs"

	"github.com/bwmarrin/discordgo"
)

func Leaderboard(s *discordgo.Session, m * discordgo.MessageCreate) {
	configs.GetAllUsers(s,m)
}