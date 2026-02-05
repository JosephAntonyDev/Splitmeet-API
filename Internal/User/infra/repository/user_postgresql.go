package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/entities"
)

type UserPostgreSQLRepository struct {
	conn *core.Conn_PostgreSQL
}

func NewUserPostgreSQLRepository(conn *core.Conn_PostgreSQL) *UserPostgreSQLRepository {
	return &UserPostgreSQLRepository{conn: conn}
}

func (r *UserPostgreSQLRepository) Save(user *entities.User) error {
	query := `
		INSERT INTO users (name, email, password, phone, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id`

	err := r.conn.DB.QueryRow(
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Phone,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("error al insertar usuario en BD: %v", err)
	}

	return nil
}

func (r *UserPostgreSQLRepository) GetByEmail(email string) (*entities.User, error) {
	query := `SELECT id, name, email, password, phone, created_at, updated_at FROM users WHERE email = $1`

	row := r.conn.DB.QueryRow(query, email)

	var user entities.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error buscando usuario por email: %v", err)
	}

	return &user, nil
}

func (r *UserPostgreSQLRepository) GetByID(id int64) (*entities.User, error) {
	query := `SELECT id, name, email, password, phone, created_at, updated_at FROM users WHERE id = $1`

	row := r.conn.DB.QueryRow(query, id)

	var user entities.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error buscando usuario por ID: %v", err)
	}

	return &user, nil
}

func (r *UserPostgreSQLRepository) Update(user *entities.User) error {
	query := `
		UPDATE users 
		SET name = $1, phone = $2, password = $3, updated_at = $4 
		WHERE id = $5`
	
	user.UpdatedAt = time.Now()

	result, err := r.conn.DB.Exec(
		query,
		user.Name,
		user.Phone,
		user.Password,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("error actualizando usuario: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no se encontró usuario con id %d para actualizar", user.ID)
	}

	return nil
}

func (r *UserPostgreSQLRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.conn.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error eliminando usuario: %v", err)
	}
	return nil
}