package presencelogfx

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(fx.Invoke(New))

type Params struct {
	fx.In
	Session *discordgo.Session
	Log     *zap.Logger
}

type presenceLogger struct {
	Session *discordgo.Session
	Log     *zap.Logger
}

func New(p Params) {
	pl := presenceLogger{
		Session: p.Session,
		Log:     p.Log,
	}

	p.Session.AddHandler(pl.handlePresenceUpdate)

	return
}

func (p *presenceLogger) handlePresenceUpdate(s *discordgo.Session, m *discordgo.PresenceUpdate) {
	g, err := s.State.Guild(m.GuildID)
	if err != nil {
		p.Log.Error("unknown guild", zap.String("gid", m.GuildID))
		return
	}

	ch := p.findPresenceLogChannel(g)
	if ch == nil {
		p.Log.Debug("no music presence channel on guild", zap.String("gid", m.GuildID))
		return
	}

	u, err := s.State.Member(m.GuildID, m.User.ID)
	if err != nil {
		p.Log.Error("unknown member", zap.String("gid", m.GuildID))
		return
	}

	s.ChannelMessageSend(ch.ID, fmt.Sprintf("@%s#%s is listening to %s", u.User.Username, u.User.Discriminator, spotifySongForPresence(m.Presence)))

}

func (p *presenceLogger) findPresenceLogChannel(g *discordgo.Guild) *discordgo.Channel {
	for _, ch := range g.Channels {
		// todo configurable
		if ch.Name == "music-presence" {
			return ch
		}
	}
	return nil
}

func spotifySongForPresence(p discordgo.Presence) string {
	for _, activity := range p.Activities {
		if activity.Name == "Spotify" && strings.HasPrefix(activity.Party.ID, "spotify:") {
			return fmt.Sprintf("%s by %s", activity.Details, activity.State)
		}
	}
	return ""
}
