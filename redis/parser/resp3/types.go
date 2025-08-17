package resp3

import (
	"bytes"
	"fmt"
	"github.com/noexcs/redis-go/redis/parser/resp2"
	"strconv"
	"strings"
)

// RESP data type	   Minimal protocol version	    Category		First byte
// Nulls		       RESP3	                    Simple	        _
// Booleans	           RESP3	                    Simple	        #
// Doubles	           RESP3	                    Simple	        ,
// Big numbers	       RESP3	                    Simple	        (
// Bulk errors	       RESP3	                    Aggregate	    !
// Verbatim strings    RESP3	                    Aggregate	    =
// Maps	               RESP3	                    Aggregate	    %
// Sets	               RESP3	                    Aggregate	    ~
// Pushes	           RESP3	                    Aggregate	    >

var CRLF = "\r\n"

type RespType interface {
	ToBytes() []byte
	String() string
}

// BlobString
// The general form is `$<length>\r\n<bytes>\r\n`.
// It is basically exactly like in the previous version of RESP.
//
// The string `"hello world"` is represented by the following protocol:
//
//	 ```
//		$11<CR><LF>
//		hello world<CR><LF>
//	 ```
//
// Or as an escaped string:
//
//	 ```
//		"$11\r\nhello world\r\n"
//	 ```
//
// The length field is limited to the range of an unsigned 64-bit
// integer. Zero is a valid length, so the empty string is represented by:
//
//	"$0\r\n\r\n"
type BlobString struct {
	Data []byte
}

func (b *BlobString) ToBytes() []byte {
	//TODO implement me
	panic("implement me")
}

func (b *BlobString) String() string {
	//TODO implement me
	panic("implement me")
}

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

// Number `:<number>\r\n`
type Number struct {
	int64
}

type Null struct {
	string
}

func MakeNull() Null {
	return Null{"_\r\n"}
}

// Double `,<floating-point-number>\r\n`
type Double struct {
}

// Boolean `#t\r\n` or `#f\r\n`
type Boolean struct {
}

// BlobError `!<length>\r\n<bytes>\r\n`
type BlobError struct {
}

// VerbatimString `=<length>\r\n<bytes>\r\n`.
//
// The first three bytes provide information about the format
// of the following string, which can be `txt` for plain text, or `mkd` for
// markdown. The fourth byte is always `:`.
type VerbatimString struct {
}

// BigNumber `(<big number>\r\n`
type BigNumber struct {
}

// Nulls
// Null Bulk String, Null Arrays and Nulls
// Due to historical reasons, RESP2 features two specially crafted values for representing null values of bulk strings and arrays.
// This duality has always been a redundancy that added zero semantical value to the protocol itself.
// The null type, introduced in RESP3, aims to fix this wrong.
type Nulls struct {
}

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

// Map `%<numelements><CR><LF> ... numelements other types ...`
//
// The number of following elements must be even.
// Maps represent a sequence of field-value items.
type Map struct {
	Data   map[string]*RespType
	Length int
}

// Set `~<numelements><CR><LF> ... numelements other types ...`
type Set struct {
}

var nullsBytes = []byte("_" + resp2.CRLF)

func (n *Nulls) ToBytes() []byte {
	return nullsBytes
}

func (n *Nulls) String() string {
	return "(nil)"
}

var nulls = &Nulls{}

func MakeNulls() *Nulls {
	return nulls
}
