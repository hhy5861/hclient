package hclient

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

type (
	query struct {
		client *resty.Client
	}

	Query struct {
		body       interface{}
		debug      bool
		params     map[string]string
		header     map[string]string
		timeout    time.Duration
		response   *Response
		targetUrl  string
		skipVerify bool
		ctx        context.Context
		values     url.Values
	}
)

func NewQuery(c *resty.Client) *query {
	return &query{
		client: c,
	}
}

func (q *query) Query(content interface{}) map[string]string {
	v := reflect.ValueOf(content)
	switch v.Kind() {
	case reflect.String:
		return q.String(v.String())
	case reflect.Struct, reflect.Ptr:
		return q.Struct(v.Interface())
	case reflect.Map:
		return q.Map(v.Interface())
	default:
		return nil
	}
}

func (q *query) String(content string) map[string]string {
	params := make(map[string]string)

	if err := json.Unmarshal([]byte(content), &params); err == nil {
		return params
	} else {
		if queryData, err := url.ParseQuery(content); err == nil {
			for k, queryValues := range queryData {
				for _, queryValue := range queryValues {
					params[k] = queryValue
				}
			}
		}
	}

	return params
}

func (q *query) Struct(content interface{}) map[string]string {
	var (
		val    map[string]interface{}
		params = make(map[string]string)
	)

	if marshalContent, err := json.Marshal(content); err == nil {
		if err := json.Unmarshal(marshalContent, &val); err == nil {
			for k, v := range val {
				var queryVal string

				switch t := v.(type) {
				case string:
					queryVal = t
				case float64:
					queryVal = strconv.FormatFloat(t, 'f', -1, 64)
				case time.Time:
					queryVal = t.Format(time.RFC3339)
				default:
					j, err := json.Marshal(v)
					if err != nil {
						continue
					}

					queryVal = string(j)
				}

				params[k] = queryVal
			}
		}
	}

	return params
}

func (q *query) Map(content interface{}) map[string]string {
	return q.Struct(content)
}

func (q *Query) Add(key, value string) {
	q.params[key] = value
}
