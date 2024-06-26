package command

import (
	"github.com/noexcs/redis-go/database"
	"github.com/noexcs/redis-go/redis/parser"
	"github.com/noexcs/redis-go/redis/parser/resp2"
	"strconv"
)

func init() {
	RegisterCommand("set", execSet, nil, nil, 3, FlagWrite)
	RegisterCommand("get", execGet, nil, nil, 2, FlagReadonly)
	RegisterCommand("getrange", execGetRange, nil, nil, 4, FlagReadonly)
}

// Original: SET key value [NX | XX] [GET] [EX seconds | PX milliseconds | EXAT unix-time-seconds | PXAT unix-time-milliseconds | KEEPTTL]
// Implementation: SET key value
//
// Return
//
//	Simple string reply: OK if SET was executed correctly.
//	Null reply: (nil) if the SET operation was not performed
//		because the user specified the NX or XX option but the condition was not met.
func execSet(db database.DB, args *resp2.Array) *parser.Response {
	if args.Length != 3 {
		return &parser.Response{
			Args: nil,
			Err: &parser.Error{
				Kind:    "ERR",
				Message: "WRONG number of arguments for 'SET' command",
			},
		}
	}

	data := args.Data
	key := (*data[1]).String()
	value := (*data[2]).String()
	db.SetValue(key, value)

	return &parser.Response{
		Args: &resp2.SimpleString{Data: "OK"},
		Err:  nil,
	}
}

// GET key
// Return
// Bulk string reply: the value of key, or nil when key does not exist.
func execGet(db database.DB, args *resp2.Array) *parser.Response {
	data := args.Data
	key := (*data[1]).String()
	value, exist := db.GetValue(key)
	if exist {
		if v, ok := value.(string); ok {
			return &parser.Response{Args: &resp2.BulkString{Data: []byte(v)}, Err: nil}
		} else {
			return &parser.Response{Args: &resp2.SimpleError{Kind: "WRONGTYPE", Data: "Operation against a key holding the wrong kind of value"}, Err: nil}
		}
	} else {
		return &parser.Response{Args: resp2.MakeNullBulkString(), Err: nil}
	}
}

// GETRANGE key start end
// Returns the substring of the string value stored at key, determined by the offsets start and end (both are inclusive).
// Negative offsets can be used in order to provide an offset starting from the end of the string.
// So -1 means the last character, -2 the penultimate and so forth.
//
// The function handles out of range requests by limiting the resulting range to the actual length of the string.
//
// Return
// Bulk string reply
func execGetRange(db database.DB, args *resp2.Array) *parser.Response {
	data := args.Data
	key := (*data[1]).String()
	start := (*data[2]).(*resp2.BulkString).Data
	end := (*data[3]).(*resp2.BulkString).Data

	startInt, err1 := strconv.Atoi(string(start))
	endInt, err2 := strconv.Atoi(string(end))
	if err1 != nil || err2 != nil {
		return &parser.Response{Args: &resp2.SimpleError{Kind: "ERR", Data: "Wrong range for the value."}}
	}

	value, exist := db.GetValue(key)
	if exist {
		if v, ok := value.(string); ok {
			startInt = startInt % len(v)
			endInt = endInt % len(v)
			if startInt > endInt {
				return &parser.Response{Args: &resp2.BulkString{Data: []byte("")}}
			}
			return &parser.Response{Args: &resp2.BulkString{Data: []byte(v[startInt : endInt+1])}, Err: nil}
		} else {
			return &parser.Response{Args: &resp2.BulkString{Data: []byte("")}}
		}
	}
	return &parser.Response{Args: &resp2.BulkString{Data: []byte("")}}
}
