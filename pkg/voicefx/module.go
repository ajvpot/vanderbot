package voicefx

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/ajvpot/vanderbot/pkg/discordfx"
)

var Module = fx.Options(fx.Provide(Register))

type Params struct {
	fx.In
	Session *discordgo.Session
	Log     *zap.Logger
}

type Result struct {
	fx.Out

	Commands []*discordfx.ApplicationCommandWithHandler `group:"commands,flatten"`
}

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

func Register(p Params) Result {
	p.Session.AddHandler(func(s *discordgo.Session, m *discordgo.VoiceServerUpdate) {
		p.Log.Debug("VoiceServerUpdate", zap.Reflect("event", m))
	})
	p.Session.AddHandler(func(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
		p.Log.Debug("VoiceStateUpdate", zap.Reflect("event", m))
	})
	return Result{Commands: []*discordfx.ApplicationCommandWithHandler{{
		Command: discordgo.ApplicationCommand{
			Name:        "voicejoinme",
			Description: "join your voice server",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			g, err := s.State.Guild(i.GuildID)
			if err != nil {
				return
			}

			found := false
			for _, vs := range g.VoiceStates {
				if vs.UserID == i.Member.User.ID {
					cv, _ := s.ChannelVoiceJoin(i.GuildID, vs.ChannelID, true, false)
					cv.AddHandler(func(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
						p.Log.Debug("voiceSpeakingUpdate", zap.Reflect("gid", vc.GuildID), zap.Reflect("cid", vc.ChannelID), zap.Reflect("speaking", vs))
					})
					found = true
				}
			}

			if found {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("ok"),
					},
				})
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("are you in a voice channel?"),
				},
			})
		},
	}, {
		Command: discordgo.ApplicationCommand{
			Name:        "voiceleave",
			Description: "leave your voice channel",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			vc, ok := s.VoiceConnections[i.GuildID]
			if !ok {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "i'm not in your channel",
					},
				})
				return
			}
			vc.Disconnect()
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "ok",
				},
			})
		},
	}}}
}
