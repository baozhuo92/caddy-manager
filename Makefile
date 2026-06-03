BINARY_NAME = caddy_server
BUILD_DIR = build
IMAGE_NAME = caddy-server

.PHONY: all mac linux windows clean docker-build

all: mac

mac:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

linux:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME) .

windows:
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME).exe .

docker-build: linux
	docker build -t $(IMAGE_NAME) .

clean:
	@rm -rf $(BUILD_DIR)
