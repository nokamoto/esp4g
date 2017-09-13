package main

import (
	"flag"
	"net"
	"fmt"
	"log"
	"google.golang.org/grpc"
	calc "github.com/nokamoto/esp4g/examples/calc/protobuf"
	"io"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"golang.org/x/net/context"
	"github.com/golang/protobuf/ptypes/empty"
)

type CalcServer struct {}

func (CalcServer)AddAll(stream calc.CalcService_AddAllServer) error {
	sum := int64(0)
	for {
		operand, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("AddAll: %v", err)
			return status.Error(codes.Aborted, err.Error())
		}
		sum = sum + operand.GetX() + operand.GetY()
		log.Printf("AddAll: %v %v", operand, sum)
	}
	return stream.SendAndClose(&calc.Sum{Z: sum})
}

func (CalcServer)AddDeffered(list *calc.OperandList, stream calc.CalcService_AddDefferedServer) error {
	log.Printf("AddDeffered: %v", list)
	for _, operand := range list.Operand {
		if err := stream.Send(&calc.Sum{Z: operand.GetX() + operand.GetY()}); err != nil {
			log.Printf("AddDeffered: %v", err)
			return status.Error(codes.Aborted, err.Error())
		}
	}
	return nil
}

func (CalcServer)AddAsync(stream calc.CalcService_AddAsyncServer) error {
	for {
		operand, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("AddAsync: %v", err)
			return status.Error(codes.Aborted, err.Error())
		}
		log.Printf("AddAsync: %v", operand)
		if err = stream.Send(&calc.Sum{Z: operand.GetX() + operand.GetY()}); err != nil {
			log.Printf("AddAsync: %v", err)
			return status.Error(codes.Aborted, err.Error())
		}
	}
	return nil
}

func (CalcServer)Check(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

var (
	port = flag.Int("p", 8000, "The gRPC server port")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("listen %v port", *port)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	calc.RegisterCalcServiceServer(server, CalcServer{})
	calc.RegisterHealthCheckServiceServer(server, CalcServer{})

	log.Println("start esp4g-calc-server...")
	server.Serve(lis)
}
