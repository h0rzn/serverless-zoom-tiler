SRC_DIR=tiler
TF_DIR=architecture

# start Localstack container for local dev
up:
	docker-compose up -d

# stop Localstack container
down:
	docker-compose down

zip:
	cd $(SRC_DIR) && \
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap cmd/main.go && \
	zip tiler_code.zip bootstrap && \
	mv tiler_code.zip ../tiler_code.zip && \
	rm bootstrap
	@echo "Build complete: $(SRC_DIR)/tiler_code.zip"


# init terraform
init:
	terraform -chdir=$(TF_DIR) init

# deploy with terraform
deploy: zip
	terraform -chdir=$(TF_DIR) apply -auto-approve -replace="aws_lambda_function.tiling_handler"
