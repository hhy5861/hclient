package plugins

import (
	"github.com/go-resty/resty/v2"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"strconv"
)

type (
	Hook struct {
		option *Options
		metric *Prometheus
	}
)

const (
	operationName = "resty"
)

var (
	metric = NewPrometheus(operationName)
)

func NewHook(opts ...Option) *Hook {
	opt := &Options{
		tracer: opentracing.GlobalTracer(),
	}

	for _, fn := range opts {
		fn(opt)
	}

	return &Hook{
		option: opt,
		metric: metric,
	}
}

func (h *Hook) BeforeRequest(c *resty.Client, req *resty.Request) error {
	return h.option.BeforeRequest(c, req)
}

func (h *Hook) PreRequestHook(c *resty.Client, req *http.Request) error {
	h.metric.clientReqCnt.WithLabelValues(req.Method, req.Method).Inc()
	h.metric.clientReqSz.WithLabelValues(req.Method, req.URL.Path).Observe(float64(computeApproximateRequestSize(req)))

	return nil
}

func (h *Hook) AfterResponse(c *resty.Client, resp *resty.Response) error {
	code := strconv.Itoa(resp.StatusCode())
	path := resp.Request.RawRequest.URL.Path

	h.metric.clientReqDur.WithLabelValues(code, resp.Request.Method, path).Observe(resp.Time().Seconds())
	h.metric.clientResSz.WithLabelValues(code, resp.Request.Method, path).Observe(float64(resp.Size()))

	return h.option.AfterResponse(c, resp)
}
