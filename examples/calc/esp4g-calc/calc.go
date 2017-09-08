package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	calc "github.com/nokamoto/esp4g/examples/calc/protobuf"
	"google.golang.org/grpc"
	"io"
	"golang.org/x/net/context"
)

var (
	host = flag.String("h", "localhost", "The gRPC server host")
	port = flag.Int("p", 9000, "The gRPC server port")
	interval = flag.Int("n", 1, "Wait n seconds")
	apikey = flag.String("k", "guest", "The gRPC request 'x-api-key'")
)

type PerRPCCredentials struct {}

func (PerRPCCredentials)GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{"x-api-key": *apikey}, nil
}

func (PerRPCCredentials)RequireTransportSecurity() bool {
	return false
}

func clientSideStreaming(client calc.CalcServiceClient, count int64) {
	fmt.Print("client side: ")
	if stream, err := client.AddAll(context.Background()); err != nil {
		fmt.Println(err)
	} else {
		i := 0
		for i < 5 {
			x := count + int64(i * 2)
			y := count + int64(i * 2) + 1
			if i != 0 {
				fmt.Print(" + ")
			}
			fmt.Printf("%d + %d", x, y)
			if err = stream.Send(&calc.Operand{X: x, Y: y}); err != nil {
				fmt.Println(":", err)
				break
			}
			i = i + 1
		}
		if sum, err := stream.CloseAndRecv(); err != nil {
			fmt.Println(":", err)
		} else {
			fmt.Printf(" = %v\n", sum)
		}
	}
}

func serverSideStreaming(client calc.CalcServiceClient, count int64) {
	fmt.Print("server side: ")
	list := calc.OperandList{}
	i := 0
	for i < 5 {
		if i != 0 {
			fmt.Print(" + ")
		}
		x := count + int64(i * 2)
		y := count + int64(i * 2) + 1
		fmt.Printf("(%d + %d)", x, y)
		list.Operand = append(list.Operand, &calc.Operand{X: x, Y: y})
		i = i + 1
	}

	if stream, err := client.AddDeffered(context.Background(), &list); err != nil {
		fmt.Println(":", err)
	} else {
		fmt.Print(" = ")
		for {
			sum, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(":", err)
				return
			}
			fmt.Print(sum)
		}
		fmt.Println()
	}
}

func bidirectionalStreaming(client calc.CalcServiceClient, count int64) {
	fmt.Print("bidirectional: ")
	if stream, err := client.AddAsync(context.Background()); err != nil {
		fmt.Println(err)
	} else {
		i := int64(0)
		for i < 5 {
			x := count + int64(i * 2)
			y := count + int64(i * 2) + 1
			if i != 0 {
				fmt.Print(" + ")
			}
			fmt.Printf("%d + %d", x, y)
			if err = stream.Send(&calc.Operand{X: x, Y: y}); err != nil {
				fmt.Println(":", err)
				return
			}
			i = i + 1
		}

		fmt.Print(" = ")

		i = int64(0)
		for i < 5 {
			sum, err := stream.Recv()
			if err != nil {
				fmt.Println(":", err)
				return
			}
			fmt.Print(sum)
			i = i + 1
		}

		if err = stream.CloseSend(); err != nil {
			fmt.Println(":", err)
		} else {
			fmt.Println()
		}
	}
}

func main() {
	flag.Parse()

	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithPerRPCCredentials(PerRPCCredentials{})}

	con, err := grpc.Dial(fmt.Sprintf("%s:%d", *host, *port), opts...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer con.Close()

	client := calc.NewCalcServiceClient(con)

	count := int64(0)
	for count >= 0 {
		clientSideStreaming(client, count)
		serverSideStreaming(client, count)
		bidirectionalStreaming(client, count)

		time.Sleep(time.Duration(*interval) * time.Second)

		count = count + 1
	}
}
