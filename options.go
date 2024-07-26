package hclient

import (
	"github.com/opentracing/opentracing-go"
)

// SetRemotes set client remote service list
func SetRemotes(remote map[string]*Remote) Option {
	return func(c *ConfigCache) {
		for k, v := range remote {
			if v.Service == "" {
				v.Service = k
			}
		}

		c.remotes = remote
	}
}

// SetResponse set call client response result
func SetResponse(r IResponse) Option {
	return func(c *ConfigCache) {
		c.resp = r
	}
}

// Opentracing set open trace client reporter
func Opentracing(t opentracing.Tracer) Option {
	return func(c *ConfigCache) {
		c.trace = t
	}
}
