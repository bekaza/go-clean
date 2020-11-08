run:
	go run main.go
test:
	go test ./...
mockgen:
	go run github.com/golang/mock/mockgen --build_flags=--mod=vendor -package mocks -destination=./domain/mock/user.go -source=./domain/user.go
