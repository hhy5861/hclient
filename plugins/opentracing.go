package plugins

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hhy5861/nethttp"
	"net/url"
)

const (
	defaultComponentName = "net/http"
)

func (opt *Options) BeforeRequest(c *resty.Client, req *resty.Request) error {
	ul, err := url.Parse(req.URL)
	if err != nil {
		return err
	}

	req.SetContext(nethttp.TraceWithContext(
		opt.tracer, req.Context(),
		nethttp.OperationName(fmt.Sprintf("RESTY %s: %s", req.Method, ul.Path)),
		nethttp.ComponentName(defaultComponentName),
		nethttp.ClientTrace(true),
		nethttp.InjectSpanContext(true),
	))

	return nil
}

func (opt *Options) AfterResponse(c *resty.Client, resp *resty.Response) error {
	tracer := nethttp.TracerFromRequest(resp.Request.RawRequest)
	if tracer != nil {
		defer tracer.Finish()
	}

	return nil
}
