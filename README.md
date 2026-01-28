# 🛒 Kasir API (Point of Sale Backend)

**Tugas Kelas Golang** - A robust RESTful API built with **Golang** and **PostgreSQL** (Neon.tech) for managing a cashier system. This project demonstrates Clean Architecture, relational database management, and safe transaction handling.

![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Neon-336791?style=flat&logo=postgresql)
![License](https://img.shields.io/badge/License-MIT-green)

## 🌟 Key Features

* **Clean Architecture**: Separation of concerns (Handler -> Service -> Repository).
* **Relational Data**: Products are automatically linked to Categories via `LEFT JOIN`.
* **Safe Deletion**: Deleting a Category **does not** delete its products. Instead, products are safely moved to an "Uncategorized" system category (Transaction supported).
* **Auto-Reindexing**: ID sequence automatically resets after deletion to prevent huge ID gaps.
* **Smart Response**: Creating a product immediately returns the full object with category details.

## 🛠️ Tech Stack

* **Language**: Go (Golang)
* **Database**: PostgreSQL (hosted on [Neon.tech](https://neon.tech))
* **Configuration**: Viper (Environment variables)
* **Driver**: `lib/pq`

---

## 🚀 Getting Started

### Prerequisites

* Go installed (v1.18+)
* PostgreSQL database (Local or Cloud)

### Installation

1.  **Clone the repository**
    ```bash
    git clone [https://github.com/username/kasir_api_1.git](https://github.com/username/kasir_api_1.git)
    cd kasir_api_1
    ```

2.  **Setup Environment Variables**
    Create a `.env` file in the root directory:
    ```env
    PORT=8080
    DATABASE_URL="postgres://user:password@host/dbname?sslmode=require"
    ```

3.  **Run the Server**
    ```bash
    go mod tidy
    go run main.go
    ```
    You should see: `Whoosh!! 🚀 Server running on [ADDRESS]:8080`

---

## 📚 API Documentation

### 1️⃣ Health Check

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Check if API is running |

### 2️⃣ Categories (Tugas 1)

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/category` | Get list of all categories |
| `POST` | `/api/category` | Create a new category |
| `GET` | `/api/category/{id}` | Get specific category details |
| `PUT` | `/api/category/{id}` | Update category data |
| `DELETE` | `/api/category/{id}` | Delete category (Safe Delete) |

### 3️⃣ Products

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/product` | Get list of products (with Category info) |
| `POST` | `/api/product` | Create a new product |
| `GET` | `/api/product/{id}` | Get specific product details |
| `PUT` | `/api/product/{id}` | Update product data |
| `DELETE` | `/api/product/{id}` | Delete product |

---

## 📦 Test (cURL Examples)

You can copy and paste these commands into your terminal to test the API.

### Health Check
```bash
curl http://localhost:8080/health

---

## 🗄️ Database Schema

### Table: `category`
| Column | Type | Notes |
| :--- | :--- | :--- |
| `id` | SERIAL | Primary Key (Auto Increment) |
| `name` | VARCHAR | Category Name |
| `description` | TEXT | Optional description |

### Table: `product`
| Column | Type | Notes |
| :--- | :--- | :--- |
| `id` | SERIAL | Primary Key (Auto Increment) |
| `name` | VARCHAR | Product Name |
| `price` | INT | Product Price |
| `stock` | INT | Product Stock |
| `category_id` | INT | Foreign Key (Linked to `category.id`) |

> **Note:**
> * Relasi antar tabel menggunakan **Foreign Key** pada `category_id`.
> * Jika Kategori dihapus, produk **tidak** ikut terhapus (Logic *Safe Delete* di backend memindahkannya ke ID 1).

---

## 🤝 Contributing

Contributions are always welcome!

1.  Fork the repository
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request


Happy Coding! 👨‍💻👩‍💻
