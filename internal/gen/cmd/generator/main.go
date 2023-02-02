package main

import (
	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	postgres2 "github.com/go-jet/jet/v2/postgres"

	_ "github.com/lib/pq"
)

func generateSql() {
	t := template.Default(postgres2.Dialect).
		UseSchema(func(schema metadata.Schema) template.Schema {
			return template.DefaultSchema(schema).
				UseModel(template.DefaultModel().
					UseTable(func(table metadata.Table) template.TableModel {
						return template.DefaultTableModel(table).
							UseField(func(column metadata.Column) template.TableModelField {
								defaultTableModelField := template.DefaultTableModelField(column)
								if column.DataType.Name == "jsonb" || column.DataType.Name == "json" {
									defaultTableModelField.Type = template.Type{
										ImportPath: "encoding/json",
										Name:       "json.RawMessage",
									}
								}
								return defaultTableModelField
							})
					}),
				)
		})
	err := postgres.GenerateDSN(
		"postgres://postgres:postgres@localhost:5499/vanderbot?sslmode=disable",
		"public",
		".",
		t,
	)

	if err != nil {
		panic("failed to generate jet generated sql")
	}
}

func main() {
	generateSql()
}
