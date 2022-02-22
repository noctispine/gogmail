testdb:
	go test ./db -v -cover ./...
build:
	go build -o email
