package main

import (
    "fmt"
    "strings"
    "os"
    "flag"
    dg "github.com/bwmarrin/discordgo"
)

const (
    fcg = "190999489362788353"
    wew = ":wew:477529149964025877"
)

var (
    Token string
)

func init() {
    flag.StringVar(&Token, "t", "", "Bot token")
    flag.Parse()

    if Token == "" {
        flag.Usage()
        os.Exit(1)
    }
}

func main() {
    discord, err := dg.New("Bot " + Token)
    if err != nil {
        fmt.Println("error creating Discord session, ", err)
        return
    }
    defer discord.Close()

    // Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
    }

    fmt.Println("Chickenbot now running")

    <-make(chan struct{})
}

var (
    mid string
)

func messageCreate(s *dg.Session, m *dg.MessageCreate) {
    // Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

    // auto react to wew
    if strings.Contains(m.Content, "wew") {
        wewReact(s, m.Message)
    }

    if !strings.HasPrefix(m.Content, "ch!") {
        return
    }

    content := strings.Fields(m.Content)
    if len(content) < 2 {
        return
    }

    content = content[1:]

    switch content[0] {
    case "hi":
        s.ChannelMessageSend(m.ChannelID, "bok bok")
    }

}

func wewReact(s *dg.Session, m *dg.Message) {
    fmt.Println("detected wew")
    c, err := s.Channel(m.ChannelID)
    if err != nil {
        fmt.Println("error getting channel, ", err)
        return
    }

    // Works only for FCG discord
    if c.GuildID == fcg {
        err = s.MessageReactionAdd(c.ID, m.ID, wew)
        if err != nil {
            fmt.Println("error reacting to wew, ", err)
            return
        }
    }

}
