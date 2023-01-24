package voicefx

import (
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
		Handler: makeHandleVoiceJoinMe(p),
	}, {
		Command: discordgo.ApplicationCommand{
			Name:        "voiceleave",
			Description: "leave your voice channel",
		},
		Handler: makeHandleVoiceLeave(p),
	}}}
}
