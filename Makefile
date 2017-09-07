all: deps protoc install

deps:
	go get google.golang.org/grpc
	go get github.com/golang/protobuf/protoc-gen-go

protoc:
	protoc --go_out=plugins=grpc:. examples/ping/protobuf/service.proto

install:
	go install ./examples/ping/esp4g-ping
	go install ./examples/ping/esp4g-ping-server
