# Переменные
APP_NAME := app
DOCKER_IMAGE := student-service
DOCKER_TAG := latest

# Компиляция и сборка приложения
build:
	go build -o app cmd/main.go

# Запуск приложения
run:
	docker-compose up --build

build:
	GO111MODULE="on" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

# # Остановка и удаление Docker-контейнера
# docker-stop:
# 	docker stop $(APP_NAME)
# 	docker rm $(APP_NAME)

# # Очистка бинарных файлов и Docker-образов
# clean:
# 	rm -f $(APP_NAME)
# 	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG)
