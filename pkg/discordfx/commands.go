package discordfx

import (
	"context"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type HandlerFunc func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate)

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
	Commands  []*ApplicationCommandWithHandler `group:"commands""`
	Log       *zap.Logger
	Lifecycle fx.Lifecycle
}

type interactionHelper struct {
	i *discordgo.InteractionCreate
}

type InteractionHelper interface {
	GetInteraction() *discordgo.InteractionCreate
	Respond() (<-chan *discordgo.WebhookEdit, error)
	RespondWebhook() (<-chan *discordgo.WebhookEdit, error)
}

func NewInteractionHelper(i *discordgo.InteractionCreate) InteractionHelper {
	return &interactionHelper{i: i}
}

func RegisterCommands(p RegisterCommandsParams) error {
	handlerMap := make(map[string]HandlerFunc)
	registeredCommands := make([]*discordgo.ApplicationCommand, 0, len(p.Commands))

	ctx, cancel := context.WithCancel(context.Background())

	// Set up commands
	for _, commandHandler := range p.Commands {
		handlerMap[commandHandler.Command.Name] = commandHandler.Handler
		p.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ctx, cancel := context.WithTimeout(ctx, time.Second*15)
			defer cancel()

			p.Log.Debug("invoking command handler", zap.String("handlerName", i.ApplicationCommandData().Name))
			if h, ok := handlerMap[i.ApplicationCommandData().Name]; ok {
				h(ctx, s, i)
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
		cancel()
		registeredCommands, err := p.Session.ApplicationCommands(p.Session.State.User.ID, "")
		if err != nil {
			log.Fatalf("Could not fetch registered commands: %v", err)
		}
		for _, v := range registeredCommands {
			err := p.Session.ApplicationCommandDelete(p.Session.State.User.ID, v.GuildID, v.ID)
			if err != nil {
				p.Log.Error("failed to delete command", zap.Error(err))
			}
			p.Log.Debug("deleted command", zap.String("handlerName", v.Name), zap.String("id", v.ID))
		}
		return nil
	}})
	return nil
}
