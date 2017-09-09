all: deps protoc gtest install

deps:
	go get google.golang.org/grpc
	go get github.com/golang/protobuf/protoc-gen-go
	go get gopkg.in/yaml.v2
	go get go.uber.org/zap

protoc:
	protoc -I./examples/ping/protobuf --include_imports --include_source_info ./examples/ping/protobuf/service.proto --descriptor_set_out ./examples/ping/protobuf/descriptor.pb
	protoc -I./examples/calc/protobuf --include_imports --include_source_info ./examples/calc/protobuf/service.proto --descriptor_set_out ./examples/calc/protobuf/descriptor.pb
	protoc --go_out=plugins=grpc:. examples/ping/protobuf/service.proto
	protoc --go_out=plugins=grpc:. examples/calc/protobuf/service.proto
	protoc --go_out=plugins=grpc:. ./protobuf/service.proto

gtest:
	protoc -I./examples/ping/protobuf --include_imports --include_source_info ./examples/ping/protobuf/service.proto --descriptor_set_out ./esp4g-test/unary-descriptor.pb
	protoc -I./examples/calc/protobuf --include_imports --include_source_info ./examples/calc/protobuf/service.proto --descriptor_set_out ./esp4g-test/stream-descriptor.pb
	go test ./esp4g-test

install:
	go install ./examples/ping/esp4g-ping
	go install ./examples/ping/esp4g-ping-server
	go install ./examples/calc/esp4g-calc
	go install ./examples/calc/esp4g-calc-server
	go install ./esp4g
	go install ./esp4g-extension
