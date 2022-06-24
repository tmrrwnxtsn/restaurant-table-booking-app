package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// NewDB устанавливает соединение с базой данных по переданной строке подключения.
func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
