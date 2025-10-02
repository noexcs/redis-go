package command

import (
	"fmt"
	"github.com/noexcs/redis-go/database"
	"github.com/noexcs/redis-go/redis/parser"
	"github.com/noexcs/redis-go/redis/parser/resp"
	"github.com/noexcs/redis-go/redis/parser/resp2"
	"strconv"
	"strings"
	"time"
)

func init() {
	RegisterCommand("set", execSet, nil, nil, -3, FlagWrite)
	RegisterCommand("get", execGet, nil, nil, 2, FlagReadonly)
	RegisterCommand("getrange", execGetRange, nil, nil, 4, FlagReadonly)
	RegisterCommand("incr", execIncr, nil, nil, 2, FlagWrite)
}

// Reference: https://redis.io/docs/latest/commands/set/
//
// Syntax:
//
// SET key value [NX | XX] [GET] [EX seconds | PX milliseconds | EXAT unix-time-seconds | PXAT unix-time-milliseconds | KEEPTTL]
func execSet(db database.DB, args []resp.RespValue) *parser.Response {
	data := args
	key := data[1].String()
	value := data[2].String()
	db.SetValue(key, value)

	options := make(map[string]interface{})
	for i := 3; i < len(data); i++ {
		Option := strings.ToUpper(data[i].String())
		switch Option {
		case "NX", "XX":
			{
				if _, ok := options["NXXX"]; ok && options["NXXX"] != Option {
					return &parser.Response{
						Args: &resp2.SimpleError{
							Kind: "ERR",
							Data: "Different NX and XX are set repeatedly.",
						},
						Err: nil,
					}
				}
				options["NXXX"] = Option
			}
		case "GET":
			options["GET"] = struct{}{}
		case "EX", "PX", "EXAT", "PXAT", "KEEPTTL":
			{
				// If two different TTLs are set repeatedly.
				if _, ok := options["TTL"]; ok && options["TTL"] != Option {
					return &parser.Response{
						Args: resp2.MakeNullBulkString(),
						Err:  nil,
					}
				}
				if Option == "KEEPTTL" {
					options["TTL"] = Option
				} else if i+1 >= len(data) {
					return &parser.Response{
						Args: &resp2.SimpleError{
							Kind: "ERR",
							Data: fmt.Sprintf("Insufficient parameter for option %s.", Option),
						},
						Err: nil,
					}
				} else {
					duration := data[i+1]
					atoi, err := strconv.ParseInt(duration.String(), 10, 64)
					if err != nil {
						return &parser.Response{
							Args: &resp2.SimpleError{
								Kind: "ERR",
								Data: fmt.Sprintf("Invalid parameter for option %s.", Option),
							},
							Err: nil,
						}
					}
					options["TTL"] = []interface{}{Option, atoi}
					i++
				}
			}
		}
	}

	existedValue, exist := db.GetValue(key)
	if _, ok := options["NXXX"]; ok {
		if (options["NXXX"] == "NX" && !exist) || (options["NXXX"] == "XX" && exist) {
			if ttlOption, ok := options["TTL"]; ok {
				if ttl, ok := ttlOption.([]interface{}); ok {
					d := ttl[1].(int64)
					if ttl[0] == "EX" {
						db.SetValueWithExpiration(key, value, time.Now().Add(time.Duration(d)*time.Second))
					} else if ttl[0] == "PX" {
						db.SetValueWithExpiration(key, value, time.Now().Add(time.Duration(d)*time.Millisecond))
					} else if ttl[0] == "EXAT" {
						db.SetValueWithExpiration(key, value, time.Unix(d, 0))
					} else if ttl[0] == "PXAT" {
						db.SetValueWithExpiration(key, value, time.Unix(d, 0))
					} else if ttl[0] == "KEEPTTL" {
						db.SetValueWithKeepTTL(key, value)
					}
				}
			} else {
				db.SetValue(key, value)
			}
		} else {
			simpleError := &resp2.SimpleError{
				Kind: "ERR",
				Data: "The key does not exist.",
			}
			if exist {
				simpleError.Data = "The key has already exist."
			}
			return &parser.Response{
				Args: simpleError,
				Err:  nil,
			}
		}
	} else {
		if ttlOption, ok := options["TTL"]; ok {
			if ttl, ok := ttlOption.([]interface{}); ok {
				d := ttl[1].(int64)
				if ttl[0] == "EX" {
					db.SetValueWithExpiration(key, value, time.Now().Add(time.Duration(d)*time.Second))
				} else if ttl[0] == "PX" {
					db.SetValueWithExpiration(key, value, time.Now().Add(time.Duration(d)*time.Millisecond))
				} else if ttl[0] == "EXAT" {
					db.SetValueWithExpiration(key, value, time.Unix(d, 0))
				} else if ttl[0] == "PXAT" {
					db.SetValueWithExpiration(key, value, time.Unix(d, 0))
				} else if ttl[0] == "KEEPTTL" {
					db.SetValueWithKeepTTL(key, value)
				}
			}
		} else {
			db.SetValue(key, value)
		}
	}
	if _, ok := options["GET"]; ok {
		if exist {
			return &parser.Response{Args: &resp2.BulkString{Data: []byte(existedValue.(string))}, Err: nil}
		} else {
			return &parser.Response{Args: resp2.MakeNullBulkString(), Err: nil}
		}
	} else {
		if exist {
			return &parser.Response{Args: &resp2.SimpleString{Data: "OK"}, Err: nil}
		} else {
			return &parser.Response{Args: resp2.MakeNullBulkString(), Err: nil}
		}
	}
}

// GET key
// Return
// Bulk string reply: the value of key, or nil when key does not exist.
func execGet(db database.DB, args []resp.RespValue) *parser.Response {
	data := args
	key := data[1].String()
	value, exist := db.GetValue(key)
	if exist {
		if v, ok := value.(string); ok {
			return &parser.Response{Args: &resp2.BulkString{Data: []byte(v)}, Err: nil}
		} else {
			return &parser.Response{
				Args: &resp2.SimpleError{
					Kind: "WRONGTYPE",
					Data: "Operation against a key holding the wrong kind of value",
				},
				Err: nil,
			}
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
func execGetRange(db database.DB, args []resp.RespValue) *parser.Response {
	data := args
	key := data[1].String()
	start := data[2].(*resp2.BulkString).Data
	end := data[3].(*resp2.BulkString).Data

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

// Increments the number stored at key by one.
// If the key does not exist, it is set to 0 before performing the operation.
// An error is returned if the key contains a value of the wrong type or contains a string
// that can not be represented as integer. This operation is limited to 64-bit signed integers.
//
// Note: this is a string operation because Redis does not have a dedicated integer type.
// The string stored at the key is interpreted as a base-10 64-bit signed integer to execute the operation.
//
// Redis stores integers in their integer representation, so for string values that actually hold an integer,
// there is no overhead for storing the string representation of the integer.
func execIncr(db database.DB, args []resp.RespValue) *parser.Response {
	data := args
	key := data[1].String()
	value, exist := db.GetValue(key)
	if !exist {
		db.SetValue(key, "1")
		return &parser.Response{Args: &resp2.Integer{Data: 1}}
	} else {
		i, err := strconv.ParseInt(value.(string), 0, 64)
		if err != nil {
			return &parser.Response{Args: &resp2.SimpleError{
				Kind: "ERR",
				Data: "Operation against a key holding the wrong kind of value",
			}}
		}
		s := strconv.FormatInt(i+1, 10)
		db.SetValue(key, s)
		return &parser.Response{Args: &resp2.Integer{Data: i + 1}}
	}
}
