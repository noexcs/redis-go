package resp2

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

//  RESP data type	 Minimal protocol version	Category		First byte

//  Simple strings	 RESP2	    				Simple			+
//  Simple Errors	 RESP2	    				Simple			-
//  Integers	     RESP2	    				Simple	    	:
//  Bulk strings	 RESP2	    				Aggregate		$
//  Arrays	         RESP2	    				Aggregate		*

// resp2: https://redis.io/docs/reference/protocol-spec

var CRLF = "\r\n"

type RespType interface {
	ToBytes() []byte
	String() string
}

// ====================SimpleString======================

// SimpleString 以(+)开头，以(\r\n)结尾，中间为字符串
// 例如：+OK\r\n
type SimpleString struct {
	Data string
}

func (r *SimpleString) ToBytes() []byte {
	return []byte("+" + r.Data + CRLF)
}

func (r *SimpleString) String() string {
	return r.Data
}

var oKSimpleString = SimpleString{Data: "OK"}

func MakeSimpleString(data string) *SimpleString {
	return &SimpleString{Data: data}
}

func MakeOKSimpleString() *SimpleString {
	return &oKSimpleString
}

func MakePONGSimpleString() *SimpleString {
	return &SimpleString{Data: "PONG"}
}

// ====================SimpleError======================

// SimpleError 以(-)开头，以(\r\n)结尾，中间为字符串
// 例如：-ERR unknown command 'asdf'\r\n
//
//	-WRONGTYPE Operation against a key holding the wrong kind of value
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

func MakeSimpleError(kind string, data string) *SimpleError {
	return &SimpleError{Kind: kind, Data: data}
}

// =====================Integer=====================

// Integer 以:开头，以(\r\n)结尾，中间为整数，:[<+|->]<value>\r\n
// 例如：:1000\r\n
type Integer struct {
	Data int64
}

func (r Integer) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(r.Data, 10) + CRLF)
}

func (r *Integer) String() string {
	return strconv.FormatInt(r.Data, 10)
}

// ====================BulkString======================

// BulkString 表示一个二进制字符串，$<length>\r\n<data>\r\n
// 以$<length>\r\n开头，以<data>\r\n结尾，中间为字符串
// 例如：$5\r\nhello\r\n
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

func MakeNullBulkString() *BulkString {
	return &BulkString{Data: nil}
}

// =====================Array=====================

// Array 表示一个数组，*<number-of-elements>\r\n<element-1>...<element-n>
// 例如：*3\r\n+OK\r\n+OK\r\n+OK\r\n
type Array struct {
	Data   []*RespType
	Length int
}

func (r *Array) ToBytes() []byte {
	if r.Data == nil {
		return []byte("*-1\r\n")
	}
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(len(r.Data)) + CRLF)
	for _, datum := range r.Data {
		buf.Write((*datum).ToBytes())
	}
	return buf.Bytes()
}

func (r *Array) String() string {
	var builder strings.Builder

	for i, datum := range r.Data {
		builder.WriteString((*datum).String())
		if i != r.Length-1 {
			builder.WriteString(" ")
		}
	}

	return builder.String()
}
