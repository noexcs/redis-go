package simpleDB

import (
	"github.com/noexcs/redis-go/database/datastruct"
	"time"
)

type SingleDB struct {
	index int
	data  *datastruct.BPlusTree[interface{}]
}

func (d *SingleDB) GetValue(key string) (any, bool) {
	k, v, exist := d.data.Get(key)
	if !exist {
		return nil, false
	}
	if volatileKey, ok := k.(*datastruct.VolatileKey); ok {
		if volatileKey.ExpiredAt.Before(time.Now()) {
			d.Delete(key)
			return nil, false
		}
		return v, true
	}
	return v, true
}

func (d *SingleDB) SetValueWithExpiration(key string, value any, expiration time.Time) {
	d.data.Insert(&datastruct.VolatileKey{
		Name:      key,
		ExpiredAt: expiration,
	}, value, true)
}

func (d *SingleDB) SetValueWithKeepTTL(key string, value any) {
	d.data.Insert(key, value, false)
}

func (d *SingleDB) SetValue(key string, value any) {
	d.data.Insert(key, value, true)
}

func (d *SingleDB) Delete(key string) bool {
	return d.data.Delete(key)
}

func (d *SingleDB) FlushDb() {
	d.data.Clear()
}

func NewSingleDB() *SingleDB {
	return &SingleDB{index: 0, data: datastruct.MakeBPlusTree[interface{}]()}
}
