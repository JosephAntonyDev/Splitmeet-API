package repository

import (
	"database/sql"
	"fmt"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/JosephAntonyDev/splitmeet-api/internal/product/domain/entities"
)

type ProductPostgreSQLRepository struct {
	conn *core.Conn_PostgreSQL
}

func NewProductPostgreSQLRepository(conn *core.Conn_PostgreSQL) *ProductPostgreSQLRepository {
	return &ProductPostgreSQLRepository{conn: conn}
}

func (r *ProductPostgreSQLRepository) GetByID(id int64) (*entities.Product, error) {
	query := `
		SELECT id, category_id, name, presentation, size, default_price, is_predefined, created_by, created_at 
		FROM products 
		WHERE id = $1`

	row := r.conn.DB.QueryRow(query, id)

	var product entities.Product
	var categoryID, createdBy sql.NullInt64
	var presentation, size sql.NullString
	var defaultPrice sql.NullFloat64

	err := row.Scan(
		&product.ID,
		&categoryID,
		&product.Name,
		&presentation,
		&size,
		&defaultPrice,
		&product.IsPredefined,
		&createdBy,
		&product.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error al buscar producto por ID: %v", err)
	}

	if categoryID.Valid {
		product.CategoryID = &categoryID.Int64
	}
	if presentation.Valid {
		product.Presentation = presentation.String
	}
	if size.Valid {
		product.Size = size.String
	}
	if defaultPrice.Valid {
		product.DefaultPrice = &defaultPrice.Float64
	}
	if createdBy.Valid {
		product.CreatedBy = &createdBy.Int64
	}

	return &product, nil
}

func (r *ProductPostgreSQLRepository) GetByCategory(categoryID int64) ([]entities.Product, error) {
	query := `
		SELECT id, category_id, name, presentation, size, default_price, is_predefined, created_by, created_at 
		FROM products 
		WHERE category_id = $1 
		ORDER BY is_predefined DESC, name ASC`

	rows, err := r.conn.DB.Query(query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener productos por categoría: %v", err)
	}
	defer rows.Close()

	return r.scanProducts(rows)
}

func (r *ProductPostgreSQLRepository) Search(searchQuery string, categoryID *int64) ([]entities.Product, error) {
	var query string
	var args []interface{}

	if categoryID != nil {
		query = `
			SELECT id, category_id, name, presentation, size, default_price, is_predefined, created_by, created_at 
			FROM products 
			WHERE (LOWER(name) LIKE LOWER($1) OR LOWER(presentation) LIKE LOWER($1)) 
			AND category_id = $2
			ORDER BY is_predefined DESC, name ASC`
		args = []interface{}{"%" + searchQuery + "%", *categoryID}
	} else {
		query = `
			SELECT id, category_id, name, presentation, size, default_price, is_predefined, created_by, created_at 
			FROM products 
			WHERE LOWER(name) LIKE LOWER($1) OR LOWER(presentation) LIKE LOWER($1)
			ORDER BY is_predefined DESC, name ASC`
		args = []interface{}{"%" + searchQuery + "%"}
	}

	rows, err := r.conn.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error al buscar productos: %v", err)
	}
	defer rows.Close()

	return r.scanProducts(rows)
}

func (r *ProductPostgreSQLRepository) Save(product *entities.Product) error {
	query := `
		INSERT INTO products (category_id, name, presentation, size, default_price, is_predefined, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id`

	var categoryID, createdBy sql.NullInt64
	var presentation, size sql.NullString
	var defaultPrice sql.NullFloat64

	if product.CategoryID != nil {
		categoryID = sql.NullInt64{Int64: *product.CategoryID, Valid: true}
	}
	if product.Presentation != "" {
		presentation = sql.NullString{String: product.Presentation, Valid: true}
	}
	if product.Size != "" {
		size = sql.NullString{String: product.Size, Valid: true}
	}
	if product.DefaultPrice != nil {
		defaultPrice = sql.NullFloat64{Float64: *product.DefaultPrice, Valid: true}
	}
	if product.CreatedBy != nil {
		createdBy = sql.NullInt64{Int64: *product.CreatedBy, Valid: true}
	}

	err := r.conn.DB.QueryRow(
		query,
		categoryID,
		product.Name,
		presentation,
		size,
		defaultPrice,
		product.IsPredefined,
		createdBy,
		product.CreatedAt,
	).Scan(&product.ID)

	if err != nil {
		return fmt.Errorf("error al insertar producto: %v", err)
	}
	return nil
}

func (r *ProductPostgreSQLRepository) scanProducts(rows *sql.Rows) ([]entities.Product, error) {
	var products []entities.Product

	for rows.Next() {
		var product entities.Product
		var categoryID, createdBy sql.NullInt64
		var presentation, size sql.NullString
		var defaultPrice sql.NullFloat64

		err := rows.Scan(
			&product.ID,
			&categoryID,
			&product.Name,
			&presentation,
			&size,
			&defaultPrice,
			&product.IsPredefined,
			&createdBy,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear producto: %v", err)
		}

		if categoryID.Valid {
			product.CategoryID = &categoryID.Int64
		}
		if presentation.Valid {
			product.Presentation = presentation.String
		}
		if size.Valid {
			product.Size = size.String
		}
		if defaultPrice.Valid {
			product.DefaultPrice = &defaultPrice.Float64
		}
		if createdBy.Valid {
			product.CreatedBy = &createdBy.Int64
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar productos: %v", err)
	}

	return products, nil
}
