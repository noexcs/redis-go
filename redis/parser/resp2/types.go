package resp2

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/noexcs/redis-go/redis/parser/resp"
)

var CRLF = "\r\n"

type SimpleString struct {
	Data string
}

func (r *SimpleString) ToBytes() []byte {
	return []byte("+" + r.Data + CRLF)
}
func (r *SimpleString) String() string {
	return r.Data
}

type SimpleError struct {
	Kind string
	Data string
}

func (r *SimpleError) ToBytes() []byte {
	return []byte("-" + r.Kind + " " + r.Data + CRLF)
}
func (r *SimpleError) String() string {
	return fmt.Sprintf("%s %s", r.Kind, r.Data)
}

type Integer struct {
	Data int64
}

func (r *Integer) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(r.Data, 10) + CRLF)
}
func (r *Integer) String() string {
	return strconv.FormatInt(r.Data, 10)
}

type BulkString struct {
	Data []byte
}

func (r *BulkString) ToBytes() []byte {
	if r.Data == nil {
		return []byte("$-1\r\n")
	}
	return []byte("$" + strconv.Itoa(len(r.Data)) + CRLF + string(r.Data) + CRLF)
}
func (r *BulkString) String() string {
	if r.Data == nil {
		return "(nil)"
	}
	return string(r.Data)
}

type Array struct {
	Data   []resp.RespValue
	Length int
}

func (r *Array) ToBytes() []byte {
	if r.Data == nil {
		return []byte("*-1\r\n")
	}
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(len(r.Data)) + CRLF)
	for _, datum := range r.Data {
		buf.Write(datum.ToBytes())
	}
	return buf.Bytes()
}
func (r *Array) String() string {
	var builder strings.Builder
	for i, datum := range r.Data {
		builder.WriteString(datum.String())
		if i != r.Length-1 {
			builder.WriteString(" ")
		}
	}
	return builder.String()
}

var OKSimpleString = SimpleString{Data: "OK"}

func MakeSimpleString(data string) *SimpleString {
	return &SimpleString{Data: data}
}

func MakeOKSimpleString() *SimpleString {
	return &OKSimpleString
}

func MakePONGSimpleString() *SimpleString {
	return &SimpleString{Data: "PONG"}
}

func MakeNullBulkString() *BulkString {
	return &BulkString{Data: nil}
}

func MakeSimpleError(kind string, data string) *SimpleError {
	return &SimpleError{Kind: kind, Data: data}
}
