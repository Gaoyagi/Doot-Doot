# Doot-Doot bot
## About
A Golang Discod muisc bot that join calls and streams music to everyone in the call.
Inspired by the need to replace Rythm and Groovy bots after their takedown and the need to relearn GoLang 
(I also learned why they were taken down lol)

## Direct Dependencies:
- github.com/bwmarrin/dgvoice
- github.com/spf13/viper

## Prerequisite programs:
- Golang
- FFMPEG
- GCC

## How to use:
- clone this repo
- create a .env file
- aquire my bot token (you can't)
- cd to Doot-Doot
- go mod tidy
- go build
- ./Doot-Doot
- Commands:
    - @Doot-Doot bind: binds the bot to a channel and will readd further commands from that channel
    - !play <song file>.mp3: adds song to the play queue and will join the voice call if it isn't in it already
    - !skip: stops the currently playing song and goes to the next song
    - !stop: stops the current song and the bot leaves the voice call

## Future plans/current problems:
- As of right now it only works with 1 server at a time
- needs my bot token to work still
- needs to be published on discord 
- switch the audio stream source from predownloaded audio files to audio live streams


## Author 
George Aoyagi



