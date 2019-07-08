package main

import (
	"encoding/binary"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

//HexColor is a struct gives RGB values
type HexColor struct {
	Color color.Attribute
	R     int
	G     int
	B     int
}

//Msg is a composition of Color.New printf functions
func Msg(MsgType, format string, a ...interface{}) {

	// TODO: Add support for changing color by configuration

	Error := color.New(color.FgRed, color.Bold)
	Info := color.New(color.FgYellow, color.Bold)
	Head := color.New(color.FgCyan, color.Bold)
	Text := color.New(color.FgWhite)

	switch MsgType {
	case "Error":
		Error.Printf(format, a...)
	case "Info":
		Info.Printf(format, a...)
	case "Head":
		Head.Printf(format, a...)
	case "Text":
		Text.Printf(format, a...)
	default:
		Text.Printf(format, a...)
	}
}

//Clear clears the terminal
func Clear() {

	var c *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
	case "linux":
		c = exec.Command("clear")
	case "windows":
		c = exec.Command("cmd", "/c", "cls")
	default:
		Msg(InfoMsg, "Clear function not supported on current OS\n")
	}
	if c != nil {
		c.Stdout = os.Stdout
		c.Run()
	}
}

//Header simply prints a header containing state/session information
func Header() {
	Msg(InfoMsg, "Welcome, %s!\n\n", GetNick(State.Guild, State.Session.User))
	Msg(InfoMsg, "Guild: %s, Channel: %s\n", State.Guild.Name, State.Channel.Name)
}

//ReceivingMessageParser parses receiving message for mentions, images and MultiLine and returns string array
func ReceivingMessageParser(m *discordgo.Message) []string {
	Message, err := m.ContentWithMoreMentionsReplaced(Session.DiscordGo)

	if err != nil {
		Message = m.ContentWithMentionsReplaced()
	}

	//Parse images
	for _, Attachment := range m.Attachments {
		Message = Message + " " + Attachment.URL
	}

	// MultiLine comment parsing
	Messages := strings.Split(Message, "\n")

	return Messages
}

//PrintMessages prints amount of Messages to CLI
func PrintMessages(Amount int) {
	for Key, m := range State.Messages {
		if Key >= len(State.Messages)-Amount {
			Messages := ReceivingMessageParser(m)

			for _, Msg := range Messages {
				//log.Printf("> %s > %s\n", UserName(m.Author.Username), Msg)
				MessagePrint(m.Timestamp, m.Author, State.Guild, Msg)

			}
		}
	}
}

//Notify uses Notify-Send from libnotify to send a notification when a mention arrives.
func Notify(m *discordgo.Message) {
	Channel, err := State.Session.DiscordGo.Channel(m.ChannelID)
	if err != nil {
		Msg(ErrorMsg, "(NOT) Channel Error: %s\n", err)
	}
	Guild, err := State.Session.DiscordGo.Guild(Channel.GuildID)
	if err != nil {
		Msg(ErrorMsg, "(NOT) Guild Error: %s\n", err)
	}
	Title := "@" + m.Author.Username + " : " + Guild.Name + "/" + Channel.Name
	cmd := exec.Command("notify-send", Title, m.ContentWithMentionsReplaced())
	err = cmd.Start()
	if err != nil {
		Msg(ErrorMsg, "(NOT) Check if libnotify is installed, or disable notifications.\n")
	}

}

//MessagePrint prints one correctly formatted Message to stdout
func MessagePrint(Time discordgo.Timestamp, user *discordgo.User, guild *discordgo.Guild, Content string) {
	TimeStamp, _ := time.Parse(time.RFC3339, string(Time))
	LocalTime := TimeStamp.Local().Format("2006/01/02 15:04:05")
	UserName := GetColor(guild, user)

	log.SetFlags(0)
	log.Printf("%s > %s > %s\n", LocalTime, UserName(GetNick(guild, user)), Content)
	log.SetFlags(log.LstdFlags)
}

//GetNick gets the nickname of a user in a guild, or his username if he doesn't have one
func GetNick(guild *discordgo.Guild, user *discordgo.User) string {
	nick := user.Username
	member, err := Session.DiscordGo.State.Member(guild.ID, user.ID)
	if err == nil && member.Nick != "" {
		nick = member.Nick
	}

	return nick
}

//GetColor gets a Sprint for a specified user and specified guild with the first role's color
func GetColor(guild *discordgo.Guild, user *discordgo.User) func(a ...interface{}) string {
	var MemberRoleColor = 0xffffff
	Member, err := Session.DiscordGo.State.Member(guild.ID, user.ID)
	if err == nil && len(Member.Roles) > 0 {
		MemberRole, err := Session.DiscordGo.State.Role(guild.ID, Member.Roles[0])
		if err == nil {
			MemberRoleColor = MemberRole.Color
		}
	}

	Color := ColorMatch(MemberRoleColor)
	return color.New(Color).SprintFunc()
}

//ColorMatch compares HEX->DEC colorcoding and returns the closest ANSI color
func ColorMatch(colorinput int) color.Attribute {
	var Result float64
	var ColorResult color.Attribute
	Result = 10000

	//log.Println(colorinput)

	var ANSIColors []HexColor
	ANSIColors = append(ANSIColors, HexColor{color.FgRed, 255, 0, 0})
	ANSIColors = append(ANSIColors, HexColor{color.FgGreen, 0, 128, 0})
	ANSIColors = append(ANSIColors, HexColor{color.FgYellow, 255, 255, 0})
	ANSIColors = append(ANSIColors, HexColor{color.FgBlue, 0, 0, 255})
	ANSIColors = append(ANSIColors, HexColor{color.FgMagenta, 255, 0, 255})
	ANSIColors = append(ANSIColors, HexColor{color.FgCyan, 0, 255, 255})
	ANSIColors = append(ANSIColors, HexColor{color.FgWhite, 255, 255, 255})
	HexNumber := [4]byte{}
	binary.BigEndian.PutUint32(HexNumber[:], uint32(colorinput))
	InputStruct := HexColor{color.FgBlack, int(HexNumber[1]), int(HexNumber[2]), int(HexNumber[3])}

	for _, acolor := range ANSIColors {
		DiffSum := dis(acolor.R, InputStruct.R) + dis(acolor.G, InputStruct.G) + dis(acolor.B, InputStruct.B)
		TestResult := math.Sqrt(DiffSum)
		if TestResult < Result {
			Result = TestResult
			ColorResult = acolor.Color
		}
	}

	return ColorResult
}

func dis(a, b int) float64 {
	return float64((a - b) * (a - b))
}
