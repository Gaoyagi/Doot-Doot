package main

import(
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"strings"
)

var bound bool = false
var chlBound string

func main() {
	// load the .env file
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	//create discord session
	session, err := discordgo.New("Bot " + viper.GetString("BOT_TOKEN"))
	if err!=nil {
		fmt.Println("Unable to create Discord session", err)
		return
	}

	//this add handler handles whenever a message is created
	//the message created is a callback
	session.AddHandler(msgCreate)

	// Open a websocket connection to Discord and begin listening.
	err = session.Open()
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
	session.Close()

	
	


}

//checks everytime a message is created
//should look for channel origin, an @, commands
func msgCreate(session *discordgo.Session, msg *discordgo.MessageCreate) {
	split := strings.Split(msg.Content, " ")
	//if unbound
	if !bound{
		// checks to see if this message will bind it (command: @bot bind)
		if msg.Mentions[0].ID==session.State.User.ID{
			if len(split)==2 && split[1]=="bind"{
				bound = true
				chlBound = msg.ChannelID
				fmt.Println("doot-doot is now bound")
			}
		} 
	//if bound
	} else {
		// check to see if fired message was in the bounded to channel
		if msg.ChannelID != chlBound{
			fmt.Println("wrong channel")
			return
		}
		// disregard messages from itself
		if msg.Author.ID == session.State.User.ID {
			return
		} 

		switch split[0] {
			case "!play":
				fmt.Println("this is the play command")
			case "!skip":
				fmt.Println("this is the play command")
			case "!pause":
				fmt.Println("this is the pause command")
			default:
				fmt.Println("invalid command")
		}

	}
}