all: vet fmt build

build-quick:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/ipamas ./cmd/main
	docker build -f build-tools/Dockerfile -t ipamas .

build:
	go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/ipamas ./cmd/main
	docker build -f build-tools/Dockerfile -t ipamas .

vet:
	go vet ./...

fmt:
	go list -f '{{.Dir}}' ./... | grep -v /vendor/ | xargs -L1 gofmt -l

tidy:
	go mod tidy

vendor:
	go mod vendor

protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        api/ipam.proto
