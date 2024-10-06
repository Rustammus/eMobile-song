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
- **swagger available on http://localhost:8082/swagger/index.html**

### Docker:
- **build image**
```
docker build -t obuhovskaia11/emobile-song:latest .
```
- **or pull from Docker Hub**
```
docker pull obuhovskaia11/emobile-song:latest
```
- **run container**
```
docker run -d --name serverSong -p 8082:8082 --env-file app_dev.env obuhovskaia11/emobile-song:latest
```

### Docker compose(app and database):
- **run docker-compose**
```
docker-compose up --build
```

### Kubernetes (app and service):
- **copy manifests from `/k8s/`**
- **create database resource:**
```
kubectl create -f ./database/db-config.yaml
kubectl create -f ./database/db-deployment.yaml
kubectl create -f ./database/db-service.yaml
```
- **create app resources:**
- **set APP_HOST="your_node_external_ip" env in `app-config.yaml`**
```
kubectl create -f ./app/app-config.yaml
kubectl create -f ./app/app-deployment.yaml
kubectl create -f ./app/app-service.yaml
```