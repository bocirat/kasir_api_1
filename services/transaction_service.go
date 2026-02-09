package services

import (
	"database/sql"
	"errors"
	"kasir_api_1/models"
	"kasir_api_1/repositories"
	"time"
)

type TransactionService struct {
	repo        *repositories.TransactionRepository
	productRepo *repositories.ProductRepository
	db          *sql.DB
}

func NewTransactionService(repo *repositories.TransactionRepository, productRepo *repositories.ProductRepository, db *sql.DB) *TransactionService {
	return &TransactionService{
		repo:        repo,
		productRepo: productRepo,
		db:          db,
	}
}

func (s *TransactionService) Checkout(request models.CheckoutRequest) (*models.Transaction, error) {
	// 1. Start Database Transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	// Make sure rollback if panic or error occurs
	defer tx.Rollback()

	// 2. Prepare Transaction object
	transaction := &models.Transaction{
		CreatedAt: time.Now(),
	}

	totalAmount := 0
	var details []models.TransactionDetail

	// test fixing here.
	for _, item := range request.Items {
		// get product data
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, errors.New("product not found: " + err.Error())
		}

		subtotal := product.Price * item.Quantity
		totalAmount += subtotal

		detail := models.TransactionDetail{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		}
		details = append(details, detail)

		// 4. Decrease Stock (using transaction tx)
		err = s.productRepo.DecreaseStock(tx, product.ID, item.Quantity)
		if err != nil {
			return nil, errors.New("failed to decrease stock " + product.Name + ": " + err.Error())
		}
	}

	transaction.TotalAmount = totalAmount

	// 5. Save Transaction Header
	err = s.repo.Create(tx, transaction)
	if err != nil {
		return nil, err
	}

	// 6. Save Transaction Details
	for i := range details {
		details[i].TransactionID = transaction.ID
		err = s.repo.CreateDetail(tx, &details[i])
		if err != nil {
			return nil, err
		}
	}

	transaction.Details = details

	// 7. Commit Transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) GetDailyReport(date string) (*models.DailyReport, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return s.repo.GetDailySummary(date)
}
