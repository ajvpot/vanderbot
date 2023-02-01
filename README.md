# Discord Utility Bot

[![PkgGoDev](https://pkg.go.dev/badge/github.com/ajvpot/vanderbot)](https://pkg.go.dev/github.com/ajvpot/vanderbot)
[![Discord](https://img.shields.io/badge/Discord-%235865F2.svg?style=flat&logo=discord&logoColor=white)](https://discord.com/api/oauth2/authorize?client_id=1066509605556977724&permissions=8&scope=applications.commands%20bot)

## Configuration
* Edit config.yml
* Create `secrets.yml`:
    ```yaml
    discord:
      token: Bot changeme
    ```

## Modules/Usage
### Discord module
* Connects the bot to Discord, provides callbacks for registering event handlers.
* [Configurable](https://pkg.go.dev/github.com/ajvpot/vanderbot/pkg/discordfx#BotConfig)
### "fed" module
* Log Spotify presence messages for a guild to a channel.
* Log deleted messages for a guild to a channel.
* [Configurable](https://pkg.go.dev/github.com/ajvpot/vanderbot/pkg/fedfx#Config)

### ublock
* Check what resources are blocked from loading on a page by uBlock Origin with `/ublock [url]`

### Multichannel Audio Recording
* Join the user's channel with `/join`
* Leave the current voice channel with `/leave`
* Output to separate audio files per speaker
* Option to output to separate audio outputs live for OBS etc.
  * is there some library for ipc audio?

### Call Summary
* Join channels when they become active
* Stopwatch for total call duration
* Stopwatch for individual speakers
* Post summary to text when everyone leaves or when commanded

## TODO
* [x] fedfx
  * [x] logging messages to stdout
  * [x] post to admin channel when deleted message
  * [x] Log spotify presence
* [x] HTTP Server
  * [x] Health check
* [ ] Interaction helper
  * [x] 15s context for interactions
  * [ ] Integrate ublock command
* [ ] Voice helper
  * [ ] Speaker identification
  * [ ] Decoding https://github.com/bwmarrin/dgvoice
  * [ ] Recording https://github.com/bwmarrin/discordgo/blob/master/examples/voice_receive/main.go
  * [ ] Map SSRC to users from speaking updates
  * [ ] JACK Output
    * https://github.com/xthexder/go-jack
    * windows? https://jackaudio.org/faq/jack_on_windows.html
* [ ] Add functionality to discordgo for tracking ssrc (?)
* [ ] Persistent message store in database for other modules to use and support of multiple clients, upgrades without losing messages.
* rm kodata from the repo, add an env var for runtime environment, move config to folder, load config from runtime env file, copy entire config to kodata
* yeet the user's discord token out of browser profile for local operation?