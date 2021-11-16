package db

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// NewConnection returns a new database connection with the schema applied, if not already
// applied.
func NewConnection(conn string) (db *sqlx.DB, err error) {

	log.Info("connecting to postgres database...")
	if db, err = sqlx.Connect("postgres", conn); err != nil {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		for range ticker.C {
			if db, err = sqlx.Connect("postgres", conn); err == nil {
				break
			}
		}
	}
	log.Info("connected to postgres database")

	log.Info("verifying postgres connection...")
	if err := db.Ping(); err != nil {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		for range ticker.C {
			if err := db.Ping(); err == nil {
				break
			}
		}
	}
	log.Info("verified postgres connection")

	return db, nil
}
