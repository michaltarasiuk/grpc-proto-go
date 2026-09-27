.PHONY: gen tidy run-ping run-caller demo

# Regenerate Go stubs from .proto files.
gen:
	protoc \
		-I api \
		--go_out=gen --go_opt=paths=source_relative \
		--go-grpc_out=gen --go-grpc_opt=paths=source_relative \
		ping/v1/ping.proto

tidy:
	go mod tidy

run-ping:
	go run ./cmd/ping-server

run-caller:
	go run ./cmd/caller-server

demo:
	@echo "Terminal 1: make run-ping"
	@echo "Terminal 2: make run-caller"
	@echo "Then:       curl 'http://localhost:8080/call?msg=hi'"
