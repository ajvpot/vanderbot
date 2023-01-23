package discordfx

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const ConfigurationKey = "discord"

type HandlerFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)

// ApplicationCommandWithHandler is a top level Command with a Handler function.
// Subcommands are TODO
// todo: define our own handler func? want to use channels to post back messages.
type ApplicationCommandWithHandler struct {
	Command discordgo.ApplicationCommand
	Handler HandlerFunc
	GuildID string
}

type NewSessionParams struct {
	fx.In
	Commands  []*ApplicationCommandWithHandler `group:"commands"`
	Log       *zap.Logger
	Lifecycle fx.Lifecycle
	Config    config.Provider
}

type NewSessionResult struct {
	fx.Out
	Session *discordgo.Session
}

type BotConfig struct {
	Token string
}

func NewDiscordSession(p NewSessionParams) (NewSessionResult, error) {
	handlerMap := make(map[string]HandlerFunc)
	registeredCommands := make([]*discordgo.ApplicationCommand, 0, len(p.Commands))

	cfg := BotConfig{}

	err := p.Config.Get(ConfigurationKey).Populate(&cfg)
	if err != nil {
		p.Log.Error("failed loading config", zap.Error(err))
		return NewSessionResult{}, err
	}

	s, err := discordgo.New(cfg.Token)
	if err != nil {
		p.Log.Error("invalid bot parameters", zap.Error(err))
		return NewSessionResult{}, err
	}

	instrumentSession(s, p.Log)

	// Set up commands
	for _, commandHandler := range p.Commands {
		handlerMap[commandHandler.Command.Name] = commandHandler.Handler
		s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			p.Log.Debug("invoking command handler", zap.String("handlerName", i.ApplicationCommandData().Name))
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
			p.Log.Debug("registered command", zap.String("name", ccmd.Name), zap.String("id", ccmd.ID))
			registeredCommands = append(registeredCommands, ccmd)
		}
		return nil
	}, OnStop: func(ctx context.Context) error {
		for _, v := range registeredCommands {
			p.Log.Debug("deleting command handler", zap.String("handlerName", v.Name), zap.String("id", v.ID))
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
func instrumentSession(s *discordgo.Session, p *zap.Logger) {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		p.Info("Logged in", zap.String("username", s.State.User.Username), zap.String("discriminator", s.State.User.Discriminator))
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		p.Debug("MessageCreate", zap.Reflect("event", m))
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.PresenceUpdate) {
		p.Debug("PresenceUpdate", zap.Reflect("event", m))
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.VoiceServerUpdate) {
		p.Debug("VoiceServerUpdate", zap.Reflect("event", m))
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
		p.Debug("VoiceStateUpdate", zap.Reflect("event", m))
	})

}

var Module = fx.Options(fx.Provide(NewDiscordSession))
