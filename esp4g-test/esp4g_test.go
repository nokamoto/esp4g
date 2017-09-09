package esp4g_test

import (
	"testing"
	"google.golang.org/grpc"
	ping "github.com/nokamoto/esp4g/examples/ping/protobuf"
	calc "github.com/nokamoto/esp4g/examples/calc/protobuf"
	"golang.org/x/net/context"
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
		if len(service.requests) != 1 || *req != service.requests[0] {
			t.Errorf("unexpected request: %v %v", req, service.requests)
		}
		if len(service.responses) != 1 || *res != service.responses[0] {
			t.Errorf("unexpected response: %v %v", res, service.responses)
		}
	})
}

func TestClientSideStreamingProxy(t *testing.T) {
	withServers(t, STREAM_DESCRIPTOR, CONFIG, func(con *grpc.ClientConn, _ *PingService, service *CalcService) {
		client := calc.NewCalcServiceClient(con)

		if stream, err := client.AddAll(context.Background()); err != nil {
			t.Error(err)
		} else {
			i := int64(0)
			sum := int64(0)
			for i < 5 {
				x := i * 2
				y := (i * 2) + 1
				req := calc.Operand{X: x, Y: y}

				if err = stream.Send(&req); err != nil {
					t.Error(err)
				}

				i = i + 1
				sum = sum + x + y
			}

			if res, err := stream.CloseAndRecv(); err != nil {
				t.Error(err)
			} else {
				if sum != res.Z {
					t.Errorf("%v != %v", sum, res)
				}

				i := int64(0)
				for i < 5 {
					x := i * 2
					y := (i * 2) + 1

					if len(service.allRequests) != 1 ||
						len(service.allRequests[0]) != 5 ||
						service.allRequests[0][i].GetX() != x ||
						service.allRequests[0][i].GetY() != y {
						t.Errorf("unexpected request: %v %v %v", service.allRequests[0][i], x, y)
					}

					i = i + 1
				}

				if len(service.allResponses) != 1 || *res != service.allResponses[0] {
					t.Errorf("unexpected response: %v %v", res, service.allResponses[0])
				}
			}
		}
	})
}


