package usecase

import (
	
	"fmt"
	"gambling-bot/configs"
	
	
	"sync"

	"math/rand"

	"github.com/bwmarrin/discordgo"
	//"google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

var sessions = make(map[string]*DataSession)
var mutex = sync.Mutex{}

type DataSession struct{
	valueBot	int
	valueUser	int
	Money		int
}



func StartBlackjack(s *discordgo.Session, m *discordgo.MessageCreate){
	mutex.Lock()
	defer mutex.Unlock()

	msgRef := &discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	}
	

	session, exist := sessions[m.Author.ID]

	if !exist {
		session = &DataSession {
			valueBot: rand.Intn(11)+1,
			valueUser: rand.Intn(11)+1,
		}
		sessions[m.Author.ID] = session
	}

	configs.AddUserData(m.Author.Username, 1000000)


	embed := &discordgo.MessageEmbed{
		Title:       "🃏 Blackjack Game",
		Color:       0x00ff00, // Warna hijau
	
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("🧍 Pemain: %s", m.Author.Username),
				Value:  fmt.Sprintf("**Kartu Kamu:** %d\n", session.valueUser),
				Inline: true,
			},
			{
				Name:   "🤖 Bot (Host)",
				Value:  fmt.Sprintf("**Kartu Bot:** %d\n", session.valueBot),
				Inline: true,
			},
		},
	
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Ketik 'hit' untuk ambil kartu lagi, atau 'stand' untuk berhenti.",
		},
	}
		
	s.ChannelMessageSendEmbedReply(m.ChannelID, embed, msgRef)
	
}

func AddCard(s *discordgo.Session, m *discordgo.MessageCreate){
	mutex.Lock()
	defer mutex.Unlock()

	session, exist := sessions[m.Author.ID]
	if !exist {
		s.ChannelMessageSend(m.ChannelID, "Anda belum memulai game! Ketik `blackjack` untuk mulai.")
		return
	}

	msgRef := &discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	}	

	session.valueBot += rand.Intn(11)+1
	session.valueUser += rand.Intn(11)+1

	// s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Kartu %s: %d", m.Author.Username, session.valueUser))
	// s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Kartu Host: %d",session.valueBot))

	embed := &discordgo.MessageEmbed{
		Title:       "🃏 Blackjack - Update Kartu",
		Description: fmt.Sprintf("**%s menarik kartu baru!**", m.Author.Username),
		Color:       0x00ff00, // Warna hijau
	
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("🧍 Pemain: %s", m.Author.Username),
				Value:  fmt.Sprintf("Kartu Baru: %d\n**Total Poin:** %d", rand.Intn(11)+1, session.valueUser),
				Inline: true,
			},
			{
				Name:   "🤖 Bot (Host)",
				Value:  fmt.Sprintf("Kartu Baru: %d\n**Total Poin:** %d", rand.Intn(11)+1, session.valueBot),
				Inline: true,
			},
		},
	
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Ketik 'hit' untuk ambil kartu lagi, atau 'stand' untuk berhenti.",
		},
	}
	
	

	s.ChannelMessageSendEmbedReply(m.ChannelID, embed, msgRef)

	if session.valueBot == 21 && session.valueUser == 21 {
		s.ChannelMessageSendReply(m.ChannelID, "🤝 **DRAWW!** Kedua pihak mendapatkan Blackjack! Pertandingan berakhir seri.", msgRef)
		delete(sessions, m.Author.ID)
		return
	}
	
	if session.valueBot > 21 && session.valueUser > 21 {
		s.ChannelMessageSendReply(m.ChannelID, "🤝 **DRAWW!** Kedua pihak melebihi 21. Pertandingan berakhir tanpa pemenang.", msgRef)
		delete(sessions, m.Author.ID)
		return
	}
	
	if session.valueBot > 21 {
		s.ChannelMessageSendReply(m.ChannelID, "🎉 **Anda Menang!** Bot melebihi 21. Selamat, Anda mendapatkan 50.000!", msgRef)
		delete(sessions, m.Author.ID)
		configs.UpdateUserDataWin(m.Author.Username, 50000)
		return
	}
	
	if session.valueUser > 21 {
		s.ChannelMessageSendReply(m.ChannelID, "💔 **Anda Kalah!** Skor Anda melebihi 21. Cobalah lagi lain kali!", msgRef)
		delete(sessions, m.Author.ID)
		configs.UpdateUserDataLose(m.Author.Username, 50000)
		return
	}
	
	if session.valueUser == 21 {
		s.ChannelMessageSendReply(m.ChannelID, "🎉 **Blackjack! Anda Menang!** Selamat, Anda mencapai skor sempurna 21!", msgRef)
		delete(sessions, m.Author.ID)
		configs.UpdateUserDataWin(m.Author.Username, 50000)
		return
	}
	
	if session.valueBot == 21 {
		s.ChannelMessageSendReply(m.ChannelID, "💔 **Anda Kalah!** Bot mendapatkan Blackjack. Coba keberuntungan Anda lagi!", msgRef)
		delete(sessions, m.Author.ID)
		configs.UpdateUserDataLose(m.Author.Username, 50000)
		return
	}
	


	

	
}

