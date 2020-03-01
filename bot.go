package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"regexp"
	"io/ioutil"
	
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	
	// If the message begings with the string "!rep" reply with a
	// message saying a command was received and send a message
	// notifying each user that was given rep.
	matched, err := regexp.MatchString(`^!rep`, m.Content)

	if err != nil {
		fmt.Println("error processing regexp,", err)
		return
	}


	if matched == true {
		s.ChannelMessageSend(m.ChannelID, "Somebody sent me a command.")

		for _, user := range m.Mentions {
			s.ChannelMessageSend(m.ChannelID, "<@"+user.ID+">"+", you were given rep!")

			// Add rep in database. Create database entry for mentioned user, if none exists.
			
		}
	}
}

// Function for checking existence of string in a file (database).
// Function returns a simple boolean value.
func StringExists(str, filepath string) (bool, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("error reading file,", err)
		return false, err
	}

	stringExists, err := regexp.Match(str, file)
	if err != nil {
		fmt.Println("error matching string", err)
		return false, err
	}
	return stringExists, err
}
