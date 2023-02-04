package presencefx

import (
	"database/sql"
	"encoding/json"

	"github.com/bwmarrin/discordgo"
	"github.com/go-jet/jet/v2/postgres"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/ajvpot/vanderbot/internal/gen/vanderbot/public/model"
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
	GetPresenceHistory(userID string) ([]*discordgo.Presence, error)
}

type Config struct {
}

const ConfigKey = "presenceStore"

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

	p.Session.AddHandler(s.handlePresenceUpdate)

	return Result{Store: &s}, nil
}

func (p *store) handlePresenceUpdate(s *discordgo.Session, m *discordgo.PresenceUpdate) {
	serializedMessage, err := json.Marshal(m.Presence)
	if err != nil {
		p.Log.Error("error serializing created message", zap.Error(err))
		return
	}
	insertStmt := table.Presence.INSERT(table.Presence.GuildID, table.Presence.Blob).VALUES(m.GuildID, serializedMessage)
	_, err = insertStmt.Exec(p.db)
	if err != nil {
		p.Log.Error("error inserting presence", zap.Error(err))
	}
}

func (p *store) GetPresenceHistory(userID string) ([]*discordgo.Presence, error) {
	stmt := table.Presence.SELECT(table.Presence.Blob).WHERE(table.Presence.UserID.EQ(postgres.String(userID)))
	var dest []struct {
		model.Presence
	}
	err := stmt.Query(p.db, &dest)
	if err != nil {
		p.Log.Error("error reading message", zap.Error(err))
		return nil, err
	}

	presences := make([]*discordgo.Presence, 0, len(dest))
	for i, row := range dest {
		presences[i] = &discordgo.Presence{}
		err = json.Unmarshal(row.Blob, presences[i])
		if err != nil {
			p.Log.Error("error deserializing presence", zap.Error(err))
			return nil, err
		}

	}
	return presences, nil
}
