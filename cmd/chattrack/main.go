package main

import (
	"strings"

	"go.uber.org/config"
	"go.uber.org/fx"

	"github.com/ajvpot/vanderbot/pkg/configfx"
	"github.com/ajvpot/vanderbot/pkg/discordfx"
	"github.com/ajvpot/vanderbot/pkg/recordfx"
	"github.com/ajvpot/vanderbot/pkg/voicefx"
	"github.com/ajvpot/vanderbot/pkg/zapfx"
)

const staticConfig = `
discord:
  token: ${DISCORD_BOT_TOKEN}
  intent: 32767 # IntentsAll


logging:
  level: debug
  development: true
  encoding: console
`

func main() {
	app := fx.New(
		fx.Supply(
			fx.Annotate(config.Source(strings.NewReader(staticConfig)), fx.As(new(config.YAMLOption)), fx.ResultTags(`group:"configopts""`)),
		),

		configfx.Module,
		zapfx.Module,

		// commands
		voicefx.Module,
		recordfx.Module,

		// discord
		discordfx.Module,

		/*fx.WithLogger(func() fxevent.Logger {
			return &fxevent.NopLogger
		}),*/
	)

	app.Run()
}
