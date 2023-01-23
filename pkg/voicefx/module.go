package voicefx

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(fx.Invoke(Register))

type Params struct {
	fx.In
	Session *discordgo.Session
	Log     *zap.Logger
}

type Result struct {
	fx.Out
}

func Register(p Params) {
	joined := false
	p.Session.AddHandler(func(s *discordgo.Session, m *discordgo.VoiceServerUpdate) {
		p.Log.Debug("VoiceServerUpdate", zap.Reflect("event", m))
	})
	p.Session.AddHandler(func(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
		p.Log.Debug("VoiceStateUpdate", zap.Reflect("event", m))
		if !joined {
			cv, _ := s.ChannelVoiceJoin(m.GuildID, m.ChannelID, true, false)
			cv.AddHandler(func(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
				p.Log.Debug("voiceSpeakingUpdate", zap.Reflect("conn", vc), zap.Reflect("speaking", vs))
			})
			joined = true
		}
	})
}
