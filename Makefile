run:
	@echo 'Running server.'
	export `less .env | xargs`; go run cmd/geolocalization/main.go