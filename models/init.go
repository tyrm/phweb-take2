package models

import (
	"database/sql"
	"errors"
	"github.com/gobuffalo/packr/v2"
	"github.com/juju/loggo"
	_ "github.com/lib/pq" // include for sql
	"github.com/rubenv/sql-migrate"
)

var db *sql.DB
var logger *loggo.Logger

var (
	ErrDoesNotExist = errors.New("does not exist")
	ErrInvalid      = errors.New("invalid")
)

func Close() {
	err := db.Close()
	if err != nil {
		logger.Errorf("Could not close database: %s", err.Error())
	}
}

func Init(connectionString string) error {
	newLogger := loggo.GetLogger("migrations")
	logger = &newLogger

	// Connect to Database
	logger.Debugf("Connecting to Database")
	dbClient, err := sql.Open("postgres", connectionString)
	if err != nil {
		logger.Criticalf("Could not connect to database: %s", err)
		return err
	}
	db = dbClient
	db.SetMaxIdleConns(5)

	// Do Migration
	logger.Debugf("Loading Migrations")
	migrations := &migrate.PackrMigrationSource{
		Box: packr.New("migrations", "./migrations"),
	}

	logger.Debugf("Applying Migrations")
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if n > 0 {
		logger.Infof("Applied %d migrations!\n", n)
	}
	if err != nil {
		logger.Criticalf("Coud not migrate database: %s", err)
		return err
	}

	return nil
}

func GetDBConn() *sql.DB {
	return db
}
