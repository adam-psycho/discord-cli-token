package DiscordState

import "github.com/adam-psycho/discordgo"

//State is the current state of the attached client
type State struct {
	Guild         *discordgo.Guild
	Channel       *discordgo.Channel
	Channels      []*discordgo.Channel
	Members       map[string]*discordgo.Member
	Messages      []*discordgo.Message
	Session       *Session
	MessageAmount int  //Amount of Messages to keep in State
	Enabled       bool //Toggles State for Event handling
}

//Session contains the 'state' of the attached server
type Session struct {
	User      *discordgo.User
	Token     string
	DiscordGo *discordgo.Session
	Guilds    []*discordgo.UserGuild
}
