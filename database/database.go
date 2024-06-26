package database

type DB interface {
	GetValue(key string) (value any, exist bool)

	SetValue(key string, value any)

	Delete(key string)

	FlushDb()
}
