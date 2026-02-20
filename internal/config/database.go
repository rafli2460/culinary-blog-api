package config

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/rafli2460/culinary-blog-api/pkg/logger"
	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Database struct {
	Read  *sqlx.DB
	Write *sqlx.DB
}

func InitDB() *Database {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dbDialect := os.Getenv("DB_DIALECT")
	writeHost := os.Getenv("DB_WRITE_HOST")
	readHost := os.Getenv("DB_READ_HOST")

	writeDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, writeHost, dbPort, dbName)
	readDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, readHost, dbPort, dbName)

	writeDB, err := sqlx.Connect(dbDialect, writeDSN)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Can't connect to Writer DB")
	}
	log.Info().Msg("Successfully connect to Writer DB")

	readDB, err := sqlx.Connect(dbDialect, readDSN)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Can't connect to Reader DB")
	}
	log.Info().Msg("Successfully connect to Reader DB")

	runMigrations(writeDB)

	return &Database{
		Read:  readDB,
		Write: writeDB,
	}
}

func (db *Database) Close() {
	if db.Read != nil {
		db.Read.Close()
	}

	if db.Write != nil {
		db.Write.Close()
	}

	log.Info().Msg("Connection has been closed")
}

func runMigrations(db *sqlx.DB) {
	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		logger.SystemError("error creating database driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		logger.SystemError("error initiating database migration")
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Info().Msg("Database migration: no changes to apply")
		} else {
			logger.LogError(err, "failed to run migration")
		}
	} else {
		log.Info().Msg("Database migration completed successfully")
	}

}
