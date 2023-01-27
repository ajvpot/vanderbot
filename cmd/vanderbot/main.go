package main

import (
	"go.uber.org/fx"

	"github.com/ajvpot/blocksaas/app/mortar/chromefx"

	"github.com/ajvpot/vanderbot/pkg/commands/ublockfx"
	"github.com/ajvpot/vanderbot/pkg/configfx"
	"github.com/ajvpot/vanderbot/pkg/discordfx"
	"github.com/ajvpot/vanderbot/pkg/fedfx"
	"github.com/ajvpot/vanderbot/pkg/systemfx"
	"github.com/ajvpot/vanderbot/pkg/voicefx"
	"github.com/ajvpot/vanderbot/pkg/zapfx"
)

func main() {
	app := fx.New(
		configfx.Module,
		zapfx.Module,
		systemfx.Module,

		chromefx.Module,

		// commands
		ublockfx.Module,
		voicefx.Module,
		fedfx.Module,

		// discord
		discordfx.Module,

		/*fx.WithLogger(func() fxevent.Logger {
			return &fxevent.NopLogger
		}),*/
	)

	app.Run()
}
