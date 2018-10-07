package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-ini/ini"
	"github.com/nlopes/slack"
	"github.com/olekukonko/tablewriter"
)

func populateChannels(api *slack.Client) {
	chs, err := api.GetChannels(false)
	AssertTrue(err == nil, fmt.Sprintf("error while retriving channels: %v", err))
	for _, channel := range chs {
		channels[channel.ID] = channel.Name
	}
}

func showChannels() {
	Println("Channels")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name"})

	for k, v := range channels {
		table.Append([]string{k, v})
	}
	if debug {
		table.Render()
	}
	Println()
}

func populateGroups(api *slack.Client) {
	grs, err := api.GetGroups(false)
	AssertTrue(err == nil, fmt.Sprintf("error while retriving groups: %v", err))

	for _, group := range grs {
		groups[group.ID] = group.Name
	}
}

func showGroups() {
	Println("Groups")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name"})

	for k, v := range groups {
		table.Append([]string{k, v})
	}
	if debug {
		table.Render()
	}
	Println()
}

func colorSupportedTerminal() bool {
	term, ok := os.LookupEnv("TERM_PROGRAM")

	if ok && (term == "iTerm.app" || term == "xterm") {
		return true
	}

	return false
}

func populateUsers(api *slack.Client) {
	usrs, err := api.GetUsers()
	AssertTrue(err == nil, fmt.Sprintf("error while retriving users: %v", err))

	for _, user := range usrs {
		users[user.ID] = user.Name
	}
}

func showUsers() {
	Println("Users")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name"})

	for k, v := range users {
		table.Append([]string{k, v})
	}
	if debug {
		table.Render()
	}
	Println()
}

// ReadFile reads from file
func ReadFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

// ReadConfig reads ini config
func ReadConfig(data []byte) ([]string, []string, []string, map[string]string) {

	cfg, err := ini.LoadSources(ini.LoadOptions{
		AllowNestedValues:   true,
		UnparseableSections: []string{"allowed", "approval required", "peer approver"},
	}, data)

	if err != nil {
		Println(fmt.Sprintf("Failed to read file: %v", err))
		os.Exit(1)
	}

	allowed := []string{}
	approvalRequired := []string{}
	peers := []string{}
	accounts := map[string]string{}

	for _, section := range cfg.Sections() {
		sectionName := section.Name()

		switch sectionName {
		case "DEFAULT":
			for _, key := range section.Keys() {
				name := key.Name()
				switch name {
				case "auth_command":
					authCommand = key.Value()
				case "session_expired_message":
					sessionExpiredMessage = key.Value()
				}
			}
		case "allowed":
			body := section.Body()
			lines := strings.Split(body, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				allowed = append(allowed, line)
			}
		case "approval required":
			body := section.Body()
			lines := strings.Split(body, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				approvalRequired = append(approvalRequired, line)
			}
		case "peer approver":
			body := section.Body()
			lines := strings.Split(body, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				peers = append(peers, line)
			}
		case "accounts":
			for _, key := range section.Keys() {
				name := key.Name()
				value := key.Value()
				accounts[name] = value
			}
		}

	}

	return allowed, approvalRequired, peers, accounts
}

// CreateConfigFile creates config.ini file
func CreateConfigFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	defer file.Close()

	content := `[allowed]


[approval required]


[peer approver]


`
	_, err = file.Write([]byte(content))

	return err
}
