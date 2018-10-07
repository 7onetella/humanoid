package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/nlopes/slack"
)

var channels = map[string]string{}
var groups = map[string]string{}
var users = map[string]string{}
var logger = &MessageLogger{debug: true, color: false}
var botID string
var api *slack.Client
var rtm *slack.RTM
var allowedCommmands = []string{}
var approvalRequired = []string{}
var peers = []string{}
var debug bool

func init() {
	debugEnv := os.Getenv("SLACK_BOT_DEBUG")
	if len(debugEnv) > 0 && "true" == debugEnv {
		debug = true
	}

	token, ok := os.LookupEnv("SLACK_BOT_USER_OAUTH_ACCESS_TOKEN")
	if !ok {
		Println("SLACK_BOT_USER_OAUTH_ACCESS_TOKEN is required")
		os.Exit(1)
	}

	botID, ok = os.LookupEnv("SLACK_BOT_MEMBER_ID")
	if !ok {
		Println("SLACK_BOT_MEMBER_ID is required")
		os.Exit(1)
	}

	api = slack.New(token)
	slack.SetLogger(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags))
	api.SetDebug(false)

	logger.color = colorSupportedTerminal()

	populateChannels(api)

	showChannels()

	populateGroups(api)

	showGroups()

	populateUsers(api)

	showUsers()

	home, err := homedir.Dir()
	if err != nil {
		Println("error while accessing home directory")
	}

	configFile := home + "/.humanoid.config.ini"

	configContent, err := ReadFile(configFile)
	// if this is first run and the user has not set up the config file
	if _, ok := err.(*os.PathError); ok {
		Println("initializing " + configFile)
		err = CreateConfigFile(configFile)
		if err != nil {
			Println(err)
			os.Exit(1)
		}
	}

	allowedCommmands, approvalRequired, peers = ReadConfig(configContent)

}

func main() {

	e := makeExecutionPoint()
	e = makeCheckAllowedCommandMiddleWare()(e)
	e = makeCheckApprovalRequiredCommandMiddleWare()(e)
	e = makeCheckForApprovalKeywordMiddleWare()(e)

	rtm = api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch event := msg.Data.(type) {

		case *slack.MessageEvent:
			if len(event.User) == 0 {
				continue
			}
			channel := event.Msg.Channel

			PrintMessageEvent(api, event)

			if IsMessageNotDirectedAtBot(event) {
				fmt.Println()
				continue
			}

			message := RemoveMention(event)

			req := decodeRequest(message, channel)

			resp, err := e(req)
			if err != nil {
				fmt.Println(err)
				continue
			}

			encodeResponse(resp)

		case *slack.PresenceChangeEvent:
		case *slack.LatencyReport:
		case *slack.HelloEvent:
		case *slack.ConnectedEvent:

		case *slack.RTMError:
			Println(fmt.Sprintf("Error: %s\n", event.Error()))
		case *slack.InvalidAuthEvent:
			Println("Invalid credentials")
			return

		default:
		}
	}
}

// AssertTrue if the condition is not met, logs assert error and exits
func AssertTrue(condition bool, errmsg string) {
	if !condition {
		Println(errmsg)
		os.Exit(1)
	}
}

// Execute parses out command and executes it
func Execute(message string) string {
	cmd, args := GetCommandAndArgs(message)
	output, err := Exec(cmd, args)
	if err != nil {
		Println(err.Error())
	}
	return output
}

// Exec executes commands
func Exec(cmd string, args []string) (string, error) {
	output, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}

// RemoveMention remove mention
func RemoveMention(event *slack.MessageEvent) string {
	textAfterMention := strings.Replace(event.Msg.Text, "<@"+botID+"> ", "", -1)
	return textAfterMention
}

// GetCommandAndArgs returns command and its arguments
func GetCommandAndArgs(textAfterMention string) (string, []string) {
	tokens := strings.Fields(textAfterMention)

	if len(tokens) == 1 {
		return tokens[0], []string{}
	}

	return tokens[0], tokens[1:]
}

// IsMessageNotDirectedAtBot is message directed at bot or not
func IsMessageNotDirectedAtBot(event *slack.MessageEvent) bool {
	// check if we have a DM, or standard channel post
	direct := strings.HasPrefix(event.Msg.Channel, "D")

	// if NOT direct message or NOT mention + message in a channel
	return !direct && !strings.Contains(event.Msg.Text, "@"+botID)
}

// Respond replies back with triple quote
func Respond(rtm *slack.RTM, msg, channel string) {
	if len(msg) > 0 {
		rtm.SendMessage(rtm.NewOutgoingMessage("```"+msg+"```", channel))
	}
}

// ReplaceIDWithMention replaces id with mention name e.g. replace <@FSDFSF> with @morgan
func ReplaceIDWithMention(s string) string {
	ids := []string{}

	for _, f := range strings.Fields(s) {
		var re = regexp.MustCompile(`(?s)\<@(\w+)\>`)

		matches := re.FindStringSubmatch(f)
		if len(matches) == 2 {
			ids = append(ids, matches[1])
		}
	}

	out := s
	for _, id := range ids {
		out = strings.Replace(out, "<@"+id+">", "@"+users[id], -1)
	}

	return out
}

func decodeRequest(message, channel string) BotRequest {
	req := BotRequest{
		message: message,
		channel: channel,
	}
	return req
}

func encodeResponse(resp BotResponse) {

	Respond(rtm, resp.message, resp.channel)

	Println()
}
