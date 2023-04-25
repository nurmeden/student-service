APP_NAME=app
DOCKER_TAG=latest

.PHONY: build docker clean run-mongo

build:
	go build -o app .

docker:
	docker build -t app:latest .

clean:
	rm -f app

run-mongo:
	docker run -p 27017:27017 -d --rm mongo:latest

run:
	docker run -p 8080:8080 -d --rm app