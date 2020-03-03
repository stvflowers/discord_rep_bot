package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"regexp"
	"io/ioutil"
	"strconv"
	"strings"
	
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
	
	// Check if message is a command to the bot or not.
	matched, err := regexp.MatchString(`^!rep`, m.Content)

	if err != nil {
		fmt.Println("error processing regexp,", err)
		return
	}


	if matched == true {
		s.ChannelMessageSend(m.ChannelID, "Somebody sent me a command.")

		// Give rep to each user mentioned in the message.
		for _, user := range m.Mentions {
			// s.ChannelMessageSend(m.ChannelID, "<@"+user.ID+">"+", you were given rep!")

			// Declaring variable for database filename
			var database string
			database = `database.txt`
			
			// Add rep in database. Create database entry for mentioned user, if none exists.
			// Check if user exists in database. Note String() function format.
			stringExists, err := StringExists(user.String(), database)

			if err != nil {
				fmt.Println("error checking existence of user in database,", err)
				return
			}

			if stringExists == true {
				// Increment rep of user, in databse, by 1.
				// Get rep value for user from database.
				// Increment rep value for user by 1.
				// Replace old rep value with new rep value for user, in database.

				// Variable for username#discriminator in db.
				dbentry := user.String()
				err := UpdateRep(dbentry, database)
				
				if err != nil {
					fmt.Println("error updating rep,", err)
					return
				}
				
				rep, err := GetUserRep(dbentry, database)

				if err != nil {
					fmt.Println("error getting user rep from database,", err)
					return
				}

				// Notification for user saying that user got rep and how much rep they now have.
				s.ChannelMessageSend(m.ChannelID, "<@"+user.ID+">"+", you now have "+rep+" rep!")
			} else {
				// Create new entry for user, in database, and give 1 rep to user.
				err := AppendStringToFile(user.String()+`=1`, database)

				if err != nil {
					fmt.Println("error creating new entry for user "+user.String()+" in databse,", err)
					return
				}
				}
			}

	}
}

// Function for checking existence of string in a file (database).
// Function returns a simple boolean value and an error value.
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

// Function for appending a string to a file.
func AppendStringToFile(str, filepath string) error {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()
	
	if err != nil {
		fmt.Println("error opening file,", err)
		return err
	}
	
	if _, err = f.WriteString(str); err != nil {
		fmt.Println("error writing to file,", err)
		return err
	}
	
	return err
}

// Function for reading a file line-by-line and doing an
// operation on lines containing a specified string.
func UpdateRep(dbentry, file string) error {
	input, err := ioutil.ReadFile(file)
	
	if err != nil {
		fmt.Println("error reading file,", err)
		return err
	}
	
	lines := strings.Split(string(input), "\n")
	
	for i, line := range lines {
		// Check if line matches dbentry.
		if strings.Contains(line, dbentry) == true {
			// Increment rep by 1.
			// uhd: username#discriminator
			re0 := regexp.MustCompile(`^[^=]+`)
			uhd := re0.FindString(line)
			
			// old rep
			re1 := regexp.MustCompile(`[^=]+$`)
			oldRep := re1.FindString(line)
			
			// convert to int
			oldRepInt, err := strconv.ParseInt(oldRep, 10, 64)
			if err != nil {
				fmt.Println("error converting string to int,", err)
				return err
			}
			
			// Increment rep
			var newRepInt int64
			newRepInt = oldRepInt + 1
			
			// convert back to string
			var newRepStr string
			newRepStr = strconv.FormatInt(newRepInt, 10)
			
			// replace line with updated entry
			lines[i] = uhd+"="+newRepStr
		}
	}
	
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(file, []byte(output), 0644)
	if err != nil {
		fmt.Println("error writing file,", err)
		return err
	}

	return err
}

// Function for getting rep of a specific user
// uhd: username#descriptor
func GetUserRep(uhd, database string) (string, error) {
	// Find line containing string uhd in database
	input, err := ioutil.ReadFile(database)
	
	if err != nil {
		fmt.Println("error reading file,", err)
		return "", err
	}
	
	lines := strings.Split(string(input), "\n")
	
	for _, line := range lines {
		// Check for uhd in database.
			if strings.Contains(line, uhd) == true {
				
				// Assign rep to a variable
				re := regexp.MustCompile(`[^=]+$`)
				rep := re.FindString(line)
				
				return rep, err
			}
	}
	
	return "no dbentry found", err
}
