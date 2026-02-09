package models

type ProductSales struct {
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	TotalQty    int    `json:"total_qty"`
	TotalSales  int    `json:"total_sales"`
}

type DailyReport struct {
	Date              string         `json:"date"`
	TotalRevenue      int            `json:"total_revenue"`
	TotalTransactions int            `json:"total_transactions"`
	SoldItems         []ProductSales `json:"sold_items"`
}
