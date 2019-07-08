package main

import (
	"strings"

	"github.com/adam-psycho/discordgo"
)

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated user has access to.
func newMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Global Mentions
	Mention := "<@!" + State.Session.User.ID + ">"
	if strings.Contains(m.Message.Content, Mention) {
		go Notify(m.Message)
	}

	// Do nothing when State is disabled
	if !State.Enabled {
		return
	}

	//State Messages
	if m.ChannelID == State.Channel.ID {
		State.AddMessage(m.Message)

		Messages := ReceivingMessageParser(m.Message)

		for _, Msg := range Messages {
			MessagePrint(m.Timestamp, m.Author, State.Guild, Msg)
			//log.Printf("> %s > %s\n", UserName(m.Author.Username), Msg)
		}
	}
}
