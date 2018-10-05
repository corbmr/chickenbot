package main

import (
    "fmt"
    "strings"
    "github.com/bwmarrin/discordgo"
)


func main() {
    token := "NDk3MDU0NzU0NTc3NTgwMDMz.DplNPA.UFVHE7h-JaikAGd2wsxIrhv-3v8"
    dg, err := discordgo.New("Bot " + token)
    if err != nil {
        fmt.Println("error creating Discord session, ", err)
        return
    }
    defer dg.Close()

    // Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
    }

    fmt.Println("Chickenbot now running")

    <-make(chan struct{})
}

const (
    fcg = "190999489362788353"
    wew = ":wew:477529149964025877"
)

var (
    mid string
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    // Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

    // auto react to wew
    if strings.Contains(m.Content, "wew") {
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
            }
        }
    }

    if !strings.HasPrefix(m.Content, "ch!") {
        return
    }

    content := strings.Split(m.Content, " ")
    if len(content) < 2 {
        return
    }

    content = content[1:]

    switch content[0] {
    case "hi":
        s.ChannelMessageSend(m.ChannelID, "bok bok")
    }

}
