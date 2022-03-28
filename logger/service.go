package main

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/geoffjay/plantd/core/util"
	"github.com/geoffjay/plantd/logger/db"
	"github.com/geoffjay/plantd/logger/db/migrations"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	db         *sql.DB
	migrations []db.Migration
}

// NewService constructs and instance of a service type.
func NewService() *Service {
	return &Service{
		migrations: make([]db.Migration, 0),
	}
}

func (s *Service) run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	s.openDB()
	defer s.db.Close()

	s.runMigrations()

	go func() {
		log.Debug("run logging service")
	}()

	<-ctx.Done()
	log.Debug("exiting service")
}

func connectionStr() string {
	host := util.Getenv("PLANTD_LOGGER_TSDB_HOST", "localhost")
	port := util.Getenv("PLANTD_LOGGER_TSDB_PORT", "5432")
	user := util.Getenv("PLANTD_LOGGER_TSDB_USER", "admin")
	password := util.Getenv("PLANTD_LOGGER_TSDB_PASSWORD", "plantd")
	database := util.Getenv("PLANTD_LOGGER_TSDB_DATABASE", "plantd_development")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		database,
	)
}

func (s *Service) openDB() {
	var err error

	log.Debug("connecting to database")

	s.db, err = sql.Open("postgres", connectionStr())
	if err != nil {
		log.Panic(err)
	}
}

func (s *Service) runMigrations() {
	log.Debug("loading database migrations")
	s.registerMigration(&migrations.CreateMetricsTable{DB: s.db})

	log.Debug("running database migrations")

	// XXX: should these rollback if one fails?
	for _, migration := range s.migrations {
		if err := migration.Up(); err != nil {
			log.Panic(err)
		}
	}
}

func (s *Service) registerMigration(migration db.Migration) {
	s.migrations = append(s.migrations, migration)
}
