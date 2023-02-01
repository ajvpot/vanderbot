package main

import (
	"github.com/go-jet/jet/v2/generator/postgres"
	_ "github.com/lib/pq"
)

func generateSql() {
	err := postgres.GenerateDSN(
		"postgres://postgres:postgres@localhost:5499/vanderbot?sslmode=disable",
		"public",
		".",
	)

	if err != nil {
		panic("failed to generate jet generated sql")
	}
}

func main() {
	generateSql()
}
