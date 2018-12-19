serve:
	PORT=50000 go run server.go

protoc:
	protoc services/hello/*.proto --go_out=plugins=grpc:codegen/go --js_out=codegen/js --python_out=codegen/python

grpcc:
	grpcc --proto ./service/hello/hello.proto --address 127.0.0.1:50000

call-from-go:
	ENDPOINT=localhost:50000 go run clients/client.go

call-from-js:
	ENDPOINT=localhost:50000 node clients/client.js

call-from-python:
	ENDPOINT=localhost:50000 python clients/client.js
