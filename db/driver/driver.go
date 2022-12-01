package driver

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifeTime = 5 * time.Minute

func ConnectSQL(dsn string) (*DB, error) {
	conn, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	conn.SetMaxOpenConns(maxOpenDbConn)
	conn.SetMaxIdleConns(maxIdleDbConn)
	conn.SetConnMaxLifetime(maxDbLifeTime)

	if err := testDB(conn); err != nil {
		return nil, err
	}
	dbConn.SQL = conn
	return dbConn, nil
}

func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db, nil
}

func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}
