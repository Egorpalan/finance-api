run:
	docker-compose up --build
	goose -dir=migrations postgres "postgres://postgres:postgres@localhost:5436/finance?sslmode=disable" down
	goose -dir=migrations postgres "postgres://postgres:postgres@localhost:5436/finance?sslmode=disable" up



test:
	go test -v ./tests/...
