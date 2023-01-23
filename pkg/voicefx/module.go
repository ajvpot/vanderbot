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

	var conn *discordgo.VoiceConnection

	p.Session.AddHandler(func(s *discordgo.Session, m *discordgo.VoiceServerUpdate) {
		p.Log.Debug("VoiceServerUpdate", zap.Reflect("event", m))
	})
	p.Session.AddHandler(func(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
		p.Log.Debug("VoiceStateUpdate", zap.Reflect("event", m))
		if m.ChannelID != "" {
			var err error
			conn, err = s.ChannelVoiceJoin(m.GuildID, m.ChannelID, true, true)
			if err != nil {
				p.Log.Error("failed to join voice channel", zap.Error(err))
			}
		} else {
			conn.Close()
		}
	})
}
