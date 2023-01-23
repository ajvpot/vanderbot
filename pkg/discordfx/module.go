package discordfx

import (
	"context"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type HandlerFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)

// ApplicationCommandWithHandler is a top level Command with a Handler function.
// Subcommands are TODO
type ApplicationCommandWithHandler struct {
	Command discordgo.ApplicationCommand
	Handler HandlerFunc
	GuildID string
}

type NewSessionParams struct {
	fx.In
	Commands  []*ApplicationCommandWithHandler
	Client    *http.Client
	Log       zap.Logger
	Lifecycle fx.Lifecycle
	Config    config.Provider
}

type NewSessionResult struct {
	fx.Out
	Session *discordgo.Session
}

func NewDiscordSession(p NewSessionParams) (NewSessionResult, error) {
	handlerMap := make(map[string]HandlerFunc)
	registeredCommands := make([]*discordgo.ApplicationCommand, 0, len(p.Commands))

	s := &discordgo.Session{Client: p.Client}
	instrumentSession(s, p.Log)

	// Set up commands
	for _, commandHandler := range p.Commands {
		handlerMap[commandHandler.Command.Name] = commandHandler.Handler
		s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if h, ok := handlerMap[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		})
	}

	p.Lifecycle.Append(fx.Hook{OnStart: func(ctx context.Context) error {
		err := s.Open()
		if err != nil {
			return err
		}
		for _, v := range p.Commands {
			ccmd, err := s.ApplicationCommandCreate(s.State.User.ID, v.GuildID, &v.Command)
			if err != nil {
				return err
			}
			registeredCommands = append(registeredCommands, ccmd)
		}
		return nil
	}, OnStop: func(ctx context.Context) error {
		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, v.GuildID, v.ID)
			if err != nil {
				p.Log.Error("failed to delete command", zap.Error(err))
				//return err
			}
		}
		return s.Close()
	}})

	return NewSessionResult{Session: s}, nil
}

// instrumentSession adds log handlers for the Session.
func instrumentSession(s *discordgo.Session, p zap.Logger) {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		p.Info("Logged in", zap.String("username", s.State.User.Username), zap.String("discriminator", s.State.User.Discriminator))
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		p.Debug("MessageCreate", zap.Reflect("event", m))
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.PresenceUpdate) {
		p.Debug("PresenceUpdate", zap.Reflect("event", m))
	})
}

var Module = fx.Options(fx.Provide(NewDiscordSession))
