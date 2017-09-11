package extension

import "github.com/prometheus/client_golang/prometheus"

type Histogram struct {
	Namespace string `yaml:"namespace"`

	Subsystem string `yaml:"subsystem"`

	Name string `yaml:"name"`

	Help string `yaml:"help"`

	Buckets []float64 `yaml:"buckets"`
}

func (h Histogram)histogramOpts() prometheus.HistogramOpts {
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
	Logging bool `yaml:"logging"`

	Prometheus Prometheus `yaml:"prometheus"`
}

type Rule struct {
	Selector string `yaml:"selector"`

	AllowUnregisteredCalls bool `yaml:"allow_unregistered_calls"`

	RegisteredApiKey []string `yaml:"registered_api_keys"`
}

type Usage struct {
	Rules []Rule `yaml:"rules"`
}

type Config struct {
	Logs Logs `yaml:"logs"`

	Usage Usage `yaml:"usage"`
}
