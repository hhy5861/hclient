package plugins

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
)

var (
	clientReqCnt = &Metric{
		ID:          "clientReqCnt",
		Name:        "client_requests_total",
		Description: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		Type:        "counter_vec",
		Args:        []string{"method", "url"}}

	clientReqSz = &Metric{
		ID:          "clientReqSz",
		Name:        "client_request_size_bytes",
		Description: "The HTTP request sizes in bytes.",
		Type:        "summary_vec",
		Args:        []string{"method", "url"}}

	clientReqDur = &Metric{
		ID:          "clientReqDur",
		Name:        "client_request_duration_seconds",
		Description: "The HTTP request latencies in seconds.",
		Type:        "histogram_vec",
		Args:        []string{"code", "method", "url"}}

	clientResSz = &Metric{
		ID:          "clientResSz",
		Name:        "client_response_size_bytes",
		Description: "The HTTP response sizes in bytes.",
		Type:        "summary_vec",
		Args:        []string{"code", "method", "url"}}

	standardClientMetrics = []*Metric{
		clientReqCnt,
		clientReqSz,
		clientReqDur,
		clientResSz,
	}
)

type Metric struct {
	MetricCollector prometheus.Collector
	ID              string
	Name            string
	Description     string
	Type            string
	Args            []string
}

type Prometheus struct {
	clientReqCnt *prometheus.CounterVec
	clientReqSz  *prometheus.SummaryVec
	clientResSz  *prometheus.SummaryVec
	clientReqDur *prometheus.HistogramVec
	MetricsList  []*Metric
}

func NewPrometheus(subsystem string, customMetricsList ...[]*Metric) *Prometheus {
	var metricsList []*Metric

	if len(customMetricsList) > 1 {
		panic("Too many args. NewPrometheus( string, <optional []*Metric> ).")
	} else if len(customMetricsList) == 1 {
		metricsList = customMetricsList[0]
	}

	for _, metric := range standardClientMetrics {
		metricsList = append(metricsList, metric)
	}

	promMetrics := &Prometheus{
		MetricsList: metricsList,
	}

	promMetrics.registerMetrics(subsystem)
	return promMetrics
}

func NewMetric(m *Metric, subsystem string) prometheus.Collector {
	var metric prometheus.Collector

	switch m.Type {
	case "counter_vec":
		metric = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "counter":
		metric = prometheus.NewCounter(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "gauge_vec":
		metric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "gauge":
		metric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "histogram_vec":
		metric = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "histogram":
		metric = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "summary_vec":
		metric = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "summary":
		metric = prometheus.NewSummary(
			prometheus.SummaryOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	}

	return metric
}

func (p *Prometheus) registerMetrics(subsystem string) {
	for _, metricDef := range p.MetricsList {
		metric := NewMetric(metricDef, subsystem)
		if err := prometheus.Register(metric); err != nil {
			log.Println("[Prometheus registered]", err)
		}

		switch metricDef {
		case clientReqCnt:
			p.clientReqCnt = metric.(*prometheus.CounterVec)
		case clientReqDur:
			p.clientReqDur = metric.(*prometheus.HistogramVec)
		case clientResSz:
			p.clientResSz = metric.(*prometheus.SummaryVec)
		case clientReqSz:
			p.clientReqSz = metric.(*prometheus.SummaryVec)
		}

		metricDef.MetricCollector = metric
	}
}

func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}

	s += len(r.Host)
	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}

	return s
}
