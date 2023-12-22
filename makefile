swag:
	swag init --parseDependency --parseInternal --parseDepth 1 -d cmd/api

run:
	go run cmd/api/main.go

gen:
	go run cmd/generate/main.go

api:
	go run cmd/api/main.go

up:
	go run cmd/migrate/main.go -up

up_test:
	go run cmd/migrate/main.go -up -test

down:
	go run cmd/migrate/main.go -down

down_test:
	go run cmd/migrate/main.go -down -test

to:
	go run cmd/migrate/main.go -to -version 1

build-m:
	docker build -t account-migrate:latest -f deploy/migrate/linux/Dockerfile .

build-a:
	docker build -t account:latest -f deploy/api/linux/Dockerfile .

run-m:
	docker run --name account-migrate --rm --network="host" account-migrate:latest

run-a:
	docker run --name account -p 5487:5487 --network="host" -v ${PWD}/docker/log:/app/log account:latest


# docker exec -it <container-id> sh