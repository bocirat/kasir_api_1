package repositories

import (
	"database/sql"
	"kasir_api_1/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *sql.Tx, t *models.Transaction) error {
	query := `
		INSERT INTO transaction (total_amount, created_at)
		VALUES ($1, $2)
		RETURNING id
	`
	// Use tx if provided, else db
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow(query, t.TotalAmount, t.CreatedAt)
	} else {
		row = r.db.QueryRow(query, t.TotalAmount, t.CreatedAt)
	}

	err := row.Scan(&t.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionRepository) CreateDetail(tx *sql.Tx, td *models.TransactionDetail) error {
	query := `
		INSERT INTO transaction_detail (transaction_id, product_id, quantity, subtotal)
		VALUES ($1, $2, $3, $4)
	`

	var err error
	if tx != nil {
		_, err = tx.Exec(query, td.TransactionID, td.ProductID, td.Quantity, td.Subtotal)
	} else {
		_, err = r.db.Exec(query, td.TransactionID, td.ProductID, td.Quantity, td.Subtotal)
	}

	return err
}

func (r *TransactionRepository) GetDailySummary(date string) (*models.DailyReport, error) {
	// 1. Get total revenue and transaction count
	queryTotal := `
		SELECT 
			COALESCE(SUM(total_amount), 0), 
			COUNT(id)
		FROM transaction 
		WHERE DATE(created_at) = $1
	`
	var report models.DailyReport
	report.Date = date

	err := r.db.QueryRow(queryTotal, date).Scan(&report.TotalRevenue, &report.TotalTransactions)
	if err != nil {
		return nil, err
	}

	// 2. Get product sales breakdown
	queryProducts := `
		SELECT 
			p.id, p.name, 
			SUM(td.quantity) as total_qty, 
			SUM(td.subtotal) as total_sales
		FROM transaction_detail td
		JOIN transaction t ON td.transaction_id = t.id
		JOIN product p ON td.product_id = p.id
		WHERE DATE(t.created_at) = $1
		GROUP BY p.id, p.name
	`

	rows, err := r.db.Query(queryProducts, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.ProductSales
	for rows.Next() {
		var p models.ProductSales
		if err := rows.Scan(&p.ProductID, &p.ProductName, &p.TotalQty, &p.TotalSales); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	report.SoldItems = products
	return &report, nil
}
