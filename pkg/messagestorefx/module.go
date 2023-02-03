package messagestorefx

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-jet/jet/v2/postgres"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/ajvpot/vanderbot/internal/gen/vanderbot/public/table"
)

var Module = fx.Options(fx.Provide(New))

type Params struct {
	fx.In
	Session   *discordgo.Session
	Log       *zap.Logger
	Lifecycle fx.Lifecycle
	Config    config.Provider
	DB        *sql.DB
}

type Result struct {
	fx.Out
	Store Store
}

type Store interface {
	GetMessage(messageID string) (*discordgo.Message, error)
}

type Config struct {
}

const ConfigKey = "messageStore"

type store struct {
	Log    *zap.Logger
	config Config
	db     *sql.DB
}

func New(p Params) (Result, error) {
	s := store{
		Log:    p.Log,
		config: Config{},
		db:     p.DB,
	}

	err := p.Config.Get(ConfigKey).Populate(&s.config)
	if err != nil {
		return Result{}, err
	}

	p.Session.AddHandler(s.handleMessageCreate)
	p.Session.AddHandler(s.handleMessageEdit)
	p.Session.AddHandler(s.handleMessageDelete)

	return Result{Store: &s}, nil
}

func (p *store) handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	p.Log.Info("chat create", zap.Reflect("payload", m))

	p.logMessage(m.Message, false)
}

func (p *store) logMessage(m *discordgo.Message, isDelete bool) {
	serializedMessage, err := json.Marshal(m)
	if err != nil {
		p.Log.Error("error serializing created message", zap.Error(err))
		return
	}
	insertStmt := table.Message.INSERT(table.Message.Blob, table.Message.IsDelete).VALUES(serializedMessage, isDelete)
	_, err = insertStmt.Exec(p.db)
	if err != nil {
		p.Log.Error("error inserting message", zap.Error(err))
		return
	}
}

func (p *store) handleMessageEdit(s *discordgo.Session, m *discordgo.MessageEdit) {
	p.Log.Info("chat edit", zap.Reflect("payload", m))
}

func (p *store) handleMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	p.Log.Info("chat delete", zap.Reflect("payload", m))

	p.logMessage(m.Message, true)
}

func (p *store) GetMessage(messageID string) (*discordgo.Message, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	stmt := table.Message.SELECT(table.Message.Blob).WHERE(table.Message.MessageID.EQ(postgres.String(messageID)))
	res, err := stmt.Rows(ctx, p.db)
	if err != nil {
		p.Log.Error("error retreiving message", zap.Error(err))
		return nil, err
	}

	var msg discordgo.Message
	return &msg, res.Scan(&msg)
}
