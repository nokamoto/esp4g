package utils

import (
	"io/ioutil"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"fmt"
)

func ReadFileDescriptorSet(file string) (*descriptor.FileDescriptorSet, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	fds := &descriptor.FileDescriptorSet{}

	return fds, proto.Unmarshal(data, fds)
}

func GetFullMethodName(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto, method *descriptor.MethodDescriptorProto) string {
	return fmt.Sprintf("/%s.%s/%s", file.GetPackage(), service.GetName(), method.GetName())
}

func Methods(fds *descriptor.FileDescriptorSet) []string {
	methods := []string{}

	for _, file := range fds.GetFile() {
		for _, service := range file.GetService() {
			for _, method := range service.GetMethod() {
				methods = append(methods, GetFullMethodName(file, service, method))
			}
		}
	}

	return methods
}
