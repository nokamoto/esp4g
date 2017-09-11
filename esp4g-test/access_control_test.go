package esp4g_test

import (
	"testing"
	"google.golang.org/grpc"
	ping "github.com/nokamoto/esp4g/examples/ping/protobuf"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

func expectErrorCode(t *testing.T, err error, code codes.Code) {
	if err == nil {
		t.Error("expected error but actual nil")
	}

	stat, ok := status.FromError(err)

	if !ok {
		t.Errorf("expected status error but actual not: %v", err)
	}

	if stat.Code() != code {
		t.Errorf("%v != %v", stat, code)
	}
}

func TestDenyUnregisteredCalls(t *testing.T) {
	config := "access_control-deny.yaml"

	withServers(t, UNARY_DESCRIPTOR, config, func(con *grpc.ClientConn, _ *PingService, _ *CalcService) {
		preflightPing(t, con)

		_, err := callPing(con, &ping.Ping{X: 100})

		expectErrorCode(t, err, codes.Unauthenticated)
	})

	withServers(t, STREAM_DESCRIPTOR, config, func(con *grpc.ClientConn, _ *PingService, _ *CalcService) {
		preflightCalc(t, con)

		operands, operandList, _, _ := makeTestCase()

		_, err := callCalcCStream(con, operands)

		expectErrorCode(t, err, codes.Unauthenticated)

		_, err = callCalcSStream(con, operandList)

		expectErrorCode(t, err, codes.Unauthenticated)

		_, err = callCalcBStream(con, operands)

		expectErrorCode(t, err, codes.Unauthenticated)
	})
}
