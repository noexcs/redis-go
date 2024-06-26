package resp3

import "github.com/noexcs/redis-go/redis/parser/resp2"

//  RESP data type	   Minimal protocol version	    Category		First byte
//  Nulls		       RESP3	                    Simple	        _
//  Booleans	       RESP3	                    Simple	        #
//  Doubles	           RESP3	                    Simple	        ,
//  Big numbers	       RESP3	                    Simple	        (
//  Bulk errors	       RESP3	                    Aggregate	    !
//  Verbatim strings   RESP3	                    Aggregate	    =
//  Maps	           RESP3	                    Aggregate	    %
//  Sets	           RESP3	                    Aggregate	    ~
//  Pushes	           RESP3	                    Aggregate	    >

// Nulls
// Null Bulk String, Null Arrays and Nulls
// Due to historical reasons, RESP2 features two specially crafted values for representing null values of bulk strings and arrays.
// This duality has always been a redundancy that added zero semantical value to the protocol itself.
// The null type, introduced in RESP3, aims to fix this wrong.
type Nulls struct {
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
