package core

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Conn_PostgreSQL struct {
	DB *sql.DB
}

func GetDBPool() (*Conn_PostgreSQL, error) {

	dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        return nil, fmt.Errorf("la variable de entorno DB_URL está vacía")
    }

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la base de datos: %w", err)
	}

	db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
        db.Close()
		return nil, fmt.Errorf("error al verificar la conexión (ping): %w", err)
	}

    fmt.Println("Conexión a PostgreSQL exitosa")
	return &Conn_PostgreSQL{DB: db}, nil
}

// Wrapper simple para Exec (Insert, Update, Delete)
func (conn *Conn_PostgreSQL) Execute(query string, values ...interface{}) (sql.Result, error) {
    // DB.Exec ya maneja el prepare/exec internamente de forma optimizada para uso único
	result, err := conn.DB.Exec(query, values...)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando query: %w", err)
	}
	return result, nil
}

func (conn *Conn_PostgreSQL) Query(query string, values ...interface{}) (*sql.Rows, error) {
	rows, err := conn.DB.Query(query, values...)
	if err != nil {
		return nil, fmt.Errorf("error en select query: %w", err)
	}
	return rows, nil
}