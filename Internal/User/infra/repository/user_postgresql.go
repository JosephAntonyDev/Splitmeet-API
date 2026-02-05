package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"

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
		INSERT INTO users (username, name, email, password, phone, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id`

	err := r.conn.DB.QueryRow(
		query,
		user.Username,
		user.Name,
		user.Email,
		user.Password,
		user.Phone,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("error al insertar usuario: %v", err)
	}
	return nil
}

func (r *UserPostgreSQLRepository) GetByEmail(email string) (*entities.User, error) {
	query := `SELECT id, username, name, email, password, phone, created_at, updated_at FROM users WHERE email = $1`

	row := r.conn.DB.QueryRow(query, email)

	var user entities.User
	err := row.Scan(
		&user.ID,
		&user.Username,
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
	query := `SELECT id, username, name, email, password, phone, created_at, updated_at FROM users WHERE id = $1`

	row := r.conn.DB.QueryRow(query, id)

	var user entities.User
	err := row.Scan(
		&user.ID,
		&user.Username,
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
		SET username = $1, email = $2, name = $3, phone = $4, password = $5, updated_at = $6 
		WHERE id = $7`

	user.UpdatedAt = time.Now()

	result, err := r.conn.DB.Exec(
		query,
		user.Username,
		user.Email,
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

func (r *UserPostgreSQLRepository) GetUsersByIDs(ids []int64) ([]entities.User, error) {

	query := `
        SELECT id, username, name, email, phone 
        FROM users 
        WHERE id = ANY($1)`

	rows, err := r.conn.DB.Query(query, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("error buscando usuarios por lote: %v", err)
	}
	defer rows.Close()

	var users []entities.User

	for rows.Next() {
		var u entities.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Name, &u.Email, &u.Phone); err != nil {
			return nil, fmt.Errorf("error escaneando usuario: %v", err)
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserPostgreSQLRepository) GetByUsername(username string) (*entities.User, error) {
	query := `SELECT id, username, name, email, phone, password, created_at, updated_at FROM users WHERE username = $1`
	row := r.conn.DB.QueryRow(query, username)

	var user entities.User
	err := row.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Phone, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error buscando por username: %v", err)
	}
	return &user, nil
}

func (r *UserPostgreSQLRepository) SearchByUsername(query string, limit int) ([]entities.User, error) {
	sqlQuery := `SELECT id, username, name, email, phone FROM users WHERE username ILIKE $1 LIMIT $2`

	rows, err := r.conn.DB.Query(sqlQuery, "%"+query+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("error buscando usuarios: %v", err)
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var u entities.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Name, &u.Email, &u.Phone); err != nil {
			return nil, fmt.Errorf("error escaneando usuario: %v", err)
		}
		users = append(users, u)
	}

	return users, nil
}
