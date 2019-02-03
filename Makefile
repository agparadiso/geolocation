run:
	@echo 'Running server.'
	export `less .env | xargs`; go run cmd/geolocation/main.go

build:
	go build ./cmd/geolocation