package main

import (
	"fmt"

	"github.com/kyokomi/emoji"
	"github.com/nlopes/slack"
)

// MessageLogger logger specialized for slack messages
type MessageLogger struct {
	debug bool
	color bool
}

// SetDebug updates flag for enabling logging
func (ml *MessageLogger) SetDebug(debug bool) {
	ml.debug = debug
}

var logFormat = "%-20s: %s\n"

// Log logs message
func (ml *MessageLogger) Log(msg string) {
	fmt.Printf(logFormat, msg)
}

// Println prints logs
func Println(a ...interface{}) {
	if debug {
		fmt.Println(a...)
	}
}

// PrintMessageEvent prints message info
func PrintMessageEvent(api *slack.Client, event *slack.MessageEvent) {
	text := ReplaceIDWithMention(event.Msg.Text)
	var channel string
	var userName string

	user, err := api.GetUserInfo(event.User)
	if err == nil {
		userName = users[user.ID]
	}

	ch := channels[event.Msg.Channel]

	if len(ch) == 0 {
		channel = emoji.Sprint(":lock:" + groups[event.Msg.Channel])
	} else {
		channel = emoji.Sprint(":hash:" + ch)
	}

	Println(fmt.Sprintf("Request   : %s | %s | %s", text, userName, channel))
	// if debug {
	// 	table := tablewriter.NewWriter(os.Stdout)
	// 	table.SetHeader([]string{"Text", "From", "Channel"})
	// 	table.Append([]string{text, userName, channel})
	// 	table.Render()
	// }
	Println()
}
