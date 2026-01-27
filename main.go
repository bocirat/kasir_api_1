package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// =========struct_Categories==============//
type Categories struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// =========in_memory_data=====================//
var categories = []Categories{
	{ID: 1, Name: "Main Course", Description: "All Main Dish Special For You"},
	{ID: 2, Name: "Coffee", Description: "Special Coffee for your day"},
	{ID: 3, Name: "Snacks", Description: "Lite bites to cheer up your day"},
}

// =============main_start_here============================//
func main() {
	//GET_localhost:8080/api/categories/{id}
	//PUT_localhost:8080/api/categories/{id}
	//DELETE_localhost:8080/api/categories/{id}
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoriesID(w, r)
		} else if r.Method == "PUT" {
			updateCategories(w, r)
		} else if r.Method == "DELETE" {
			deleteCategories(w, r)
		}
	})

	//GET_localhost:8080/api/categories
	//POST_localhost:8080/api/categories
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)
		} else if r.Method == "POST" {
			var categoriesNew Categories
			err := json.NewDecoder(r.Body).Decode(&categoriesNew)
			if err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
				return
			}

			//===insert_categories_data_into_Categories===//
			categoriesNew.ID = len(categories) + 1
			categories = append(categories, categoriesNew)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) //==201==
			json.NewEncoder(w).Encode(categoriesNew)
		}
	})

	//localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API is running well",
		})
	})

	fmt.Println("Server is running on //localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to run server")
	}

}

//===============function_CRUD=============================//

// ====get_categories_ID==============//
func getCategoriesID(w http.ResponseWriter, r *http.Request) {
	//==parse_ID_from_URL==
	//==URL_is_/api/categories/{id}==
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Categories ID", http.StatusBadRequest)
		return
	}

	//=====search_categories_by_ID===========
	for _, p := range categories {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	//==if_not_found==
	http.Error(w, "Categories not found", http.StatusNotFound)
}

// ========update_categories=============//
func updateCategories(w http.ResponseWriter, r *http.Request) {
	//==get_ID_from_request==
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	//==change_int==
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Categories ID", http.StatusBadRequest)
		return
	}

	//==get_data_from_request==
	var updateCategories Categories
	err = json.NewDecoder(r.Body).Decode(&updateCategories)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	//==loop_find_categories_and_change_data
	for i := range categories {
		if categories[i].ID == id {
			updateCategories.ID = id
			categories[i] = updateCategories

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateCategories)
			return
		}
	}

	http.Error(w, "Categories not found", http.StatusNotFound)
}

// apply_mutex
var mu sync.Mutex

// =========delete_categories===================
func deleteCategories(w http.ResponseWriter, r *http.Request) {
	//==get_id==
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	//==change_id_int===
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Categories ID", http.StatusBadRequest)
		return
	}

	//mutex_implement_here
	mu.Lock()         //lock_access
	defer mu.Unlock() //Ensure_it_unlock_when_finished

	//==loop_find_categories_and_delete
	for i, p := range categories {
		if p.ID == id {
			//==slice_here==
			categories = append(categories[:i], categories[i+1:]...)

			//==fixing_logic_reindexing_before_return
			for i := range categories {
				categories[i].ID = i
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "success delete",
			})
			return
		}

	}

	http.Error(w, "Categories Not Found", http.StatusNotFound)
}
