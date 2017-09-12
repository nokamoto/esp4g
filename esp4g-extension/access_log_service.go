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
)

type AccessLogService struct {
	requestBytes *prometheus.HistogramVec

	responseBytes *prometheus.HistogramVec

	responseSeconds *prometheus.HistogramVec

	logger *zap.SugaredLogger
}

func observer(unary *proto.UnaryAccessLog, vec *prometheus.HistogramVec) prometheus.Observer {
	return vec.WithLabelValues(unary.GetMethod(), unary.GetStatus())
}

func observerStream(stream *proto.StreamAccessLog, vec *prometheus.HistogramVec) prometheus.Observer {
	return vec.WithLabelValues(stream.GetMethod(), stream.GetStatus())
}

func (a *AccessLogService)UnaryAccess(_ context.Context, unary *proto.UnaryAccessLog) (*empty.Empty, error) {
	rt := utils.ConvertProtoDuration(unary.GetResponseTime())
	if a.logger != nil {
		a.logger.Infow("",
			"method", unary.GetMethod(),
			"status", unary.GetStatus(),
			"response_seconds", rt.Seconds(),
			"request_bytes", unary.GetRequestSize(),
			"response_bytes", unary.GetResponseSize(),
		)
	}
	if a.responseSeconds != nil {
		observer(unary, a.responseSeconds).Observe(rt.Seconds())
	}
	if a.requestBytes != nil {
		observer(unary, a.requestBytes).Observe(float64(unary.GetRequestSize()))
	}
	if a.responseBytes != nil {
		observer(unary, a.responseBytes).Observe(float64(unary.GetResponseSize()))
	}
	return &empty.Empty{}, nil
}

func (a *AccessLogService)StreamAccess(_ context.Context, stream *proto.StreamAccessLog) (*empty.Empty, error) {
	rt := utils.ConvertProtoDuration(stream.GetResponseTime())
	if a.logger != nil {
		a.logger.Infow("",
			"method", stream.GetMethod(),
			"status", stream.GetStatus(),
			"response_seconds", rt.Seconds(),
		)
	}
	if a.responseSeconds != nil {
		observerStream(stream, a.responseSeconds).Observe(rt.Seconds())
	}
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
