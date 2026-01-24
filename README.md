### Tugas 1 kelas Golang Kasir_API

## 🛠
- GET    /health              - Check server Health
- GET    /api/categories      - Get list categories
- POST   /api/categories      - Create categories
- GET    /api/categories/{id}     - Get categories by ID
- PUT    /api/categories/{id}     - Update categories by ID
- DELETE /api/categories/{id}     - Delete categories by ID

## 📦 Test
* GET Health
    - curl http://localhost:8080/health
* GET list all categories
    - curl http://localhost:8080/api/categories
* POST categories
    - curl -X POST http://localhost:8080/api/categories -H "Content-Type: application/json" -d '{"name": "Juice","description": "All you can Juice!!"}'
* GET categories detail by ID
    - curl http://localhost:8080/api/categories/{id}
* PUT update categories data by ID
    - curl -X PUT http://localhost:8080/api/categories -H "Content-Type: application/json" -d '{"name": "Juice","description": "All you can Juice!!"}'
* DELETE delete categories by ID (currently didnt reindexing after delete)
    - curl -X DELETE http://localhost:8080/api/categories/{id}