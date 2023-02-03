package fedfx

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/ajvpot/vanderbot/pkg/discordfx"
	"github.com/ajvpot/vanderbot/pkg/messagestorefx"
)

var Module = fx.Options(fx.Invoke(New))

type Params struct {
	fx.In
	Session      *discordgo.Session
	Log          *zap.Logger
	Config       config.Provider
	MessageStore messagestorefx.Store
}

type DeletedMessageLogConfig struct {
	Channel       discordfx.Channel `yaml:"channel"`
	AllowDeletion bool              `yaml:"allowDeletion"`
}
type GuildConfig struct {
	SpotifyLogChannel discordfx.Channel       `yaml:"spotifyLogChannel"`
	DeletedMessageLog DeletedMessageLogConfig `yaml:"deletedMessageLog"`
}
type Config struct {
	Guilds map[string]GuildConfig `yaml:"guilds"`
}
type fedLogger struct {
	Session      *discordgo.Session
	Log          *zap.Logger
	config       Config
	messageStore messagestorefx.Store
	lastSong     map[string]string
}

func New(p Params) error {
	pl := fedLogger{
		Session:      p.Session,
		Log:          p.Log,
		config:       Config{},
		lastSong:     make(map[string]string),
		messageStore: p.MessageStore,
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
	p.logMessageDelete(m.GuildID, m)
}

func (p *fedLogger) handleMessageDeleteBulk(s *discordgo.Session, m *discordgo.MessageDeleteBulk) {
	for _, mid := range m.Messages {
		msg, err := s.State.Message(m.ChannelID, mid)
		if err != nil {
			p.Log.Error("failed messagedeletebulk parse")
			continue
		}
		p.logMessageDelete(m.GuildID, &discordgo.MessageDelete{
			Message:      nil,
			BeforeDelete: msg,
		})
	}
	return
}

func (p *fedLogger) logMessageDelete(gid string, m *discordgo.MessageDelete) {
	gc, ok := p.config.Guilds[gid]

	if !ok || gc.DeletedMessageLog.Channel == "" {
		return
	}

	if m.BeforeDelete == nil {
		dbm, err := p.messageStore.GetMessage(m.ID)
		if err != nil {
			p.Session.ChannelMessageSend(string(gc.DeletedMessageLog.Channel), fmt.Sprintf("[fed] Someone deleted a message but the contents were not cached in memory: %v", err))
			return
		}
		m.BeforeDelete = dbm
	}

	if gc.DeletedMessageLog.AllowDeletion &&
		string(gc.DeletedMessageLog.Channel) == m.BeforeDelete.ChannelID &&
		m.BeforeDelete.Author.String() == p.Session.State.User.String() {
		p.Log.Info("allowing deletion of deleted message log")
		return
	}

	ch, err := p.Session.State.Channel(m.BeforeDelete.ChannelID)
	if err != nil {
		p.Session.ChannelMessageSend(string(gc.DeletedMessageLog.Channel), fmt.Sprintf("[fed] @%s deleted a message:\n%v", m.BeforeDelete.Author.String(), m.BeforeDelete.Content))
		return
	}
	p.Session.ChannelMessageSend(string(gc.DeletedMessageLog.Channel), fmt.Sprintf("[fed] @%s deleted a message in #%s:\n%v", m.BeforeDelete.Author.String(), ch.Name, m.BeforeDelete.Content))
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

	s.ChannelMessageSend(string(gc.SpotifyLogChannel), fmt.Sprintf("[fed] @%s is listening to %s", u.User.String(), songName))
}
func spotifySongForPresence(p discordgo.Presence) string {
	for _, activity := range p.Activities {
		if activity.Name == "Spotify" && strings.HasPrefix(activity.Party.ID, "spotify:") {
			return fmt.Sprintf("%s by %s", activity.Details, activity.State)
		}
	}
	return ""
}
