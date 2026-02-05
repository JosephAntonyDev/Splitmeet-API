package repository

import (
	"database/sql"
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/category/domain/entities"
	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
)

type CategoryPostgreSQLRepository struct {
	conn *core.Conn_PostgreSQL
}

func NewCategoryPostgreSQLRepository(conn *core.Conn_PostgreSQL) *CategoryPostgreSQLRepository {
	return &CategoryPostgreSQLRepository{conn: conn}
}

func (r *CategoryPostgreSQLRepository) GetAll() ([]entities.Category, error) {
	query := `SELECT id, name, icon, is_active, created_at FROM categories WHERE is_active = true ORDER BY id`

	rows, err := r.conn.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener categorías: %v", err)
	}
	defer rows.Close()

	var categories []entities.Category
	for rows.Next() {
		var category entities.Category
		var icon sql.NullString

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&icon,
			&category.IsActive,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear categoría: %v", err)
		}

		if icon.Valid {
			category.Icon = icon.String
		}

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar categorías: %v", err)
	}

	return categories, nil
}

func (r *CategoryPostgreSQLRepository) GetByID(id int64) (*entities.Category, error) {
	query := `SELECT id, name, icon, is_active, created_at FROM categories WHERE id = $1 AND is_active = true`

	row := r.conn.DB.QueryRow(query, id)

	var category entities.Category
	var icon sql.NullString

	err := row.Scan(
		&category.ID,
		&category.Name,
		&icon,
		&category.IsActive,
		&category.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error al buscar categoría por ID: %v", err)
	}

	if icon.Valid {
		category.Icon = icon.String
	}

	return &category, nil
}
