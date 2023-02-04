package voicefx

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/ajvpot/vanderbot/pkg/discordfx"
)

var Module = fx.Options(fx.Provide(New))

type Params struct {
	fx.In
	Session *discordgo.Session
	Log     *zap.Logger
}

type Result struct {
	fx.Out

	Commands []*discordfx.ApplicationCommandWithHandler `group:"command,flatten"`
}

type voiceManagerState struct {
	Recording bool
}
type voiceManager struct {
	Session *discordgo.Session
	Log     *zap.Logger
	state   map[string]voiceManagerState
}

func New(p Params) Result {
	p.Session.AddHandler(func(s *discordgo.Session, m *discordgo.VoiceServerUpdate) {
		p.Log.Debug("VoiceServerUpdate", zap.Reflect("event", m))
	})
	p.Session.AddHandler(func(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
		p.Log.Debug("VoiceStateUpdate", zap.Reflect("event", m))
	})

	vm := voiceManager{
		Session: p.Session,
		Log:     p.Log,
	}

	return Result{Commands: vm.commands()}
}

func (p *voiceManager) commands() []*discordfx.ApplicationCommandWithHandler {
	return []*discordfx.ApplicationCommandWithHandler{{
		Command: discordgo.ApplicationCommand{
			Name:        "join",
			Description: "Join your voice channel.",
		},
		Handler: p.joinme,
	}, {
		Command: discordgo.ApplicationCommand{
			Name:        "leave",
			Description: "Leave your voice channel.",
		},
		Handler: p.leaveme,
	}}
}
