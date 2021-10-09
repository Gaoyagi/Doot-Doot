package main

import(
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/dgvoice"
	"github.com/spf13/viper"
)

var bound bool = false
var chlBound string
var guildID string
var queue = make([]string, 0, 0)
var voiceCall *discordgo.VoiceConnection

// load the .env file
func loadDotEnv() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
}

func main() {
	loadDotEnv()

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
	split := strings.Split(msg.Content, " ")	//split message up by word
	// if unbound
	if !bound{
		// checks to see if this message will bind it (command: @Doot-Doot bind)
		if msg.Mentions[0].ID==session.State.User.ID{
			if len(split)==2 && split[1]=="bind"{
				bound = true
				chlBound = msg.ChannelID
				fmt.Println("doot-doot is now bound")
			}
		} 
	// if bound
	} else {
		// only consider messages if not from itself and if in bounded channel
		if msg.ChannelID == chlBound && msg.Author.ID != session.State.User.ID {
			killch := make(chan bool)
			switch split[0] {
				// joins vc and then plays specified song (!play songFile)
				case "!play":
					fmt.Println("this is the play command")
					// check if bot is already in call
					if voiceCall == nil {
						vc, err := joinCall(session, msg.Author.ID)	// joins vc
						if err!=nil {
							fmt.Println("Unable to join VC")
						} else {
							voiceCall = vc
						}
					}
					// if the play command doesnt include a url/file name
					if len(split)!=2{
						session.ChannelMessageSend(chlBound, "no file name or url detected")
					} else {
						queue = append(queue, split[1])
						fmt.Println("now playing " + queue[0])
						dgvoice.PlayAudioFile(voiceCall, "./"+queue[0], killch)
						queue = queue[1:]   // removes song that is already playing from queue
					}
				// skips to the next song in queue
				case "!skip":
					fmt.Println("this is the command")
					//killch <- true
					//dgvoice.PlayAudioFile(voiceCall, "./"+queue[0], ch1)
					//queue = queue[1:]
				// stops playing music and leaves the call,
				case "!stop":
					fmt.Println("this is the stop command")
					killch <- true			// kills ffmpeg in playaudiofile, stops the current song playing
					voiceCall.Disconnect()  // leaves the vc
					queue = nil				// clears queue
					// close(killch)
					// killch = make(chan bool)
				default:
					fmt.Println("invalid command")
					session.ChannelMessageSend(chlBound, "invalid command")
			}
		} else {
			return
		}
	}
}

// joins the same voice channel as the user who requested the song
func joinCall(session *discordgo.Session, userID string) (*discordgo.VoiceConnection, error){
	// gets the guild that they belong to
	chnl, err := session.Channel(chlBound)
	if err!=nil {
		fmt.Println("Unable to obtain bound Discord channel")
		return nil, err
	}
	guildID = chnl.GuildID
	// use state.voicestate to get the voice state of the user who called it
	vs, err := session.State.VoiceState(guildID, userID)
	if err!=nil {
		fmt.Println("User is not part of a Voice Call (unable to create VoiceState)")
		return  nil, err
	}
	// use the voice state to get the voice channel the user is in
	// join that voice channel
	return session.ChannelVoiceJoin(guildID, vs.ChannelID, false, false)
}

func downloadSong(session *discordgo.Session, file string, voiceCall*discordgo.VoiceConnection) {
	
	
	/* illegal to use due to dmca and youtube terms of service
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
	}

	stream, _, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		panic(err)
	}

	file, err := os.Create("video.mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	 }*/
}

