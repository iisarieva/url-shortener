APP_NAME = server
ENTRYPOINT = cmd/server/main.go


# поднимаем проект с Redis
up:
	docker-compose up --build

# генерируем Swagger-документацию
swag:
	swag init -g $(ENTRYPOINT)

# запускаем сервер с предварительной генерацией Swagger
run: swag
	go run $(ENTRYPOINT)

# билдим Go-приложение (локально)
build:
	go build -o $(APP_NAME) $(ENTRYPOINT)

# модульная чистка
tidy:
	go mod tidy

# билдим docker-образ
docker:
	docker build -t url-shortener .

# выключаем контейнеры
down:
	docker-compose down

# чистим всё лишнее
clean:
	rm -f $(APP_NAME)
	rm -rf docs/swagger.* docs/docs.go
