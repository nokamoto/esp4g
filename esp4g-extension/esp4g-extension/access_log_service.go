package extension

import (
	"golang.org/x/net/context"
	proto "github.com/nokamoto/esp4g/protobuf"
	"github.com/golang/protobuf/ptypes/empty"
	"time"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/nokamoto/esp4g/esp4g-utils"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)

type accessLogService struct {
	logs Logs

	requestSizeHistogram *prometheus.HistogramVec

	responseSizeHistogram *prometheus.HistogramVec

	responseTimeHistogram *prometheus.HistogramVec
}

func convert(d *duration.Duration) time.Duration {
	if d == nil {
		return time.Duration(-1)
	}
	return time.Duration(d.Seconds) * time.Second + time.Duration(d.Nanos)
}

func observer(unary *proto.UnaryAccessLog, vec *prometheus.HistogramVec) prometheus.Observer {
	return vec.WithLabelValues(unary.GetMethod(), unary.GetStatus())
}

func observerStream(stream *proto.StreamAccessLog, vec *prometheus.HistogramVec) prometheus.Observer {
	return vec.WithLabelValues(stream.GetMethod(), stream.GetStatus())
}

func (a *accessLogService)UnaryAccess(_ context.Context, unary *proto.UnaryAccessLog) (*empty.Empty, error) {
	rt := convert(unary.GetResponseTime())
	if a.logs.Logging {
		utils.Logger.Infow("unary",
			"method", unary.GetMethod(),
			"status", unary.GetStatus(),
			"response_time", rt,
			"request_size", unary.GetRequestSize(),
			"response_size", unary.GetResponseSize(),
		)
	}
	if a.responseTimeHistogram != nil {
		observer(unary, a.responseTimeHistogram).Observe(rt.Seconds())
	}
	if a.requestSizeHistogram != nil {
		observer(unary, a.requestSizeHistogram).Observe(float64(unary.GetRequestSize()))
	}
	if a.responseSizeHistogram != nil {
		observer(unary, a.responseSizeHistogram).Observe(float64(unary.GetResponseSize()))
	}
	return &empty.Empty{}, nil
}

func (a *accessLogService)StreamAccess(_ context.Context, stream *proto.StreamAccessLog) (*empty.Empty, error) {
	rt := convert(stream.GetResponseTime())
	if a.logs.Logging {
		utils.Logger.Infow("stream",
			"method", stream.GetMethod(),
			"status", stream.GetStatus(),
			"response_time", rt,
		)
	}
	if a.responseTimeHistogram != nil {
		observerStream(stream, a.responseTimeHistogram).Observe(rt.Seconds())
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

func newAccessLogService(config Config, _ *descriptor.FileDescriptorSet) *accessLogService {
	c := config.Logs.Prometheus
	var requestSizeHistogram *prometheus.HistogramVec
	var responseSizeHistogram *prometheus.HistogramVec
	var responseTimeHistogram *prometheus.HistogramVec
	if c.Enable {
		lbs := labelNames()
		requestSizeHistogram = register(c.RequestSizeHistogram, lbs)
		responseSizeHistogram = register(c.ResponseSizeHistogram, lbs)
		responseTimeHistogram = register(c.ResponseTimeHistogram, lbs)

		go func() {
			http.Handle("/metrics", promhttp.Handler())
			err := http.ListenAndServe(c.Address, nil)
			utils.Logger.Fatalw("stop prometheus exporter", "err", err)
		}()

		utils.Logger.Infow("listen prometheus exporter", "address", c.Address)
	}
	return &accessLogService{
		logs: config.Logs,
		requestSizeHistogram: requestSizeHistogram,
		responseSizeHistogram: responseSizeHistogram,
		responseTimeHistogram: responseTimeHistogram,
	}
}
