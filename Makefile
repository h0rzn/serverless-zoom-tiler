BIN_NAME=tiler
SRC_DIR=tiler

zip:
	cd $(SRC_DIR) && \
	go build -o $(BIN_NAME) cmd/main.go && \
	zip tiler_code.zip $(BIN_NAME) && \
	rm $(BIN_NAME)
	@echo "Build complete: $(SRC_DIR)/tiler_code.zip"
