package database

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	zero "github.com/rs/zerolog/log"

	"github.com/rafli2460/culinary-blog-api/internal/config"
	"github.com/rafli2460/culinary-blog-api/internal/server"
)

func Init(conf map[string]string) *server.Datasources {
	var err error
	var dbWriter *sqlx.DB
	var dbReader *sqlx.DB

	dsWriter, dsReader := parseDs(conf)

	if dbWriter, err = sqlx.Connect(conf[config.DbDialeg], dsWriter); err == nil {
		dbWriter.SetConnMaxLifetime(time.Duration(1) * time.Second)
		dbWriter.SetMaxOpenConns(10)
		dbWriter.SetMaxIdleConns(10)

		zero.Log().Msg("Initializing Writer DB: Pass")
	} else {
		zero.Panic().
			Str("Context", "Connecting to Writer DB").
			Err(err).Msg("")
	}

	if dbReader, err = sqlx.Connect(conf[config.DbDialeg], dsReader); err == nil {
		dbReader.SetConnMaxLifetime(time.Duration(1) * time.Second)
		dbReader.SetMaxOpenConns(10)
		dbReader.SetMaxIdleConns(10)

		zero.Log().Msg("Initializing Reader DB: Pass")
	} else {
		zero.Panic().
			Str("Context", "Connecting to Reader DB").
			Err(err).Msg("")
	}

	ds := &server.Datasources{
		WriterDB: dbWriter,
		ReaderDB: dbReader,
	}

	return ds

}

func parseDs(conf map[string]string) (dsWriter, dsReader string) {
	hostWriter := conf[config.DbHostWriter]
	hostReader := conf[config.DbHostReader]
	port := conf[config.DbPort]
	user := conf[config.DbUser]
	pass := conf[config.DbPass]
	name := conf[config.DbName]

	dsWriter = fmt.Sprintf("%s:%s@(%s:%s)/%s", user, pass, hostWriter, port, name)
	dsReader = fmt.Sprintf("%s:%s@(%s:%s)/%s", user, pass, hostReader, port, name)

	return
}
