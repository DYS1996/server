package pgsql

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

// PGSQL contains db instance and its config
type PGSQL struct {
	instance *sql.DB
	config   *PGConfig
}

// New uses PGConfig, returns PGSQL instance
func New(c *PGConfig) (*PGSQL, error) {
	cfg, err := validConfig(c)
	if err != nil {
		return nil, fmt.Errorf("validate config: %v", err)
	}

	pgsql, err := sql.Open("postgres", `postgres://`+cfg.User+`:`+cfg.Pass+`@`+cfg.Host+`:`+strconv.Itoa(cfg.Port)+`/`+cfg.DBName+`?sslmode=disable&connect_timeout=5`)
	if err != nil {
		return nil, fmt.Errorf("open db: %v", err)
	}
	err = pgsql.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping db: %v", err)
	}
	return &PGSQL{instance: pgsql, config: cfg}, nil
}

// Close GPSQL
func (pg *PGSQL) Close() {
	pg.instance.Close()
}

func validConfig(r *PGConfig) (*PGConfig, error) {
	n := *r

	if n.Host == "" {
		n.Host = "localhost"
	}
	if n.Port == 0 {
		n.Port = 5432
	}
	if n.DBName == "" {
		n.DBName = "postgres"
	}
	if n.User == "" {
		return nil, errors.New("empty username")
	}
	if n.Pass == "" {
		return nil, errors.New("empty password")
	}
	return &n, nil
}
