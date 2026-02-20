package platform

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type PostgresDB struct {
	Connection *sql.DB
}

func newConnection(ctx context.Context, con_str string) (*sql.DB, error) {
	if strings.TrimSpace(con_str) == "" {
		return nil, fmt.Errorf("Database connection string is empty")
	}

	db, err := sql.Open("postgres", con_str)
	if err != nil {
		return nil, fmt.Errorf("Postgres connection was not established with :\n%v", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 10)

	return db, nil
}

func (p *PostgresDB) Close() {
	// wg.Add(1)
	// defer wg.Done()

	p.Connection.Close()
}

func Instrument(ctx context.Context, connection_string string) (*PostgresDB, error) {
	// wg.Add(1)
	// defer wg.Done()

	var db *sql.DB
	var err error

	for i := range 5 {
		fmt.Printf("Attempting to open postgres connection: attempt#[%d]\n", i+1)
		db, err = newConnection(ctx, connection_string)
		if err == nil {
			break
		}

		time.Sleep(time.Second * time.Duration(i))
	}

	if err != nil {
		return nil, err
	}

	fmt.Printf("Postgres connection established successfully")
	return &PostgresDB{
		Connection: db,
	},nil 
}
