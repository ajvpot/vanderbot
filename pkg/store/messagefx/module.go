package messagefx

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/go-jet/jet/v2/postgres"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/ajvpot/vanderbot/internal/gen/vanderbot/public/model"
	"github.com/ajvpot/vanderbot/internal/gen/vanderbot/public/table"
)

var Module = fx.Options(fx.Provide(New))

type Params struct {
	fx.In
	Session *discordgo.Session
	Log     *zap.Logger
	DB      *sql.DB
}

type Result struct {
	fx.Out
	Store Store
}

type Store interface {
	GetMessage(messageID string) (*discordgo.Message, error)
}

type store struct {
	Log *zap.Logger
	db  *sql.DB
}

func New(p Params) (Result, error) {
	s := store{
		Log: p.Log,
		db:  p.DB,
	}

	p.Session.AddHandler(s.handleMessageCreate)
	p.Session.AddHandler(s.handleMessageEdit)
	p.Session.AddHandler(s.handleMessageDelete)

	return Result{Store: &s}, nil
}

func (p *store) handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	p.Log.Info("chat create", zap.Reflect("payload", m))

	p.logMessage(m.Message)
}

func (p *store) logMessage(m *discordgo.Message) {
	serializedMessage, err := json.Marshal(m)
	if err != nil {
		p.Log.Error("error serializing created message", zap.Error(err))
		return
	}
	insertStmt := table.Message.INSERT(table.Message.Blob).VALUES(serializedMessage)
	_, err = insertStmt.Exec(p.db)
	if err != nil {
		p.Log.Error("error inserting message", zap.Error(err))
		return
	}
}

func (p *store) handleMessageEdit(s *discordgo.Session, m *discordgo.MessageUpdate) {
	p.Log.Info("chat edit", zap.Reflect("payload", m))

	p.logMessage(m.Message)
}

func (p *store) handleMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	p.Log.Info("chat delete", zap.Reflect("payload", m))

	p.logMessage(m.Message)
}

func (p *store) GetMessage(messageID string) (*discordgo.Message, error) {
	stmt := table.Message.SELECT(table.Message.Blob).WHERE(table.Message.MessageID.EQ(postgres.String(messageID)))
	var dest []struct {
		model.Message
	}
	err := stmt.Query(p.db, &dest)
	if err != nil {
		p.Log.Error("error reading message", zap.Error(err))
		return nil, err
	}

	for _, row := range dest {
		var msg discordgo.Message
		return &msg, json.Unmarshal(row.Blob, &msg)
	}
	return nil, errors.New("no message found")
}
