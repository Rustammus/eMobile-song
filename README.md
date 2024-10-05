# Song app

---

#### API и модели данных описанны в docs, swagger доступен по пути /swagger/index.html

## START PROJECT
- **make app_dev.env and db_dev.env from example `example.env` file**

### Local:
- **install dependencies:**
```
go mod tidy
```
- **run mock song info service:**
```
go run ./mock/infoServer/main.go 
```
- **run project:**
```
go run ./cmd/main.go 
```
- ** swagger available on http://localhost:8082/swagger/index.html **
