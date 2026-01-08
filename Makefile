SRC_DIR=tiler

zip:
	cd $(SRC_DIR) && \
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap cmd/main.go && \
	zip tiler_code.zip bootstrap && \
	rm bootstrap
	@echo "Build complete: $(SRC_DIR)/tiler_code.zip"
