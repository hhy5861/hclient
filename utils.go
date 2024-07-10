package hclient

import (
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hhy5861/hclient/protocol"
	"strings"
	"time"
)

type (
	Protocol interface {
		Formatter() string
	}
)

const (
	K8SProtocol    = "K8S"
	K8SSProtocol   = "K8SS"
	defaultTimeout = 5000 * time.Millisecond
)

func (svc *Request) SetRemote(remote string) *Request {
	cfg, ok := svc.remotes[remote]
	if ok {
		svc.query.debug = cfg.Debug
		svc.query.timeout = defaultTimeout
		if cfg.Timeout > 0 {
			svc.query.timeout = cfg.Timeout * time.Millisecond
		}

		svc.query.targetUrl = svc.getProtocol(*cfg).Formatter()
	}

	return svc
}

func (svc *Request) SetPath(path string) *Request {
	path = strings.TrimRight(strings.TrimLeft(path, "/"), "/")

	targetUrl := fmt.Sprintf("%s/%s", strings.Trim(svc.query.targetUrl, "/"), path)
	svc.query.targetUrl = strings.Trim(targetUrl, "/")

	return svc
}

func (svc *Request) SetParam(content interface{}) *Request {
	svc.query.params = NewQuery(svc.client).Query(content)

	for k, v := range svc.query.params {
		svc.query.values.Set(k, v)
	}

	return svc
}

func (svc *Request) SetBody(content interface{}) *Request {
	svc.query.body = content

	return svc
}

func (svc *Request) SkipVerify() *resty.Client {
	if svc.query.skipVerify {
		svc.client.SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: svc.query.skipVerify,
		})
	}

	return svc.client
}

func (svc *Request) getProtocol(r Remote) Protocol {
	switch strings.ToUpper(r.Protocol) {
	case K8SProtocol, K8SSProtocol:
		return &protocol.K8S{
			Schema:    r.Protocol,
			Domain:    r.Domain,
			Service:   r.Service,
			Namespace: r.Namespace,
			Port:      r.Port,
		}
	default:
		return &protocol.Domain{
			Schema: r.Protocol,
			Domain: r.Domain,
			Port:   r.Port,
		}
	}
}
