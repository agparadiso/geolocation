run:
	@echo 'Running server.'
	export `less .env | xargs`; go run cmd/geolocation/main.go

run-local:
	@echo 'Running server locally.'
	docker-compose down
	docker-compose build
	docker-compose up

build:
	go build ./cmd/geolocation

test:
	go test ./...