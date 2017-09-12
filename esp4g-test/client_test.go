package esp4g_test

import (
	ping "github.com/nokamoto/esp4g/examples/ping/protobuf"
	calc "github.com/nokamoto/esp4g/examples/calc/protobuf"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"io"
)

func callPing(con *grpc.ClientConn, x *ping.Ping) (*ping.Pong, error) {
	c := ping.NewPingServiceClient(con)
	return c.Send(context.Background(), x)
}

func callCalcCStream(con *grpc.ClientConn, xs []*calc.Operand) (*calc.Sum, error) {
	c := calc.NewCalcServiceClient(con)

	stream, err := c.AddAll(context.Background())
	if err != nil {
		return nil, err
	}

	for _, x := range xs {
		// Send may return EOF.
		if err := stream.Send(x); err != nil {
			return nil, err
		}
	}

	return stream.CloseAndRecv()
}

func callCalcSStream(con *grpc.ClientConn, x *calc.OperandList) ([]*calc.Sum, error) {
	c := calc.NewCalcServiceClient(con)

	stream, err := c.AddDeffered(context.Background(), x)
	if err != nil {
		return nil, err
	}

	ys := []*calc.Sum{}

	for {
		y, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		ys = append(ys, y)
	}

	return ys, nil
}

func callCalcBStream(con *grpc.ClientConn, xs []*calc.Operand) ([]*calc.Sum, error) {
	c := calc.NewCalcServiceClient(con)

	stream, err := c.AddAsync(context.Background())
	if err != nil {
		return nil, err
	}

	ys := []*calc.Sum{}

	for _, x := range xs {
		// Send may return EOF.
		if err := stream.Send(x); err != nil {
			return nil, err
		}

		y, err := stream.Recv()
		if err != nil {
			return nil, err
		}
		ys = append(ys, y)
	}

	if err := stream.CloseSend(); err != nil {
		return nil, err
	}

	return ys, nil
}