func StayCard(s *discordgo.Session, m *discordgo.MessageCreate){
	mutex.Lock()
	defer mutex.Unlock()

	session, exist := sessions[m.Author.ID]
	if !exist {
		s.ChannelMessageSend(m.ChannelID, "Anda belum memulai game! ketik `blackjack` untuk memulai.")
		return
	}
	
	msgRef := &discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	}

	if session.valueUser < session.valueBot {
		s.ChannelMessageSendReply(m.ChannelID, "Anda Kalah!!", msgRef)
		delete(sessions, m.Author.ID)
		return
	}


	session.valueBot += rand.Intn(11)+1

	embed := &discordgo.MessageEmbed{
		Title:       "🃏 Blackjack - Update Kartu",
		Description: fmt.Sprintf("**%s menarik kartu baru!**", m.Author.Username),
		Color:       0x00ff00, // Warna hijau
	
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("🧍 Pemain: %s", m.Author.Username),
				Value:  fmt.Sprintf("Kartu Baru: %d\n**Total Poin:** %d", rand.Intn(11)+1, session.valueUser),
				Inline: true,
			},
			{
				Name:   "🤖 Bot (Host)",
				Value:  fmt.Sprintf("Kartu Baru: %d\n**Total Poin:** %d", rand.Intn(11)+1, session.valueBot),
				Inline: true,
			},
		},
	
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Ketik 'hit' untuk ambil kartu lagi, atau 'stand' untuk berhenti.",
		},
	}
	
	s.ChannelMessageSendEmbedReply(m.ChannelID, embed, msgRef)


	if session.valueBot == 21 && session.valueUser == 21 {
		s.ChannelMessageSendReply(m.ChannelID, "🤝 **DRAWW!** Kedua pihak mendapatkan Blackjack! Pertandingan berakhir seri.", msgRef)
		delete(sessions, m.Author.ID)
		return
	}
	
	if session.valueBot > 21 && session.valueUser > 21 {
		s.ChannelMessageSendReply(m.ChannelID, "🤝 **DRAWW!** Kedua pihak melebihi 21. Pertandingan berakhir tanpa pemenang.", msgRef)
		delete(sessions, m.Author.ID)
		return
	}
	
	if session.valueBot > 21 {
		s.ChannelMessageSendReply(m.ChannelID, "🎉 **Anda Menang!** Bot melebihi 21. Selamat, Anda mendapatkan 50.000!", msgRef)
		delete(sessions, m.Author.ID)
		configs.UpdateUserDataWin(m.Author.Username, 50000)
		return
	}
	
	if session.valueUser > 21 {
		s.ChannelMessageSendReply(m.ChannelID, "💔 **Anda Kalah!** Skor Anda melebihi 21. Cobalah lagi lain kali!", msgRef)
		delete(sessions, m.Author.ID)
		configs.UpdateUserDataLose(m.Author.Username, 50000)
		return
	}
	
	if session.valueUser == 21 {
		s.ChannelMessageSendReply(m.ChannelID, "🎉 **Blackjack! Anda Menang!** Selamat, Anda mencapai skor sempurna 21!", msgRef)
		delete(sessions, m.Author.ID)
		configs.UpdateUserDataWin(m.Author.Username, 50000)
		return
	}
	
	if session.valueBot == 21 {
		s.ChannelMessageSendReply(m.ChannelID, "💔 **Anda Kalah!** Bot mendapatkan Blackjack. Coba keberuntungan Anda lagi!", msgRef)
		delete(sessions, m.Author.ID)
		configs.UpdateUserDataLose(m.Author.Username, 50000)
		return
	}

	
}

