APP_NAME := labor-app
TARGET_DIR := target
BINARY := $(TARGET_DIR)/$(APP_NAME)

.PHONY: run build clean

run:
	CGO_ENABLED=1 go run .

build:
	mkdir -p $(TARGET_DIR)
	CGO_ENABLED=1 go build -o $(BINARY) .

start: build
	./$(BINARY)

clean:
	rm -f $(APP_NAME)
	go clean
