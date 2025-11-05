package db

import (
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type Adapter struct {
	db *sql.DB
}

func NewAdapter(driverName, dataSourceName string) (*Adapter, error) {
	// connect
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal("db connection failure: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("db ping failure: %v", err)
	}
	// create table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS arith_history (date DATETIME, answer INTEGER, operation TEXT)`)
	if err != nil {
		log.Fatal("table creation failure: %v", err)
	}
	return &Adapter{db: db}, nil
}

func (da Adapter) CloseDbConnection() {
	err := da.db.Close()
	if err != nil {
		{
			log.Fatalf("db close failure %v", err)
		}
	}
}

func (da Adapter) AddToHistory(answer int32, operation string) error {
	queryString, args, err := sq.Insert("arith_history").Columns("date", "answer", "Operation").Values(time.Now(), answer, operation).ToSql()
	if err != nil {
		return err
	}
	_, err = da.db.Exec(queryString, args...)
	if err != nil {
		return err
	}
	return nil
}
