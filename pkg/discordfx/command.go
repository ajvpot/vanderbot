package discordfx

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type HandlerFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)

// ApplicationCommandWithHandler is a top level Command with a Handler function.
// Subcommands are TODO
// todo: define our own handler func? want to use channels to post back messages.
type ApplicationCommandWithHandler struct {
	Command discordgo.ApplicationCommand
	Handler HandlerFunc
	GuildID string
}

type RegisterCommandsParams struct {
	fx.In
	Session   *discordgo.Session
	Commands  []*ApplicationCommandWithHandler `group:"commands"`
	Log       *zap.Logger
	Lifecycle fx.Lifecycle
}

func RegisterCommands(p RegisterCommandsParams) error {
	handlerMap := make(map[string]HandlerFunc)
	registeredCommands := make([]*discordgo.ApplicationCommand, 0, len(p.Commands))

	// Set up commands
	for _, commandHandler := range p.Commands {
		handlerMap[commandHandler.Command.Name] = commandHandler.Handler
		p.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			p.Log.Debug("invoking command handler", zap.String("handlerName", i.ApplicationCommandData().Name))
			if h, ok := handlerMap[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		})
	}

	p.Lifecycle.Append(fx.Hook{OnStart: func(ctx context.Context) error {
		for _, v := range p.Commands {
			ccmd, err := p.Session.ApplicationCommandCreate(p.Session.State.User.ID, v.GuildID, &v.Command)
			if err != nil {
				return err
			}
			p.Log.Debug("registered command", zap.String("name", ccmd.Name), zap.String("id", ccmd.ID))
			registeredCommands = append(registeredCommands, ccmd)
		}
		return nil
	}, OnStop: func(ctx context.Context) error {
		for _, v := range registeredCommands {
			p.Log.Debug("deleting command handler", zap.String("handlerName", v.Name), zap.String("id", v.ID))
			err := p.Session.ApplicationCommandDelete(p.Session.State.User.ID, v.GuildID, v.ID)
			if err != nil {
				p.Log.Error("failed to delete command", zap.Error(err))
			}
		}
		return nil
	}})
	return nil
}
