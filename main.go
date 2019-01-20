package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	token = os.Getenv("DISCORD_TOKEN")

	quizletUsername = os.Getenv("QUIZLET_USERNAME")
	quizletPassword = os.Getenv("QUIZLET_PASSWORD")
)

func main() {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	s.State.TrackEmojis = false
	s.State.TrackMembers = false
	s.State.TrackPresences = false
	s.State.TrackRoles = false
	s.State.TrackVoice = false

	s.AddHandler(messageCreate)

	if err := s.Open(); err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}

	select {}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	c, err := s.State.Channel(m.ChannelID)
	if err != nil || c.Type != discordgo.ChannelTypeGuildText {
		return
	}

	name := strings.ToLower(c.Name)

	fn, ok := services[name]
	if !ok {
		if _, ok := serviceInitializers[name]; ok {
			s.ChannelMessageSend(c.ID, "Failed to initialize service.")
		}

		return
	}

	if err := fn(m.Content); err != nil {
		s.ChannelMessageSend(c.ID, fmt.Sprintf("An error has occurred: `%s`", err.Error()))
		return
	}

	s.ChannelMessageSend(c.ID, "Success!")
}
