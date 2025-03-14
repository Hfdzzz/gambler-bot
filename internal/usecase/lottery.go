package usecase

import (
	"fmt"
	"gambling-bot/configs"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func Spin(s *discordgo.Session, m *discordgo.MessageCreate) {
	msgRef := &discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	}

	value1 := rand.Intn(7) + 1
	value2 := rand.Intn(7) + 1
	value3 := rand.Intn(7) + 1

	spinResult := fmt.Sprintf("🎲 **%d** | 🎲 **%d** | 🎲 **%d**", value1, value2, value3)
	s.ChannelMessageSendReply(m.ChannelID, spinResult, msgRef)

	if value1 == value2 && value1 == value3 {
		winMessage := fmt.Sprintf(
			"🎉 **JACKPOT!** 🎉\nSelamat, %s! Anda mendapatkan tiga angka yang sama!\n💰 **Hadiah: 200.000**",
			m.Author.Username,
		)
		s.ChannelMessageSendReply(m.ChannelID, winMessage, msgRef)
		configs.UpdateUserDataWin(m.Author.Username, 200000)
		return
	}

	loseMessage := fmt.Sprintf(
		"💔 **Anda Kalah!**\nCoba lagi lain kali, %s! \n❌ Kehilangan: **50.000**",
		m.Author.Username,
	)
	configs.UpdateUserDataLose(m.Author.Username, 50000)
	s.ChannelMessageSendReply(m.ChannelID, loseMessage, msgRef)
}
