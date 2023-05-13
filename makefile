APP_NAME=app
DOCKER_TAG=latest

.PHONY: build docker clean run-mongo

build:
	go build -o app .

<<<<<<< HEAD
# Запуск приложения
run:
	GO111MODULE="on" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

# Сборка Docker-образа
docker-build:
	docker build -t student-service .

# Запуск Docker-контейнера
docker-run:
	docker run -d --rm -p 8000:8000 --name student-service student-service
=======
docker:
	docker build -t app:latest .

clean:
	rm -f app

run-mongo:
	docker run -p 27017:27017 -d --rm mongo:latest

run:
	docker run -p 8080:8080 -d --rm app
>>>>>>> a30aacc52501219f0830333d30b82fb05993cfd4
