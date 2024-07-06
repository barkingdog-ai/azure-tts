test:
	go test ./...

lint: 
	golangci-lint run 
format: 
	goimports -l -w .
	gofumpt -l -w .
FORCE: ;