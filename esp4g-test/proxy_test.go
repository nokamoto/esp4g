package esp4g_test

import (
	"testing"
	"google.golang.org/grpc"
	ping "github.com/nokamoto/esp4g/examples/ping/protobuf"
	"reflect"
	"github.com/nokamoto/esp4g/esp4g-extension/config"
)

func checkUnaryProxy(t *testing.T, cfg config.ExtensionConfig, apiKeys... string) {
	withServers(t, UNARY_DESCRIPTOR, cfg, apiKeys, func(con *grpc.ClientConn, service *PingService, _ *CalcService) {
		preflightPing(t, con)

		req := &ping.Ping{X: 100}
		res, err := callPing(con, req)

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

func checkClientSideStream(t *testing.T, cfg config.ExtensionConfig, apiKeys... string) {
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

		if !reflect.DeepEqual(req, service.lastAllRequests) {
			t.Errorf("unexpected request: %v %v", req, service.lastAllRequests)
		}

		if *res != *service.lastAllResponse {
			t.Errorf("unexpected response: %v %v", res, service.lastAllResponse)
		}
	})
}

func checkServerSideStream(t *testing.T, cfg config.ExtensionConfig, apiKeys... string) {
	withServers(t, STREAM_DESCRIPTOR, cfg, apiKeys, func(con *grpc.ClientConn, _ *PingService, service *CalcService) {
		preflightCalc(t, con)

		_, req, _, expected := makeTestCase()
		res, err := callCalcSStream(con, req)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(req, service.lastDefferedRequest) {
			t.Errorf("unexpected request: %v %v", req, service.lastDefferedRequest)
		}

		if !reflect.DeepEqual(res, service.lastDefferedResponses) {
			t.Errorf("unexpected response: %v %v", res, service.lastDefferedResponses)
		}

		if !reflect.DeepEqual(res, expected) {
			t.Errorf("unexpected response: %v %v", res, expected)
		}
	})
}

func checkBidirectionalStream(t *testing.T, cfg config.ExtensionConfig, apiKeys... string) {
	withServers(t, STREAM_DESCRIPTOR, cfg, apiKeys, func(con *grpc.ClientConn, _ *PingService, service *CalcService) {
		preflightCalc(t, con)

		req, _, _, expected := makeTestCase()
		res, err := callCalcBStream(con, req)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(req, service.lastAsyncRequests) {
			t.Errorf("unexpected request: %v %v", req, service.lastAsyncRequests)
		}

		if !reflect.DeepEqual(res, service.lastAsyncResponses) {
			t.Errorf("unexpected response: %v %v", res, service.lastAsyncResponses)
		}

		if !reflect.DeepEqual(res, expected) {
			t.Errorf("unexpected response: %v %v", res, expected)
		}
	})
}

func TextProxy(t *testing.T) {
	yaml, err := config.FromYamlFile("yaml/allow.yaml")
	if err != nil {
		t.Error(err)
	}
	checkUnaryProxy(t, yaml)
	checkClientSideStream(t, yaml)
	checkServerSideStream(t, yaml)
	checkBidirectionalStream(t, yaml)
}

