package main

import (
	"go.uber.org/fx"

	"github.com/ajvpot/vanderbot/pkg/commands/ublockfx"
	"github.com/ajvpot/vanderbot/pkg/configfx"
	"github.com/ajvpot/vanderbot/pkg/discordfx"
	"github.com/ajvpot/vanderbot/pkg/zapfx"
)

func main() {
	app := fx.New(
		configfx.Module,
		zapfx.Module,

		// commands
		ublockfx.Module,

		// discord
		discordfx.Module,
	)

	app.Run()
}
