all: vet fmt build

build-quick:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/ipam-as ./cmd/main
	docker build -f build-tools/Dockerfile -t ipam-as .

build:
	go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/ipam-as ./cmd/main
	docker build -f build-tools/Dockerfile -t ipam-as .

vet:
	go vet ./...

fmt:
	go list -f '{{.Dir}}' ./... | grep -v /vendor/ | xargs -L1 gofmt -l

tidy:
	go mod tidy

vendor:
	go mod vendor
