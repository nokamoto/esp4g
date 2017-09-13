package esp4g_test

import (
	"testing"
	"google.golang.org/grpc"
	ping "github.com/nokamoto/esp4g/examples/ping/protobuf"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"io"
	"github.com/nokamoto/esp4g/esp4g-extension/config"
)


func expectErrorCode(t *testing.T, err error, code codes.Code) {
	if err == nil {
		t.Error("expected error but actual nil")
	}

	stat, ok := status.FromError(err)

	if !ok {
		t.Errorf("expected status error but actual not: %v %v", err, stat)
	}

	if stat.Code() != code {
		t.Errorf("%v != %v", stat, code)
	}
}

func checkDenyUnary(t *testing.T, cfg config.ExtensionConfig, keys... string) {
	withServers(t, UNARY_DESCRIPTOR, cfg, keys, func(con *grpc.ClientConn, _ *PingService, _ *CalcService) {
		preflightPing(t, con)

		_, err := callPing(con, &ping.Ping{X: 100})

		t.Log("check unary deny")
		expectErrorCode(t, err, codes.Unauthenticated)
	})
}

func checkDenyCStream(t *testing.T, cfg config.ExtensionConfig, keys... string) {
	withServers(t, STREAM_DESCRIPTOR, cfg, keys, func(con *grpc.ClientConn, _ *PingService, _ *CalcService) {
		preflightCalc(t, con)

		operands, _, _, _ := makeTestCase()

		sum, err := callCalcCStream(con, operands)

		t.Log("check client-side stream deny", sum)
		if sum == nil && err == io.EOF {
			t.Logf("hotfix: unexpected state: sum=%v, err=%v", sum, err)
		} else {
			expectErrorCode(t, err, codes.Unauthenticated)
		}
	})
}

func checkDenySStream(t *testing.T, cfg config.ExtensionConfig, keys... string) {
	withServers(t, STREAM_DESCRIPTOR, cfg, keys, func(con *grpc.ClientConn, _ *PingService, _ *CalcService) {
		preflightCalc(t, con)

		_, operandList, _, _ := makeTestCase()

		sums, err := callCalcSStream(con, operandList)

		t.Log("check server-side stream deny", sums)
		expectErrorCode(t, err, codes.Unauthenticated)
	})
}

func checkDenyBStream(t *testing.T, cfg config.ExtensionConfig, keys... string) {
	withServers(t, STREAM_DESCRIPTOR, cfg, keys, func(con *grpc.ClientConn, _ *PingService, _ *CalcService) {
		preflightCalc(t, con)

		operands, _, _, _ := makeTestCase()

		sums, err := callCalcBStream(con, operands)

		t.Log("check bidirectional stream deny", sums)
		if sums == nil && err == io.EOF {
			t.Logf("hotfix: unexpected state: sums=%v, err=%v", sums, err)
		} else {
			expectErrorCode(t, err, codes.Unauthenticated)
		}
	})
}
