package simpleDB

import (
	"github.com/noexcs/redis-go/database/datastruct"
	"math/rand"
	"sync"
	"time"
)

type SingleDB struct {
	index int
	data  *datastruct.BPlusTree[interface{}]
	mutex sync.RWMutex
}

func (d *SingleDB) GetValue(key string) (any, bool) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	k, v, exist := d.data.Get(key)
	if !exist {
		return nil, false
	}
	if volatileKey, ok := k.(*datastruct.VolatileKey); ok {
		if volatileKey.ExpiredAt.Before(time.Now()) {
			d.mutex.RUnlock()
			d.mutex.Lock()
			d.data.Delete(key)
			d.mutex.Unlock()
			d.mutex.RLock()
			return nil, false
		}
		return v, true
	}
	return v, true
}

func (d *SingleDB) SetValueWithExpiration(key string, value any, expiration time.Time) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.data.Insert(&datastruct.VolatileKey{
		Name:      key,
		ExpiredAt: expiration,
	}, value, true)
}

func (d *SingleDB) SetValueWithKeepTTL(key string, value any) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.data.Insert(key, value, false)
}

func (d *SingleDB) SetValue(key string, value any) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.data.Insert(key, value, true)
}

func (d *SingleDB) Delete(key string) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.data.Delete(key)
}

func (d *SingleDB) FlushDb() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.data.Clear()
}

// RandomExpiredKeys 定期删除策略：随机检查并删除过期键
func (d *SingleDB) RandomExpiredKeys() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// 获取所有键的迭代器
	iterator := d.data.Iterator()
	count := 0

	// 随机检查最多20个键
	for iterator.Next() && count < 20 {
		k, _ := iterator.Value()
		if volatileKey, ok := k.(*datastruct.VolatileKey); ok {
			if volatileKey.ExpiredAt.Before(time.Now()) {
				d.data.Delete(volatileKey.Name)
			}
		}
		count++
	}

	// 释放迭代器
	iterator.Discard()
}

// DeleteExpiredKeys 定时删除策略：主动删除所有过期键
func (d *SingleDB) DeleteExpiredKeys() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// 获取所有键的迭代器
	iterator := d.data.Iterator()

	// 存储过期键的列表
	var expiredKeys []string

	// 查找所有过期键
	for iterator.Next() {
		k, _ := iterator.Value()
		if volatileKey, ok := k.(*datastruct.VolatileKey); ok {
			if volatileKey.ExpiredAt.Before(time.Now()) {
				expiredKeys = append(expiredKeys, volatileKey.Name)
			}
		}
	}

	// 释放迭代器
	iterator.Discard()

	// 删除所有过期键
	for _, key := range expiredKeys {
		d.data.Delete(key)
	}
}

// StartExpiredKeysDeletion 启动定期删除过期键的goroutine
func (d *SingleDB) StartExpiredKeysDeletion() {
	go func() {
		for {
			// 每隔一段时间随机检查过期键
			time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
			d.RandomExpiredKeys()
		}
	}()
}

func NewSingleDB() *SingleDB {
	db := &SingleDB{index: 0, data: datastruct.MakeBPlusTree[interface{}]()}
	db.StartExpiredKeysDeletion()
	return db
}
