package simpleDB

import (
	"github.com/noexcs/redis-go/database/datastruct"
)

type SingleDB struct {
	index int
	data  *datastruct.BPlusTree
}

func (d *SingleDB) GetValue(key string) (any, bool) {
	v, ok := d.data.Get(key)
	return v, ok
}

func (d *SingleDB) SetValue(key string, value any) {
	d.data.Insert(key, value)
}

func (d *SingleDB) Delete(key string) bool {
	return d.data.Delete(key)
}

func (d *SingleDB) FlushDb() {
	d.data.Clear()
}

func NewSingleDB() *SingleDB {
	return &SingleDB{index: 0, data: datastruct.MakeBPlusTree()}
}
