package discordfx

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const ConfigurationKey = "discord"

type NewSessionParams struct {
	fx.In
	Log       *zap.Logger
	Config    config.Provider
	Lifecycle fx.Lifecycle
}

type NewSessionResult struct {
	fx.Out
	Session *discordgo.Session
}

type BotConfig struct {
	Token string
}

func NewDiscordSession(p NewSessionParams) (NewSessionResult, error) {
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

	p.Lifecycle.Append(fx.Hook{OnStart: func(ctx context.Context) error {
		p.Log.Info("connecting")
		return s.Open()
	}, OnStop: func(ctx context.Context) error {
		p.Log.Info("disconnecting")
		return s.Close()
	}})

	return NewSessionResult{Session: s}, nil
}

// instrumentSession adds log handlers for the Session.
func instrumentSession(s *discordgo.Session, p *zap.Logger) {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		p.Info("connected", zap.String("username", s.State.User.Username), zap.String("discriminator", s.State.User.Discriminator))
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		p.Debug("MessageCreate", zap.Reflect("event", m))
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.PresenceUpdate) {
		p.Debug("PresenceUpdate", zap.Reflect("event", m))
	})
}
