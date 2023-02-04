package recordfx

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

type recordingManager struct {
	Session               *discordgo.Session
	Log                   *zap.Logger
	recordingStopTriggers map[string]chan struct{}
}

func New(p Params) Result {
	vm := recordingManager{
		Session:               p.Session,
		Log:                   p.Log,
		recordingStopTriggers: make(map[string]chan struct{}),
	}

	return Result{Commands: vm.commands()}
}

func (r *recordingManager) commands() []*discordfx.ApplicationCommandWithHandler {
	return []*discordfx.ApplicationCommandWithHandler{{
		Command: discordgo.ApplicationCommand{
			Name:        "record",
			Description: "Start recording.",
		},
		Handler: r.record,
	}, {
		Command: discordgo.ApplicationCommand{
			Name:        "stop",
			Description: "Stop recording.",
		},
		Handler: r.stop,
	}}
}
