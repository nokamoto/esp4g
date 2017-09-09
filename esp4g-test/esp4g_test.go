package esp4g_test

import (
	"testing"
	"google.golang.org/grpc"
	ping "github.com/nokamoto/esp4g/examples/ping/protobuf"
	calc "github.com/nokamoto/esp4g/examples/calc/protobuf"
	"golang.org/x/net/context"
	"io"
	"reflect"
)

const UNARY_DESCRIPTOR = "unary-descriptor.pb"
const STREAM_DESCRIPTOR = "stream-descriptor.pb"
const CONFIG = "config.yaml"
const PROXY_PORT = 9000
const UPSTREAM_PORT = 8000
const EXTENSION_PORT = 10000

func TestUnaryProxy(t *testing.T) {
	withServers(t, UNARY_DESCRIPTOR, CONFIG, func(con *grpc.ClientConn, service *PingService, _ *CalcService) {
		client := ping.NewPingServiceClient(con)

		req := &ping.Ping{X: 100}
		res, err := client.Send(context.Background(), req)

		if err != nil {
			t.Error(err)
		}
		if *req != *service.lastRequest {
			t.Errorf("unexpected request: %v %v", req, service.lastRequest)
		}
		if *res != *service.lastResponse {
			t.Errorf("unexpected response: %v %v", res, service.lastResponse)
		}
	})
}

func TestClientSideStreamingProxy(t *testing.T) {
	withServers(t, STREAM_DESCRIPTOR, CONFIG, func(con *grpc.ClientConn, _ *PingService, service *CalcService) {
		client := calc.NewCalcServiceClient(con)

		if stream, err := client.AddAll(context.Background()); err != nil {
			t.Error(err)
		} else {
			requests := []*calc.Operand{}

			i := int64(0)
			sum := int64(0)
			for i < 5 {
				x := i * 2
				y := (i * 2) + 1
				req := &calc.Operand{X: x, Y: y}

				if err = stream.Send(req); err != nil {
					t.Error(err)
				}

				requests = append(requests, req)

				i = i + 1
				sum = sum + x + y
			}

			if res, err := stream.CloseAndRecv(); err != nil {
				t.Error(err)
			} else {
				if sum != res.Z {
					t.Errorf("%v != %v", sum, res)
				}

				if !reflect.DeepEqual(requests, service.lastAllRequests) {
					t.Errorf("unexpected request: %v %v", requests, service.lastAllRequests)
				}

				if *res != *service.lastAllResponse {
					t.Errorf("unexpected response: %v %v", res, service.lastAllResponse)
				}
			}
		}
	})
}

func TestServerSideStreamingProxy(t *testing.T) {
	withServers(t, STREAM_DESCRIPTOR, CONFIG, func(con *grpc.ClientConn, _ *PingService, service *CalcService) {
		client := calc.NewCalcServiceClient(con)

		req := &calc.OperandList{}
		i := int64(0)
		for i < 5 {
			x := i * 2
			y := (i * 2) + 1

			req.Operand = append(req.Operand, &calc.Operand{X: x, Y: y})

			i = i + 1
		}

		res := []*calc.Sum{}

		if stream, err := client.AddDeffered(context.Background(), req); err != nil {
			t.Error(err)
		} else {
			for {
				sum, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Error(err)
					break
				}
				res = append(res, sum)
			}

			if !reflect.DeepEqual(req, service.lastDefferedRequest) {
				t.Errorf("unexpected request: %v %v", req, service.lastDefferedRequest)
			}

			if !reflect.DeepEqual(res, service.lastDefferedResponses) {
				t.Errorf("unexpected response length: %v %v", res, service.lastDefferedResponses)
			}

			i := int64(0)
			for i < 5 {
				x := i * 2
				y := (i * 2) + 1

				if res[i].GetZ() != x + y {
					t.Errorf("unexpected response: %v %v %v", res[i], x, y)
				}

				i = i + 1
			}

		}
	})
}

func TestBidirectionalStreamingProxy(t *testing.T) {
	withServers(t, STREAM_DESCRIPTOR, CONFIG, func(con *grpc.ClientConn, _ *PingService, service *CalcService) {
		client := calc.NewCalcServiceClient(con)

		if stream, err := client.AddAsync(context.Background()); err != nil {
			t.Error(err)
		} else {
			requests := []*calc.Operand{}
			responses := []*calc.Sum{}

			i := int64(0)
			for i < 5 {
				x := i * 2
				y := (i * 2) + 1
				req := &calc.Operand{X: x, Y: y}

				if err := stream.Send(req); err != nil {
					t.Error(err)
				}

				res, err := stream.Recv()
				if err != nil {
					t.Error(err)
				}

				if res.GetZ() != x + y {
					t.Errorf("unexpected response: %v %v %v", res, x, y)
				}

				requests = append(requests, req)
				responses = append(responses, res)

				i = i + 1
			}

			if err := stream.CloseSend(); err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(requests, service.lastAsyncRequests) {
				t.Errorf("unexpected request: %v %v", requests, service.lastAsyncRequests)
			}

			if !reflect.DeepEqual(responses, service.lastAsyncResponses) {
				t.Errorf("unexpected response: %v %v", responses, service.lastAsyncResponses)
			}
		}
	})
}

