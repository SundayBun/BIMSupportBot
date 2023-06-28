package repository

import (
	"BIMSupportBot/config"
	"fmt"
	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
	"time"
)

func NewPsqlDB(c *config.Config) (*sqlx.DB, error) {
	sslMode := "disable"
	if c.Postgres.SSLMode {
		sslMode = "enable"
	}

	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.DbName,
		sslMode,
		c.Postgres.Password)

	db, err := sqlx.Connect(c.Postgres.Driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(c.Postgres.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(c.Postgres.ConnMaxLifetime) * time.Second)
	db.SetMaxIdleConns(c.Postgres.MaxIdleConns)
	db.SetConnMaxIdleTime(time.Duration(c.Postgres.ConnMaxIdleTime) * time.Second)
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
