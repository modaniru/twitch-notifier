.PHONY: run
run: fmt
	go run cmd/main.go
fmt: 
	go fmt ./...
restart:
	docker compose up --force-recreate --build -d
	docker image prune -f