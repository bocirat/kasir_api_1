package repositories

import (
	"database/sql"
	"errors"
	"kasir_api_1/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// ===getall_category====================//
func (repo *CategoryRepository) GetAll() ([]models.Category, error) {
	query := "SELECT id, name, description FROM category"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

// ==getbyid_category====================//
func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	query := "SELECT id, name, description FROM category WHERE id = $1"

	var c models.Category
	err := repo.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Description)

	if err == sql.ErrNoRows {
		return nil, errors.New("kategori tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// ===create_category====================//
func (repo *CategoryRepository) Create(category *models.Category) error {
	query := "INSERT INTO category (name, description) VALUES ($1, $2) RETURNING id"

	err := repo.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
	return err
}

// ==update_category====================//
func (repo *CategoryRepository) Update(category *models.Category) error {
	query := "UPDATE category SET name = $1, description = $2 WHERE id = $3"

	result, err := repo.db.Exec(query, category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("kategori tidak ditemukan")
	}

	return nil
}

// ===delete_category,move_product,reindex_ID====================//
func (repo *CategoryRepository) Delete(id int) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queryMoveProducts := "UPDATE product SET category_id = 1 WHERE category_id = $1"
	_, err = tx.Exec(queryMoveProducts, id)
	if err != nil {
		return err
	}

	queryDeleteCategory := "DELETE FROM category WHERE id = $1"
	result, err := tx.Exec(queryDeleteCategory, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("kategori tidak ditemukan")
	}

	//===reset_sequence_id=====================//
	queryResetSeq := `
		SELECT setval(pg_get_serial_sequence('category', 'id'),
		COALESCE((SELECT MAX(id) FROM category), 1));
	`
	_, err = tx.Exec(queryResetSeq)
	if err != nil {

		return err
	}
	return tx.Commit()
}
