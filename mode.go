package hclient

type (
	IResponse interface {
		GetHttpStatus() int

		GetSuccess() int

		GetCode() int

		GetData() interface{}

		GetMessage() string

		GetStruct(data interface{}) error

		GetBody() []byte

		GetStructByBody(data interface{}) error

		IsOk() bool

		GetError() error
	}
)
