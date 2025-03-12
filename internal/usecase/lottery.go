package usecase

import (
	"fmt"
	"gambling-bot/configs"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func Spin(s *discordgo.Session, m *discordgo.MessageCreate) {
	value1 := rand.Intn(7) + 1
	value2 := rand.Intn(7) + 1
	value3 := rand.Intn(7) + 1

	spinResult := fmt.Sprintf("🎲 **%d** | 🎲 **%d** | 🎲 **%d**", value1, value2, value3)
	s.ChannelMessageSend(m.ChannelID, spinResult)

	if value1 == value2 && value1 == value3 {
		winMessage := fmt.Sprintf(
			"🎉 **JACKPOT!** 🎉\nSelamat, %s! Anda mendapatkan tiga angka yang sama!\n💰 **Hadiah: 200.000**",
			m.Author.Username,
		)
		s.ChannelMessageSend(m.ChannelID, winMessage)
		configs.UpdateUserDataWin(m.Author.Username, 200000)
		return
	}

	loseMessage := fmt.Sprintf(
		"💔 **Anda Kalah!**\nCoba lagi lain kali, %s! \n❌ Kehilangan: **50.000**",
		m.Author.Username,
	)
	configs.UpdateUserDataLose(m.Author.Username, 50000)
	s.ChannelMessageSend(m.ChannelID, loseMessage)
}
