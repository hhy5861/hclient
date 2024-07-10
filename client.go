package hclient

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"net/url"
	"time"
)

type (
	Client struct {
		request *request
	}

	ConfigCache struct {
		remotes map[string]*Remote
		resp    IResponse
		trace   opentracing.Tracer
	}

	Remote struct {
		Protocol  string        `json:"protocol" yaml:"protocol"`
		Service   string        `json:"service" yaml:"service"`
		Namespace string        `json:"namespace" yaml:"namespace"`
		Domain    string        `json:"domain" yaml:"domain"`
		Port      int           `json:"port" yaml:"port"`
		Debug     bool          `json:"debug" yaml:"debug"`
		Timeout   time.Duration `json:"timeout" yaml:"timeout"`
	}

	Option func(c *ConfigCache)
)

var (
	remotes map[string]*Remote
)

// NewClient init client
func NewClient(opts ...Option) *ConfigCache {
	cfg := &ConfigCache{
		resp:  NewResponse(),
		trace: opentracing.GlobalTracer(),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	remotes = cfg.remotes
	return cfg
}

// GlobalClient get client
func GlobalClient() *Client {
	return &Client{
		request: NewRequest(context.TODO(), remotes),
	}
}

// GlobalWithCtx get client
func GlobalWithCtx(ctx context.Context) *Client {
	return &Client{
		request: NewRequest(ctx, remotes),
	}
}

// Get request get
func (c *Client) Get(
	remote,
	path string,
	content interface{}) IResponse {

	return c.request.SetRemote(remote).SetPath(path).SetParam(content).Get()
}

// Post request post
func (c *Client) Post(
	remote,
	path string,
	content, param interface{}) IResponse {

	return c.request.SetRemote(remote).SetPath(path).SetBody(content).SetParam(param).Post()
}

// PostUrlEncode request post
func (c *Client) PostUrlEncode(
	remote,
	path string,
	content, param interface{}) IResponse {

	return c.request.SetRemote(remote).SetPath(path).SetBody(content).SetParam(param).PostUrlEncode()
}

// Put request put
func (c *Client) Put(
	remote,
	path string,
	content, param interface{}) IResponse {

	return c.request.SetRemote(remote).SetPath(path).SetBody(content).SetParam(param).Put()
}

// PostJson request post json
func (c *Client) PostJson(
	remote,
	path string,
	content, param interface{}) IResponse {

	return c.request.SetRemote(remote).SetPath(path).SetBody(content).SetParam(param).PostJson()
}

// Delete request delete
func (c *Client) Delete(
	remote,
	path string,
	content, param interface{}) IResponse {

	return c.request.SetRemote(remote).SetPath(path).SetBody(content).SetParam(param).Delete()
}

// SetHeader set request header data params
func (c *Client) SetHeader(data map[string]string) *Client {
	c.request.query.header = data

	return c
}

// SetTimeOut set request time out params
func (c *Client) SetTimeOut(timeOut time.Duration) *Client {
	c.request.query.timeout = timeOut

	return c
}

// AddParams add request time out params
func (c *Client) AddParams(key, value string) {
	c.request.query.Add(key, value)
}

// SkipVerify is skip insecure verify
func (c *Client) SkipVerify(skipVerify bool) *Client {
	c.request.query.skipVerify = skipVerify

	return c
}

func (c *Client) EnabledDebug() *Client {
	c.request.query.debug = true

	return c
}

func (c *Client) SetQueryParamsFromValues(val url.Values) *Client {
	c.request.query.values = val

	return c
}
