package fedfx

import (
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

type fedLogger struct {
	Session *discordgo.Session
	Log     *zap.Logger
}

func New(p Params) {
	pl := fedLogger{
		Session: p.Session,
		Log:     p.Log,
	}

	p.Session.AddHandler(pl.handleMessageCreate)
	p.Session.AddHandler(pl.handleMessageEdit)
	p.Session.AddHandler(pl.handleMessageDelete)
	p.Session.AddHandler(pl.handleMessageDeleteBulk)

	return
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
	for _, mid := m.Messages{
		msg, err := s.State.Message(m.ChannelID, mid)
		if err != nil{
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

}
