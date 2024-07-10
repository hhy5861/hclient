package hclient

import (
	"encoding/json"
	"log"
)

type (
	DefaultResp struct {
		Status int         `json:"status"`
		Code   int         `json:"code"`
		Msg    string      `json:"msg"`
		Data   interface{} `json:"data"`
	}
)

func NewDefaultResp() *DefaultResp {
	return &DefaultResp{
		Status: 900,
		Code:   900,
		Msg:    "the returned data structure does not conform to",
		Data:   nil,
	}
}

func (out *DefaultResp) ToString() string {
	body, err := json.Marshal(out)
	if err != nil {
		log.Println("default response errors", err)
	}

	return string(body)
}

func (out *DefaultResp) ToByte() []byte {
	body, err := json.Marshal(out)
	if err != nil {
		log.Println("default response errors", err)
	}

	return body
}
