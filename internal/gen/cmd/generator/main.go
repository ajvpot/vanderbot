package main

import (
	"github.com/go-jet/jet/v2/generator/postgres"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func generateSql() {
	err := postgres.GenerateDSN(
		"postgres://postgres:postgres@localhost:5499/vanderbot?sslmode=disable",
		"public",
		".",
	)

	if err != nil {
		log.Error().Err(err).Msg("failed to generate jet generated sql")
		return
	}
}

func main() {
	generateSql()
}
