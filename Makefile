gen-cal:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative caculator\caculatorpb\caculator.proto
run-server:
	go run caculator/server/server.go
run-client:
	go run caculator/client/client.go