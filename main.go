package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"kasir_api_1/database"
	"kasir_api_1/handlers"
	"kasir_api_1/repositories"
	"kasir_api_1/services"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
}

func main() {
	// ======viper=================//
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:        viper.GetString("PORT"),
		DatabaseURL: viper.GetString("DATABASE_URL"),
	}

	//port========//
	if config.Port == "" {
		config.Port = "8080"
	}

	//==database_setup===//
	db, err := database.InitDB(config.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Gagal inisialisasi database: %v", err)
	}
	defer db.Close()

	// ======dependency_injection=================
	//====category=================//
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	//========product===============//
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	//=========routing============//
	//=====routes_category=============//
	http.HandleFunc("/api/category", categoryHandler.HandleCategory)
	http.HandleFunc("/api/category/", categoryHandler.HandleCategoryByID)

	//==routes_product===========//
	http.HandleFunc("/api/product", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductByID)

	//========transaction===========//
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo, productRepo, db)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)
	http.HandleFunc("/api/report-today", transactionHandler.HandleDailyReport)

	//==additional_endpoint_bulk_insert========//
	http.HandleFunc("/api/product/bulk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			productHandler.CreateBulk(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	//===health=============//
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Kasir Running",
		})
	})

	//=========health_server============///
	fmt.Println("----------------------------------------")
	fmt.Printf("Whoosh!! 🚀 Server running on port:%s\n", config.Port)
	fmt.Println("----------------------------------------")

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
