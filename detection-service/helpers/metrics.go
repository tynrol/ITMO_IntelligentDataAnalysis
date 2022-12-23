package metrics

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/mock"
)

var defaultLabels = []string{
	"pod_name",
	"pod_node",
	"pod_namespace",
}

var defaultValues = []string{
	os.Getenv("LOGGER_POD_NAME"),
	os.Getenv("LOGGER_POD_NODE"),
	os.Getenv("LOGGER_POD_NAMESPACE"),
}

type PrometheusMetrics interface {
	NewCounterVec(name, help string, labels ...string) CounterVec
	NewGaugeVec(name, help string, labels ...string) GaugeVec
}

type Prometheus struct {
	namespace string
	subsystem string
}

func NewPrometheus(
	namespace string,
	subsystem string,
) *Prometheus {
	return &Prometheus{
		namespace: namespace,
		subsystem: subsystem,
	}
}

func (p *Prometheus) NewCounterVec(name, help string, labels ...string) CounterVec {
	out := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: p.namespace,
		Subsystem: p.subsystem,
		Name:      name,
		Help:      help,
	}, append(defaultLabels, labels...))

	prometheus.MustRegister(out)

	return &counterVec{counterVec: out}
}

func (p *Prometheus) NewGaugeVec(name, help string, labels ...string) GaugeVec {
	out := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: p.namespace,
		Subsystem: p.subsystem,
		Name:      name,
		Help:      help,
	}, append(defaultLabels, labels...))

	prometheus.MustRegister(out)

	return &gaugeVec{gaugeVec: out}
}

type CounterVec interface {
	WithLabelValues(labels ...string) prometheus.Counter
}

type counterVec struct {
	counterVec *prometheus.CounterVec
}

func (c *counterVec) WithLabelValues(labels ...string) prometheus.Counter {
	return c.counterVec.WithLabelValues(append(defaultValues, labels...)...)
}

type GaugeVec interface {
	WithLabelValues(labels ...string) prometheus.Gauge
}

type gaugeVec struct {
	gaugeVec *prometheus.GaugeVec
}

func (c *gaugeVec) WithLabelValues(labels ...string) prometheus.Gauge {
	return c.gaugeVec.WithLabelValues(append(defaultValues, labels...)...)
}

type CounterVecMock struct {
}

func (c *CounterVecMock) WithLabelValues(labels ...string) prometheus.Counter {
	var mock prometheus.Counter

	return mock
}

type GaugeVecMock struct {
}

func (c *GaugeVecMock) WithLabelValues(labels ...string) prometheus.Gauge {
	var mock prometheus.Gauge

	return mock
}

type PrometheusMetricsMock struct {
	mock.Mock
}

func (m *PrometheusMetricsMock) NewCounterVec(name, help string, labels ...string) CounterVec {
	return &CounterVecMock{}
}

func (m *PrometheusMetricsMock) NewGaugeVec(name, help string, labels ...string) GaugeVec {
	return &GaugeVecMock{}
}
