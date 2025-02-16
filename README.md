# Инструкция по запуску сервиса через Docker

### Установка Docker и Docker Compose
Перед запуском убедитесь, что у вас установлены:
- **Docker**: [Установка Docker](https://docs.docker.com/engine/install/)
- **Docker Compose**: [Установка Docker Compose](https://docs.docker.com/compose/install/)

Проверьте их версии:
```sh
docker --version
docker-compose --version
```

### Клонирование репозитория
```sh
git clone https://github.com/robertshagiev/avito-merch-shop.git
cd avito-merch-store
```
### Запуск контейнеров
Запускаем сервис и базу данных PostgreSQL:
```sh
docker-compose up --build -d
```

### Накатывание миграций с Goose
Для управления миграциями используется Goose

### Установка goose
``` sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```
Добавьте ```~/go/bin``` в PATH, если команда goose не работает
```sh
export PATH=$PATH:$(go env GOPATH)/bin
```
### Применение миграций
```
goose -dir ./migrations postgres "postgres://postgres:password@localhost:5432/shop?sslmode=disable" up
```
### Откат миграций (если нужно)
```sh
goose -dir ./migrations postgres "postgres://postgres:password@localhost:5432/shop?sslmode=disable" down
```
```
--build — пересобирает контейнер перед запуском.
-d — запускает контейнеры в фоновом режиме.
```
### После запуска можно проверить состояние контейнеров:
```
docker ps
```
