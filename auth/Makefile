run:
	@go run ./cmd/auth/main.go

migrate:
	@go run ./cmd/migrate/main.go

test:
	@go test ./... -v

MOCKS_DESTINATION=mocks

mocks: internal/repository/customer.go 	
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; done

proto: 
	@protoc -I protos/pb protos/pb/auth/v1/auth.proto --go_out=protos/gen --go_opt=paths=source_relative --go-grpc_out=protos/gen --go-grpc_opt=paths=source_relative

.PHONY: test
