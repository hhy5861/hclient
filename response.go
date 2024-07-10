package hclient

import (
	"encoding/json"
	"net/http"
)

type (
	Response struct {
		Status     int         `json:"status"`
		Code       int         `json:"code"`
		Data       interface{} `json:"data"`
		Msg        string      `json:"msg"`
		Body       []byte      `json:"-"`
		HttpStatus int         `json:"-"`
		err        error
	}
)

func NewResponse() *Response {
	return &Response{}
}

func (svc *Response) SetBody(body []byte, code ...int) *Response {
	status := http.StatusBadGateway
	if len(code) == 1 {
		status = code[0]
	}

	svc.Body = body
	svc.HttpStatus = status

	if len(body) > 0 {
		json.Unmarshal(body, &svc)
	}

	return svc
}

func (svc *Response) GetHttpStatus() int {
	return svc.HttpStatus
}

func (svc *Response) GetSuccess() int {
	return svc.Status
}

func (svc *Response) GetCode() int {
	return svc.Code
}

func (svc *Response) GetData() interface{} {
	return svc.Data
}

func (svc *Response) GetMessage() string {
	return svc.Msg
}

func (svc *Response) GetStruct(data interface{}) error {
	body, err := json.Marshal(svc.Data)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, data)
}

func (svc *Response) GetBody() []byte {
	return svc.Body
}

func (svc *Response) GetStructByBody(data interface{}) error {
	err := json.Unmarshal(svc.Body, &data)
	if err != nil {
		return err
	}

	return nil
}

func (svc *Response) IsOk() bool {
	return svc.HttpStatus == http.StatusOK
}

func (svc *Response) GetError() error {
	return svc.err
}

func (svc *Response) SetError(err error) *Response {
	svc.err = err
	return svc
}
