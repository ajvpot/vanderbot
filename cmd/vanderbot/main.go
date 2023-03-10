package main

import (
	"go.uber.org/fx"

	"github.com/ajvpot/blocksaas/app/mortar/chromefx"

	"github.com/ajvpot/vanderbot/pkg/command/ublockfx"
	"github.com/ajvpot/vanderbot/pkg/configfx"
	"github.com/ajvpot/vanderbot/pkg/dbfx"
	"github.com/ajvpot/vanderbot/pkg/discordfx"
	"github.com/ajvpot/vanderbot/pkg/fedfx"
	"github.com/ajvpot/vanderbot/pkg/store/messagefx"
	"github.com/ajvpot/vanderbot/pkg/store/presencefx"
	"github.com/ajvpot/vanderbot/pkg/systemfx"
	"github.com/ajvpot/vanderbot/pkg/voicefx"
	"github.com/ajvpot/vanderbot/pkg/zapfx"
)

func main() {
	app := fx.New(
		configfx.Module,
		zapfx.Module,
		systemfx.Module,
		dbfx.Module,

		chromefx.Module,

		// command
		ublockfx.Module,
		voicefx.Module,
		fedfx.Module,
		messagefx.Module,
		presencefx.Module,

		// discord
		discordfx.Module,

		/*fx.WithLogger(func() fxevent.Logger {
			return &fxevent.NopLogger
		}),*/
	)

	app.Run()
}
