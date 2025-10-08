package parser

import (
	"github.com/noexcs/redis-go/redis/parser/resp"
)

type Request struct {
	Args resp.RespValue
	Err  *resp.Error
}

func MakeErrorRequest(kind string, message string) *Request {
	return &Request{Err: &resp.Error{
		Kind:    kind,
		Message: message,
	}}
}

type Response Request

func MakeErrorResponse(kind string, message string) *Response {
	return &Response{Err: &resp.Error{
		Kind:    kind,
		Message: message,
	}}
}
