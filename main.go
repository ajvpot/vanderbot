package main

import (
	"go.uber.org/fx"

	"github.com/ajvpot/vanderbot/pkg/commands/ublockfx"
	"github.com/ajvpot/vanderbot/pkg/configfx"
	"github.com/ajvpot/vanderbot/pkg/discordfx"
	"github.com/ajvpot/vanderbot/pkg/voicefx"
	"github.com/ajvpot/vanderbot/pkg/zapfx"
)

func main() {
	app := fx.New(
		configfx.Module,
		zapfx.Module,

		// commands
		ublockfx.Module,
		voicefx.Module,

		// discord
		discordfx.Module,

		/*fx.WithLogger(func() fxevent.Logger {
			return &fxevent.NopLogger
		}),*/
	)

	app.Run()
}
