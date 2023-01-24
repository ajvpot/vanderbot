package voicefx

import (
	"github.com/bwmarrin/discordgo"
)

type Speaker struct {
	GuildID       string
	ChannelID     string
	SSRC          string
	Username      string
	Discriminator string
}

type voiceHelper struct {
	speakersBySSRC map[string]Speaker
	vc             *discordgo.VoiceConnection
	gid            string
}

// VoiceHelper manages the state for a guild.
// VoiceHelper is responsible for decoding opus audio and splitting it into PCM sample data on a separate channel per speaker.
type VoiceHelper interface {
	GetSpeakers()
	GetGuildID() string
}

// VoiceHelperProvider creates singleton VoiceHelper for each guild if it does not already exist.
type VoiceHelperProvider interface {
	GetVoiceHelper(gid string) VoiceHelper
}
