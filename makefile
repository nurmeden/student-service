# Переменные
APP_NAME := app
DOCKER_IMAGE := student-service
DOCKER_TAG := latest

# Компиляция и сборка приложения
build:
	go build -o app cmd/main.go

# Запуск приложения
run:
	GO111MODULE="on" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

# Сборка Docker-образа
docker-build:
	docker build -t student-service .

# Запуск Docker-контейнера
docker-run:
	docker run -d --rm -p 8000:8000 --name student-service student-service
