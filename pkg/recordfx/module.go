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

	Commands []*discordfx.ApplicationCommandWithHandler `group:"commands,flatten"`
}

type recordingManager struct {
	Session *discordgo.Session
	Log     *zap.Logger
}

func New(p Params) Result {
	vm := recordingManager{
		Session: p.Session,
		Log:     p.Log,
	}

	return Result{Commands: vm.commands()}
}

func (p *recordingManager) commands() []*discordfx.ApplicationCommandWithHandler {
	return []*discordfx.ApplicationCommandWithHandler{{
		Command: discordgo.ApplicationCommand{
			Name:        "record",
			Description: "Record.",
		},
		Handler: p.record,
	}, {
		Command: discordgo.ApplicationCommand{
			Name:        "stop",
			Description: "Stop recording.",
		},
		Handler: p.stop,
	}}
}
