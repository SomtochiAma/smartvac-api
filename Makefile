tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

test: tidy fmt vet
