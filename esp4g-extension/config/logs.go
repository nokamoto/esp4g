package config

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type Histogram struct {
	Namespace string `yaml:"namespace"`

	Subsystem string `yaml:"subsystem"`

	Name string `yaml:"name"`

	Help string `yaml:"help"`

	Buckets []float64 `yaml:"buckets"`
}

func (h Histogram)HistogramOpts() prometheus.HistogramOpts {
	return prometheus.HistogramOpts {
		Namespace: h.Namespace,
		Subsystem: h.Subsystem,
		Name: h.Name,
		Help: h.Help,
		Buckets: h.Buckets,
	}
}

type Histograms struct {
	ResponseSeconds *Histogram `yaml:"response_seconds"`

	RequestBytes *Histogram `yaml:"request_bytes"`

	ResponseBytes *Histogram `yaml:"response_bytes"`
}

type Prometheus struct {
	Port *int `yaml:"port"`

	Histograms Histograms `yaml:"histograms"`
}

type Logs struct {
	Zap *zap.Config `yaml:"zap"`

	Prometheus Prometheus `yaml:"prometheus"`
}
