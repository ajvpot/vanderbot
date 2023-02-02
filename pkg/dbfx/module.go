package dbfx

import (
	"database/sql"

	_ "github.com/lib/pq"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(fx.Provide(New))

type Params struct {
	fx.In
	Log       *zap.Logger
	Lifecycle fx.Lifecycle
	Config    config.Provider
}

type Result struct {
	fx.Out
	DB *sql.DB
}

type Config struct {
	URL string `yaml:"url"`
}

const ConfigKey = "db"

func New(p Params) (Result, error) {
	var cfg Config
	err := p.Config.Get(ConfigKey).Populate(&cfg)
	if err != nil {
		return Result{}, err
	}
	db, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return Result{}, err
	}
	return Result{DB: db}, nil
}
