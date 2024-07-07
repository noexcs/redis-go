package datastruct

import "time"

type VolatileKey struct {
	Name      string
	ExpiredAt time.Time
}

func NewVolatileKey(name string, expiredAt time.Time) *VolatileKey {
	return &VolatileKey{
		Name:      name,
		ExpiredAt: expiredAt,
	}
}
