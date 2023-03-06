package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dafaath/iot-server/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

var databaseUrl string

func init() {
	config := configs.GetConfig()
	databaseConfig := config.Database
	databaseUrl = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", databaseConfig.Username, databaseConfig.Password, databaseConfig.Host, databaseConfig.Port, databaseConfig.Name)
}

// This function will make a connection to the database only once.
func GetConnection() (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseUrl)
	config.MinConns = 5
	config.MaxConns = 20
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 10
	if err != nil {
		return nil, fmt.Errorf("error parsing database config %w", err)
	}

	// this returns connection pool
	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("pgx.Connect %w", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("pgx.Connect %w", err)
	}
	log.Println("Get connection from database")

	return conn, nil
}
