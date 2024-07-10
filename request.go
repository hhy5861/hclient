package hclient

import (
	"github.com/hhy5861/hclient/plugins"
	"github.com/hhy5861/nethttp"
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
)

type (
	request struct {
		remotes map[string]*Remote
		client  *resty.Client
		query   *Query
		tracer  *nethttp.Tracer
	}
)

var (
	defaultJson = NewDefaultResp().ToByte()
	cli         = &http.Client{
		Transport: &nethttp.Transport{},
	}
)

func NewRequest(ctx context.Context, remotes map[string]*Remote) *request {
	hook := plugins.NewHook()

	return &request{
		query: &Query{
			ctx:      ctx,
			header:   make(map[string]string),
			params:   make(map[string]string),
			response: NewResponse().SetBody(defaultJson),
			values:   url.Values{},
		},
		client:  resty.NewWithClient(cli).OnBeforeRequest(hook.BeforeRequest).SetPreRequestHook(hook.PreRequestHook).OnAfterResponse(hook.AfterResponse),
		remotes: remotes,
	}
}

func (svc *request) Get() IResponse {
	res, err := svc.SkipVerify().
		SetDebug(svc.query.debug).
		SetTimeout(svc.query.timeout).
		R().SetContext(svc.query.ctx).SetHeaders(svc.query.header).SetQueryParamsFromValues(svc.query.values).Get(svc.query.targetUrl)
	if err == nil && res != nil {
		return svc.query.response.SetBody(res.Body(), res.StatusCode())
	}

	return svc.query.response.SetError(err)
}

func (svc *request) Post() IResponse {
	res, err := svc.SkipVerify().
		SetDebug(svc.query.debug).
		SetTimeout(svc.query.timeout).
		R().SetContext(svc.query.ctx).SetHeaders(svc.query.header).SetBody(svc.query.body).SetQueryParamsFromValues(svc.query.values).Post(svc.query.targetUrl)
	if err == nil && res != nil {
		return svc.query.response.
			SetBody(res.Body(), res.StatusCode())
	}

	return svc.query.response.SetError(err)
}

func (svc *request) PostUrlEncode() IResponse {
	res, err := svc.SkipVerify().
		SetDebug(svc.query.debug).
		SetTimeout(svc.query.timeout).
		R().SetContext(svc.query.ctx).SetHeaders(svc.query.header).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").SetBody(svc.query.body).SetQueryParamsFromValues(svc.query.values).Post(svc.query.targetUrl)
	if err == nil && res != nil {
		return svc.query.response.SetBody(res.Body(), res.StatusCode())
	}

	return svc.query.response.SetError(err)
}

func (svc *request) PostJson() IResponse {
	body, errMsg := json.Marshal(svc.query.body)
	if errMsg != nil {
		return svc.query.response.SetBody(defaultJson, 406)
	}

	res, err := svc.SkipVerify().
		SetDebug(svc.query.debug).
		SetTimeout(svc.query.timeout).
		R().SetContext(svc.query.ctx).SetHeaders(svc.query.header).SetHeader("Content-Type", "application/json").SetBody(body).SetQueryParamsFromValues(svc.query.values).Post(svc.query.targetUrl)
	if err != nil {
		return svc.query.response.SetError(err)
	}

	return svc.query.response.SetBody(res.Body(), res.StatusCode())
}

func (svc *request) Put() IResponse {
	res, err := svc.SkipVerify().
		SetDebug(svc.query.debug).
		SetTimeout(svc.query.timeout).
		R().SetContext(svc.query.ctx).SetHeaders(svc.query.header).SetBody(svc.query.body).SetQueryParamsFromValues(svc.query.values).Put(svc.query.targetUrl)
	if err == nil && res != nil {
		return svc.query.response.SetBody(res.Body(), res.StatusCode())
	}

	return svc.query.response.SetError(err)
}

func (svc *request) Delete() IResponse {
	res, err := svc.SkipVerify().
		SetDebug(svc.query.debug).
		SetTimeout(svc.query.timeout).
		R().SetContext(svc.query.ctx).SetHeaders(svc.query.header).SetBody(svc.query.body).SetQueryParamsFromValues(svc.query.values).Delete(svc.query.targetUrl)
	if err == nil && res != nil {
		return svc.query.response.SetBody(res.Body(), res.StatusCode())
	}

	return svc.query.response.SetError(err)
}