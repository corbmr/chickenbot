package main

import (
    "fmt"
    "strings"
    "os"
    "os/signal"
    "flag"
    "time"
    . "github.com/bwmarrin/discordgo"
)

const (
    CommandPrefix = "ch!"
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
    dg, err := New("Bot " + Token)
    if err != nil {
        fmt.Println("error creating Discord session, ", err)
        return
    }

    dg.AddHandler(messageCreate)

    // Open a websocket connection to Discord and begin listening.
    err = dg.Open()
    if err != nil {
    fmt.Println("error opening connection,", err)
        return
    }

    fmt.Println("Chickenbot now running")

    s := make(chan os.Signal, 1)
    signal.Notify(s, os.Interrupt, os.Kill)
    <-s

    dg.Close()

}

func messageCreate(s *Session, m *MessageCreate) {
    // Ignore all messages created by the bot itself
    // This isn't required in this specific example but it's a good practice.
    if m.Author.Bot {
        return
    }

    // auto react to wew
    if strings.Contains(strings.ToLower(m.Content), "wew") {
        go wewReact(s, m.Message)
    }

    command := strings.Fields(m.Content)
    // empty message or not a command
    if len(command) < 1 || command[0] != CommandPrefix {
        return
    }

    // empty command, just bok
    if len(command) < 2 {
        s.ChannelMessageSend(m.ChannelID, "bok bok")
        return
    }

    switch command[1] {
    case "bomb":
        go bomb(s, m.ChannelID, m.Author)

    default:
        s.ChannelMessageSend(m.ChannelID, "Unknown command")
    }


}

func bombHandler(chID, mID string, mch chan<- *MessageReaction) func(*Session, *MessageReactionAdd) {
    return func(s *Session, mr *MessageReactionAdd) {
        if mr.ChannelID == chID && mr.MessageID == mID {
            mch <- mr.MessageReaction
        }
    }
}

const (
    ReactionsNeeded = 3
    SecondsGiven = 10
    RainTime = time.Minute
)

func content(name string, times, seconds int) string {
    return fmt.Sprintf("%s planted a chicken bomb! React %v more times within %v seconds to defuse it!", name, times, seconds)
}

func bomb(s *Session, chID string, u *User) {

    mch := make(chan *MessageReaction, 1)
    r, sec := ReactionsNeeded, SecondsGiven

    m, err := s.ChannelMessageSend(chID, content(u.Username, r, sec))
    if err != nil {
        fmt.Println("Error sending bomb message")
        return
    }

    remove := s.AddHandler(bombHandler(chID, m.ID, mch))
    defer remove()

    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for sec > 0 && r > 0 {
        select {
        case <-ticker.C:
            sec--

        case <-mch:
            r--
        }

        m, err = s.ChannelMessageEdit(chID, m.ID, content(u.Username, r, sec))
        if err != nil {
            fmt.Println("Error editing bomb message")
            return
        }

    }

    if r <= 0 {
        s.ChannelMessageSend(chID, "You have succesfully disarmed the chicken bomb")
    }

    if sec <= 0 {
        go rainChickens(s, chID, u)
    }

}

func rainChickens(s *Session, chID string, u *User) {
    c := fmt.Sprintf("You failed to disarm the chicken bomb. Chickens will rain upon this channel for the next %v", RainTime)
    s.ChannelMessageSend(chID, c)

    remove := s.AddHandler(func(s *Session, m *MessageCreate) {
        if m.ChannelID == chID {
            s.MessageReactionAdd(m.ChannelID, m.ID, "\U0001F414")
        }
    })

    <-time.After(RainTime)

    remove()

    s.ChannelMessageSend(chID, "The chicken rain has receded")
}


func wewReact(s *Session, m *Message) {
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
