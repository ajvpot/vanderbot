package fedfx

import (
	"fmt"
	"strings"

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
	DeletedMessageLogChannel string `yaml:"deletedMessageLogChannel"`
	SpotifyLogChannel        string `yaml:"spotifyLogChannel"`
}
type Config struct {
	Guilds map[string]GuildConfig `yaml:"guilds"`
}
type fedLogger struct {
	Session  *discordgo.Session
	Log      *zap.Logger
	config   Config
	lastSong map[string]string
}

func New(p Params) error {
	pl := fedLogger{
		Session:  p.Session,
		Log:      p.Log,
		config:   Config{},
		lastSong: make(map[string]string),
	}

	err := p.Config.Get("fed").Populate(&pl.config)
	if err != nil {
		return err
	}

	p.Session.AddHandler(pl.handleMessageCreate)
	// todo broken
	//p.Session.AddHandler(pl.handleMessageEdit)
	p.Session.AddHandler(pl.handleMessageDelete)
	p.Session.AddHandler(pl.handleMessageDeleteBulk)
	p.Session.AddHandler(pl.handlePresenceUpdate)
	p.Session.State.MaxMessageCount = 100000

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
	if m.BeforeDelete == nil {
		p.Log.Info("someone deleted a message but we dont have it cached")
		return
	}
	gc, ok := p.config.Guilds[m.BeforeDelete.GuildID]

	if !ok || gc.DeletedMessageLogChannel == "" {
		return
	}

	ch, err := p.Session.State.Channel(m.BeforeDelete.ChannelID)
	if err != nil {
		p.Session.ChannelMessageSend(discordfx.ChannelIDFromString(gc.DeletedMessageLogChannel), fmt.Sprintf("[fed] @%s deleted a message:\n%v", m.BeforeDelete.Author.String(), m.BeforeDelete.Content))
		return
	}
	p.Session.ChannelMessageSend(discordfx.ChannelIDFromString(gc.DeletedMessageLogChannel), fmt.Sprintf("[fed] @%s deleted a message in #%s:\n%v", m.BeforeDelete.Author.String(), ch.Name, m.BeforeDelete.Content))
}

func (p *fedLogger) handlePresenceUpdate(s *discordgo.Session, m *discordgo.PresenceUpdate) {
	gc, ok := p.config.Guilds[m.GuildID]

	if !ok || gc.SpotifyLogChannel == "" {
		return
	}

	songName := spotifySongForPresence(m.Presence)
	if songName == "" {
		p.Log.Debug("no song in spotify presence")
		return
	}

	if lastSong, ok := p.lastSong[m.User.ID]; ok {
		if songName == lastSong {
			p.Log.Debug("duplicate spotify presence")
			return
		}
	}
	p.lastSong[m.User.ID] = songName

	u, err := s.State.Member(m.GuildID, m.User.ID)
	if err != nil {
		p.Log.Error("unknown member", zap.String("gid", m.GuildID))
		return
	}

	s.ChannelMessageSend(discordfx.ChannelIDFromString(gc.SpotifyLogChannel), fmt.Sprintf("[fed] @%s is listening to %s", u.User.String(), songName))
}
func spotifySongForPresence(p discordgo.Presence) string {
	for _, activity := range p.Activities {
		if activity.Name == "Spotify" && strings.HasPrefix(activity.Party.ID, "spotify:") {
			return fmt.Sprintf("%s by %s", activity.Details, activity.State)
		}
	}
	return ""
}
