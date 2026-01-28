package repositories

import (
	"database/sql"
	"errors"
	"kasir_api_1/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// ==get_all_product======================//
func (repo *ProductRepository) GetAll() ([]models.Product, error) {
	//===query_join================//
	query := `
		SELECT
			p.id, p.name, p.price, p.stock, p.category_id,
			c.id, c.name, c.description
		FROM product p
		LEFT JOIN category c ON p.category_id = c.id
	`

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID,
			&p.Category.ID, &p.Category.Name, &p.Category.Description,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// ===getbyid================//
func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `
		SELECT
			p.id, p.name, p.price, p.stock, p.category_id,
			c.id, c.name, c.description
		FROM product p
		LEFT JOIN category c ON p.category_id = c.id
		WHERE p.id = $1
	`

	var p models.Product
	err := repo.db.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID,
		&p.Category.ID, &p.Category.Name, &p.Category.Description,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// ===create_product_saving_category=============//
func (repo *ProductRepository) Create(product *models.Product) error {
	query := `
		INSERT INTO product (name, price, stock, category_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := repo.db.QueryRow(
		query,
		product.Name,
		product.Price,
		product.Stock,
		product.CategoryID,
	).Scan(&product.ID)

	if err != nil {
		return err
	}

	fullData, err := repo.GetByID(product.ID)
	if err != nil {
		return nil
	}

	*product = *fullData

	return nil
}

// ===update_product============================//
func (repo *ProductRepository) Update(product *models.Product) error {
	//===add_category_id_to_update=====//
	query := `
		UPDATE product
		SET name = $1, price = $2, stock = $3, category_id = $4
		WHERE id = $5
	`

	result, err := repo.db.Exec(
		query,
		product.Name,
		product.Price,
		product.Stock,
		product.CategoryID,
		product.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM product WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return err
}
