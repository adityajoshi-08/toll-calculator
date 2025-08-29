format:
	@gofmt -w .

obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver ./data_receiver
	@chmod +x ./bin/receiver
	@./bin/receiver

calculator:
	@go build -o bin/calculator ./distance_calculator
	@chmod +x ./bin/calculator
	@./bin/calculator

aggregator:
	@go build -o bin/aggregator ./aggregator
	@chmod +x ./bin/aggregator
	@./bin/aggregator

proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		types/ptypes.proto

.PHONY: obu receiver calculator aggregator format
