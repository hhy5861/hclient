package plugins

import "github.com/opentracing/opentracing-go"

type (
	Options struct {
		tracer opentracing.Tracer
	}

	Option func(opt *Options)
)

func Tracer(tracer opentracing.Tracer) Option {
	return func(opt *Options) {
		opt.tracer = tracer
	}
}
