package database

import (
	"context"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/rafli2460/culinary-blog-api/internal/config"
	zero "github.com/rs/zerolog/log"
)

// Prepare prepare sql statements or exit api if fails or error
func Prepare(ctx context.Context, db *sqlx.DB, query string) *sqlx.Stmt {
	s, err := db.PreparexContext(ctx, query)
	if err != nil {
		zero.Error().Stack().
			Str("Context", "Preparing sql statement").
			Str("Query", query).
			Err(err).Msg("")

		os.Exit(config.ExitPrepareStmtFail)
	}
	return s
}

// PrepareNamed prepare sql statements with named bindvars or exit api if fails or error
func PrepareNamed(ctx context.Context, db *sqlx.DB, query string) *sqlx.NamedStmt {
	s, err := db.PrepareNamedContext(ctx, query)
	if err != nil {
		zero.Error().Stack().
			Str("Context", "Preparing sql named statement").
			Str("Query", query).
			Err(err).Msg("")

		os.Exit(config.ExitPrepareStmtFail)
	}
	return s
}
