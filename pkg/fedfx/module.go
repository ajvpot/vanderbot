package fedfx

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/ajvpot/vanderbot/pkg/discordfx"
)

var Module = fx.Options(fx.Invoke(New))

type Params struct {
	fx.In
	Session *discordgo.Session
	Log     *zap.Logger
	Config  config.Provider
}

type GuildConfig struct {
	EnableDeletedMessageLogging bool   `yaml:"enableDeletedMessageLogging"`
	LogChannel                  string `yaml:"logChannel"`
}
type Config struct {
	Guilds map[string]GuildConfig `yaml:"guilds"`
}
type fedLogger struct {
	Session *discordgo.Session
	Log     *zap.Logger
	config  Config
}

func New(p Params) error {
	pl := fedLogger{
		Session: p.Session,
		Log:     p.Log,
		config:  Config{},
	}

	err := p.Config.Get("fed").Populate(&pl.config)
	if err != nil {
		return err
	}

	for gid, gc := range pl.config.Guilds {
		if gc.EnableDeletedMessageLogging && gc.LogChannel == "" {
			return fmt.Errorf("guild has deleted message logging enabled, log channel must be defined: %s", gid)
		}
	}

	p.Session.AddHandler(pl.handleMessageCreate)
	p.Session.AddHandler(pl.handleMessageEdit)
	p.Session.AddHandler(pl.handleMessageDelete)
	p.Session.AddHandler(pl.handleMessageDeleteBulk)

	return nil
}

func (p *fedLogger) handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	p.Log.Info("chat create", zap.Reflect("payload", m))
}

func (p *fedLogger) handleMessageEdit(s *discordgo.Session, m *discordgo.MessageEdit) {
	p.Log.Info("chat edit", zap.Reflect("payload", m))
}

func (p *fedLogger) handleMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	p.logMessageDelete(m)
}

func (p *fedLogger) handleMessageDeleteBulk(s *discordgo.Session, m *discordgo.MessageDeleteBulk) {
	for _, mid := range m.Messages {
		msg, err := s.State.Message(m.ChannelID, mid)
		if err != nil {
			p.Log.Error("failed messagedeletebulk parse")
			continue
		}
		p.logMessageDelete(&discordgo.MessageDelete{
			Message:      nil,
			BeforeDelete: msg,
		})
	}
	return
}

func (p *fedLogger) logMessageDelete(m *discordgo.MessageDelete) {
	gc, ok := p.config.Guilds[m.BeforeDelete.GuildID]

	if !ok || !gc.EnableDeletedMessageLogging {
		return
	}

	p.Session.ChannelMessageSend(discordfx.ChannelIDFromString(gc.LogChannel), fmt.Sprintf("%s deleted a message:\n%v", m.BeforeDelete.Author.String(), m.BeforeDelete))
}
