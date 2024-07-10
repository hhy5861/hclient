### hclient是http(s)远程调用，基于 [resty](https://github.com/go-resty/resty) 二次封装。支持多服务调用。依赖配置管理

* 封装扩展两个监控埋点
  * [prometheus metrics](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#Counter) metrics指标监控。
  * [opentracing](https://opentracing.io/) 请求trace链路跟踪，方便定位问题，性能问题定位

* 快速使用
```go
1. 在项目中添加配置解析结构体

type (
	Config struct {
		Client *Client           `yaml:"client"`
	}

	Client struct {
		Http map[string]*hclient.Remote `yaml:"http"`
	}
)

2. 初始化配置

func (m *Injection) NewHttpClient() {
	hclient.NewClient(
		hclient.SetRemotes(m.cfg.Client.Http),
	)
}
```

* 在项目中使用DEMO
```go
package service

import (
	"github.com/hhy5861/foo-service/internal/constants"
	"github.com/hhy5861/foo-service/internal/models"
	"github.com/hhy5861/foo-service/model/form"
	"github.com/hhy5861/gclient"
	"github.com/hhy5861/hclient"
	"context"
	"net/http"
	"strconv"
)

type (
	FooService struct {
		ctx context.Context
	}
)

func NewFooService(ctx context.Context) *FooService {
	return &FooService{
		ctx: ctx,
	}
}

func (svc *FooService) GetServiceInfo(params form.FooForm) (*models.Auth, error) {
	var (
		result models.Auth
	)

    params = map[string]interface{}{
		"userId": strconv.FormatInt(params.UserId, 10),
	}

	r := hclient.GlobalWithCtx(svc.ctx).Get(constants.BarService, constants.BarServiceCommonAuthList, params)
	if r.GetHttpStatus() == http.StatusOK {
		err := r.GetStruct(&result)
		if err != nil {
			return result, err
		}
	}

	return &result, nil
}
```
* 注意service name必须和配置里的key一致，否则找不到配置

* 初始化配置项目

```go
// SetRemotes 配置调用服务所有列表map
func SetRemotes(remote map[string]*Remote) Option {
	return func(c *ConfigCache) {
		c.remotes = remote
	}
}

// SetResponse 自定义返回数据结构解析
func SetResponse(r IResponse) Option {
	return func(c *ConfigCache) {
		c.resp = r
	}
}

// Opentracing 配置tracer客服端连接器
func Opentracing(t opentracing.Tracer) Option {
	return func(c *ConfigCache) {
		c.trace = t
	}
}
```