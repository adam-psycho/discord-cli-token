package main

import (
	"strconv"
	"strings"
)

//ParseForCommands parses input for Commands, returns message if no command specified, else return is empty
func ParseForCommands(line string) string {
	// Constants
	shrug := "¯\\_(ツ)_/¯"
	tableflip := "(╯°□°）╯︵ ┻━┻"
	unflip := "┬─┬ ノ( ゜-゜ノ)"

	//One Key Commands
	switch line {
	case ":g":
		SelectGuild()
		line = ""
	case ":c":
		SelectChannel()
		line = ""
	case "/shrug":
		line = shrug
	case "/tableflip":
		line = tableflip
	case "/unflip":
		line = unflip
	default:
		// Nothing
	}

	//Argument Commands
	if strings.HasPrefix(line, ":m") {
		AmountStr := strings.Split(line, " ")
		if len(AmountStr) < 2 {
			Msg(ErrorMsg, "[:m] No Arguments \n")
			return ""
		}

		Amount, err := strconv.Atoi(AmountStr[1])
		if err != nil {
			Msg(ErrorMsg, "[:m] Argument Error: %s \n", err)
			return ""
		}

		Msg(InfoMsg, "Printing last %d messages!\n", Amount)
		State.RetrieveMessages(Amount)
		PrintMessages(Amount)
		line = ""
	}

	if strings.HasPrefix(line, "/shrug ") {
		line = line[7:len(line)] + " " + shrug
	}

	if strings.HasPrefix(line, "/tableflip ") {
		line = line[11:len(line)] + " " + tableflip
	}

	if strings.HasPrefix(line, "/unflip ") {
		line = line[8:len(line)] + " " + unflip
	}

	return line
}

//SelectGuild selects a new Guild
func SelectGuild() {
	State.Enabled = false
	SelectGuildMenu()
	SelectChannelMenu()
	State.Enabled = true
	ShowContent()
}

//SelectChannel selects a new Channel
func SelectChannel() {
	State.Enabled = false
	SelectChannelMenu()
	State.Enabled = true
	ShowContent()
}
