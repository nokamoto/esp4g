package extension

import (
	"golang.org/x/net/context"
	proto "github.com/nokamoto/esp4g/protobuf"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nokamoto/esp4g/esp4g-utils"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

type AccessLogService struct {
	requestBytes *prometheus.HistogramVec

	responseBytes *prometheus.HistogramVec

	responseSeconds *prometheus.HistogramVec

	logger *zap.SugaredLogger
}

func authority(access *proto.GrpcAccess) string {
	return strings.Join(access.Authority, ",")
}

func userAgent(access *proto.GrpcAccess) string {
	return strings.Join(access.UserAgent, ",")
}

func observer(access *proto.GrpcAccess, vec *prometheus.HistogramVec) prometheus.Observer {
	return vec.WithLabelValues(access.GetMethod(), access.GetStatus())
}

func (a *AccessLogService)grpcAccess(access *proto.GrpcAccess, keysAndValues... interface{}) {
	rt := utils.ConvertProtoDuration(access.GetResponseTime())

	if a.logger != nil {
		args := []interface{}{
			"method", access.GetMethod(),
			"status", access.GetStatus(),
			"response_seconds", rt.Seconds(),
			"authority", authority(access),
			"user_agent", userAgent(access),
		}

		args = append(args, keysAndValues...)

		a.logger.Infow("", args...)
	}

	if a.responseSeconds != nil {
		observer(access, a.responseSeconds).Observe(rt.Seconds())
	}
}

func (a *AccessLogService)UnaryAccess(_ context.Context, unary *proto.UnaryAccessLog) (*empty.Empty, error) {
	a.grpcAccess(unary.GetAccess(),
		"request_bytes", unary.GetRequestBytesSize(),
		"response_bytes", unary.GetResponseBytesSize(),
	)
	if a.requestBytes != nil {
		observer(unary.GetAccess(), a.requestBytes).Observe(float64(unary.GetRequestBytesSize()))
	}
	if a.responseBytes != nil {
		observer(unary.GetAccess(), a.responseBytes).Observe(float64(unary.GetResponseBytesSize()))
	}
	return &empty.Empty{}, nil
}

func (a *AccessLogService)StreamAccess(_ context.Context, stream *proto.StreamAccessLog) (*empty.Empty, error) {
	a.grpcAccess(stream.GetAccess())
	return &empty.Empty{}, nil
}

func labelNames() []string {
	return []string{"method", "status"}
}

func register(h *Histogram, lbs []string) *prometheus.HistogramVec {
	if h != nil {
		hv := prometheus.NewHistogramVec(h.histogramOpts(), lbs)
		prometheus.MustRegister(hv)
		utils.Logger.Infow("register prometheus histogram", "histogram", h)
		return hv
	}
	return nil
}

func NewAccessLogService(config Config, _ *descriptor.FileDescriptorSet) *AccessLogService {
	var sugar *zap.SugaredLogger

	if s, err := sugaredLogger(config); err != nil {
		utils.Logger.Fatalw("failed to create zap logger", "err", err)
	} else {
		sugar = s
	}

	c := config.Logs.Prometheus
	var requestBytes *prometheus.HistogramVec
	var responseBytes *prometheus.HistogramVec
	var responseSeconds *prometheus.HistogramVec
	if c.Port != nil {
		lbs := labelNames()
		requestBytes = register(c.Histograms.RequestBytes, lbs)
		responseBytes = register(c.Histograms.ResponseBytes, lbs)
		responseSeconds = register(c.Histograms.ResponseSeconds, lbs)

		go func() {
			http.Handle("/metrics", promhttp.Handler())
			err := http.ListenAndServe(fmt.Sprintf(":%d", *c.Port), nil)
			utils.Logger.Fatalw("stop prometheus exporter", "err", err)
		}()

		utils.Logger.Infow("listen prometheus exporter", "port", c.Port)
	}

	return &AccessLogService{
		requestBytes:    requestBytes,
		responseBytes:   responseBytes,
		responseSeconds: responseSeconds,
		logger:          sugar,
	}
}
