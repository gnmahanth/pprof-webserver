BINARY_NAME?=pprof-webserver
SRC_PATH=${PWD}
PROFILES_PATH?=${PWD}/samples

default: clean deps fmt vet test build

build:
	CGO_ENABLED=1 go build -v -race -trimpath -o $(SRC_PATH)/$(BINARY_NAME) $(SRC_PATH)/.
static:
	CGO_ENABLED=0 go build -v -a -trimpath -o $(SRC_PATH)/$(BINARY_NAME) $(SRC_PATH)/.
test:
	go test -v -race $(SRC_PATH)/...
cache-clean:
	go clean -modcache -cache -x $(SRC_PATH)
clean:
	go clean
	rm -f $(SRC_PATH)/$(BINARY_NAME)
run:
	$(SRC_PATH)/$(BINARY_NAME)
deps:
	go mod tidy -v
	go mod download
fmt:
	go fmt $(SRC_PATH)/...
vet:
	go vet $(SRC_PATH)/...
list:
	go list -m -u all
update:
	go get -u $(SRC_PATH)/...
lint:
	golangci-lint run --out-format tab
dev: clean deps fmt vet build
	./$(BINARY_NAME) --debug --storage $(PROFILES_PATH)
