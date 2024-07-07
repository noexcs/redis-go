package database

import "time"

type DB interface {
	GetValue(key string) (value any, exist bool)

	SetValue(key string, value any)

	SetValueWithExpiration(key string, value any, expiration time.Time)

	SetValueWithKeepTTL(key string, value any)

	Delete(key string) bool

	FlushDb()
}
