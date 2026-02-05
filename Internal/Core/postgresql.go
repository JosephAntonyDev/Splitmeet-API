package core

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Conn_PostgreSQL struct {
	DB *sql.DB
    // Quitamos 'Err string'. Si falla, la función constructora devuelve error.
}

// Devuelve (*Conn_PostgreSQL, error) -> Patrón estándar de Go
func GetDBPool() (*Conn_PostgreSQL, error) {
    // Asumimos que godotenv.Load() ya se hizo en el main.go
    // para no recargar el archivo en cada conexión.

	dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        return nil, fmt.Errorf("la variable de entorno DB_URL está vacía")
    }

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la base de datos: %w", err)
	}

    // Configuración recomendada para producción
	db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5) // Mantener algunas libres listas para usar

	if err := db.Ping(); err != nil {
        db.Close() // Importante cerrar si el ping falla
		return nil, fmt.Errorf("error al verificar la conexión (ping): %w", err)
	}

    fmt.Println("✅ Conexión a PostgreSQL exitosa")
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

// Wrapper para Query (Select). AHORA DEVUELVE ERROR.
func (conn *Conn_PostgreSQL) Query(query string, values ...interface{}) (*sql.Rows, error) {
	rows, err := conn.DB.Query(query, values...)
	if err != nil {
		return nil, fmt.Errorf("error en select query: %w", err)
	}
	return rows, nil
}