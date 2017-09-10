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
