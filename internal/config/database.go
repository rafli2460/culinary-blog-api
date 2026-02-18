package config

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
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
