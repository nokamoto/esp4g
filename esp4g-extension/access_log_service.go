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
	"fmt"
)

type AccessLogService struct {
	logs Logs

	requestBytesHistogram *prometheus.HistogramVec

	responseBytesHistogram *prometheus.HistogramVec

	responseSecondsHistogram *prometheus.HistogramVec
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

func (a *AccessLogService)UnaryAccess(_ context.Context, unary *proto.UnaryAccessLog) (*empty.Empty, error) {
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
	if a.responseSecondsHistogram != nil {
		observer(unary, a.responseSecondsHistogram).Observe(rt.Seconds())
	}
	if a.requestBytesHistogram != nil {
		observer(unary, a.requestBytesHistogram).Observe(float64(unary.GetRequestSize()))
	}
	if a.responseBytesHistogram != nil {
		observer(unary, a.responseBytesHistogram).Observe(float64(unary.GetResponseSize()))
	}
	return &empty.Empty{}, nil
}

func (a *AccessLogService)StreamAccess(_ context.Context, stream *proto.StreamAccessLog) (*empty.Empty, error) {
	rt := convert(stream.GetResponseTime())
	if a.logs.Logging {
		utils.Logger.Infow("stream",
			"method", stream.GetMethod(),
			"status", stream.GetStatus(),
			"response_time", rt,
		)
	}
	if a.responseSecondsHistogram != nil {
		observerStream(stream, a.responseSecondsHistogram).Observe(rt.Seconds())
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
	c := config.Logs.Prometheus
	var requestBytesHistogram *prometheus.HistogramVec
	var responseBytesHistogram *prometheus.HistogramVec
	var responseSecondsHistogram *prometheus.HistogramVec
	if c.Port != nil {
		lbs := labelNames()
		requestBytesHistogram = register(c.Histograms.RequestBytes, lbs)
		responseBytesHistogram = register(c.Histograms.ResponseBytes, lbs)
		responseSecondsHistogram = register(c.Histograms.ResponseSeconds, lbs)

		go func() {
			http.Handle("/metrics", promhttp.Handler())
			err := http.ListenAndServe(fmt.Sprintf(":%d", *c.Port), nil)
			utils.Logger.Fatalw("stop prometheus exporter", "err", err)
		}()

		utils.Logger.Infow("listen prometheus exporter", "port", c.Port)
	}
	return &AccessLogService{
		logs:                     config.Logs,
		requestBytesHistogram:    requestBytesHistogram,
		responseBytesHistogram:   responseBytesHistogram,
		responseSecondsHistogram: responseSecondsHistogram,
	}
}
