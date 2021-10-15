package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"github.com/Gaoyagi/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

var bound bool = false
var chlBound string
var guildID string
var queue = make(chan string, 20)
var voiceCall *discordgo.VoiceConnection
var killch = make(chan bool, 2)

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
		log.Println("Unable to create Discord session", err)
		return
	}

	//this add handler handles whenever a message is created
	//the message created is a callback
	session.AddHandler(msgCreate)
	// Open a websocket connection to Discord and begin listening.
	err = session.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")

	go playSong()

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
				log.Println("doot-doot is now bound")
			}
		} 
	// if bound
	} else {
		// only consider messages if not from itself and if in bounded channel
		if msg.ChannelID == chlBound && msg.Author.ID != session.State.User.ID {
			switch split[0] {
				// joins vc and then plays specified song (!play songFile)
				case "!play":
					log.Println("this is the play command")
					// check if bot is already in call
					if voiceCall == nil {
						vc, err := joinCall(session, msg.Author.ID)	// joins vc
						if err!=nil {
							log.Println("Unable to join VC")
						} else {
							voiceCall = vc
						}
					} 
					if voiceCall!=nil {
						// if the play command doesnt include a url/file name
						if len(split)!=2{
							session.ChannelMessageSend(chlBound, "no file name or url detected")
						} else {
							queue <- split[1]
						}
					}
					
				// skips to the next song in queue
				case "!skip":
					log.Println("this is skip the command")
					killch <- true
					killch <- false
				// stops playing music and leaves the call,
				case "!stop":
					log.Println("this is the stop command")
					killch <- true			// sends a value to killch, should kill ffmpeg in playaudiofile
					voiceCall.Disconnect()  // leaves the vc
					voiceCall = nil
					queue = nil				// clears queue
					killch <- false
		
				default:
					log.Println("invalid command")
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
		log.Println("Unable to obtain bound Discord channel")
		return nil, err
	}
	guildID = chnl.GuildID
	// use state.voicestate to get the voice state of the user who called it
	vs, err := session.State.VoiceState(guildID, userID)
	if err!=nil {
		log.Println("User is not part of a Voice Call (unable to create VoiceState)")
		return  nil, err
	}
	// use the voice state to get the voice channel the user is in
	// join that voice channel
	return session.ChannelVoiceJoin(guildID, vs.ChannelID, false, false)
}

func playSong() {
	for {
		song:= <-queue
		log.Println("now playing " + song)
		dgvoice.PlayAudioFile(voiceCall, "music/"+song, killch)
	}
	
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
