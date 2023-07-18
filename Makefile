build:
	go build -o ./bin/gomake ./cmd/make.go

test:
	go test ./...