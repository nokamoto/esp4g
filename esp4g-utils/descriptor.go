package utils

import (
	"io/ioutil"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func ReadFileDescriptorSet(file string) (*descriptor.FileDescriptorSet, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	fds := &descriptor.FileDescriptorSet{}

	return fds, proto.Unmarshal(data, fds)
}
