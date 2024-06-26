package parser

import (
	"github.com/noexcs/redis-go/redis/parser/resp2"
)

type Request struct {
	Args resp2.RespType
	Err  *Error
}

func MakeErrorRequest(kind string, message string) *Request {
	return &Request{Err: &Error{
		Kind:    kind,
		Message: message,
	}}
}

type Response Request

func MakeErrorResponse(kind string, message string) *Response {
	return &Response{Err: &Error{
		Kind:    kind,
		Message: message,
	}}
}
