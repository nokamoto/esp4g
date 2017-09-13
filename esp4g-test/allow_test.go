package esp4g_test

import (
	"reflect"
	"testing"
	"github.com/nokamoto/esp4g/esp4g-extension/config"
	ping "github.com/nokamoto/esp4g/examples/ping/protobuf"
	"google.golang.org/grpc"
)

func checkUnary(t *testing.T, cfg config.ExtensionConfig, apiKeys... string) {
	withServers(t, UNARY_DESCRIPTOR, cfg, apiKeys, func(con *grpc.ClientConn, service *PingService, _ *CalcService) {
		preflightPing(t, con)

		req := &ping.Ping{X: 100}
		res, err := callPing(con, req)

		if err != nil {
			t.Error(err)
		}
		if *req != *service.lastRequest() {
			t.Errorf("unexpected request: %v %v", req, service.lastRequest())
		}
		if *res != *service.lastResponse() {
			t.Errorf("unexpected response: %v %v", res, service.lastResponse())
		}
	})
}

func checkCStream(t *testing.T, cfg config.ExtensionConfig, apiKeys... string) {
	withServers(t, STREAM_DESCRIPTOR, cfg, apiKeys, func(con *grpc.ClientConn, _ *PingService, service *CalcService) {
		preflightCalc(t, con)

		req, _, expected, _ := makeTestCase()
		res, err := callCalcCStream(con, req)

		if err != nil {
			t.Error(err)
		}

		if *res != *expected {
			t.Errorf("%v != %v", res, expected)
		}

		if !reflect.DeepEqual(req, service.lastAllRequests()) {
			t.Errorf("unexpected request: %v %v", req, service.lastAllRequests())
		}

		if *res != *service.lastAllResponse() {
			t.Errorf("unexpected response: %v %v", res, service.lastAllResponse())
		}
	})
}

func checkSStream(t *testing.T, cfg config.ExtensionConfig, apiKeys... string) {
	withServers(t, STREAM_DESCRIPTOR, cfg, apiKeys, func(con *grpc.ClientConn, _ *PingService, service *CalcService) {
		preflightCalc(t, con)

		_, req, _, expected := makeTestCase()
		res, err := callCalcSStream(con, req)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(req, service.lastDefferedRequest()) {
			t.Errorf("unexpected request: %v %v", req, service.lastDefferedRequest())
		}

		if !reflect.DeepEqual(res, service.lastDefferedResponses()) {
			t.Errorf("unexpected response: %v %v", res, service.lastDefferedResponses())
		}

		if !reflect.DeepEqual(res, expected) {
			t.Errorf("unexpected response: %v %v", res, expected)
		}
	})
}

func checkBStream(t *testing.T, cfg config.ExtensionConfig, apiKeys... string) {
	withServers(t, STREAM_DESCRIPTOR, cfg, apiKeys, func(con *grpc.ClientConn, _ *PingService, service *CalcService) {
		preflightCalc(t, con)

		req, _, _, expected := makeTestCase()
		res, err := callCalcBStream(con, req)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(req, service.lastAsyncRequests()) {
			t.Errorf("unexpected request: %v %v", req, service.lastAsyncRequests())
		}

		if !reflect.DeepEqual(res, service.lastAsyncResponses()) {
			t.Errorf("unexpected response: %v %v", res, service.lastAsyncResponses())
		}

		if !reflect.DeepEqual(res, expected) {
			t.Errorf("unexpected response: %v %v", res, expected)
		}
	})
}
