package simpleDB

import (
	"github.com/noexcs/redis-go/database/datastruct"
	"math/rand"
	"time"
)

type GoMapDB struct {
	index int
	data  map[string]interface{}
	// 存储带过期时间的键
	volatileKeys map[string]*datastruct.VolatileKey
}

func (d *GoMapDB) GetValue(key string) (any, bool) {
	// 先检查是否有带过期时间的键
	if volatileKey, exists := d.volatileKeys[key]; exists {
		// 检查是否过期
		if volatileKey.ExpiredAt.Before(time.Now()) {
			// 过期则删除
			delete(d.data, key)
			delete(d.volatileKeys, key)
			return nil, false
		}
		// 未过期则返回值
		if value, ok := d.data[key]; ok {
			return value, true
		}
		return nil, false
	}

	// 检查普通键
	value, ok := d.data[key]
	return value, ok
}

func (d *GoMapDB) SetValueWithExpiration(key string, value any, expiration time.Time) {
	// 添加或更新数据
	d.data[key] = value

	// 添加或更新过期时间
	if d.volatileKeys == nil {
		d.volatileKeys = make(map[string]*datastruct.VolatileKey)
	}
	d.volatileKeys[key] = &datastruct.VolatileKey{
		Name:      key,
		ExpiredAt: expiration,
	}
}

func (d *GoMapDB) SetValueWithKeepTTL(key string, value any) {
	d.data[key] = value
	// 不修改过期时间
}

func (d *GoMapDB) SetValue(key string, value any) {
	d.data[key] = value

	// 如果该键之前是过期键，现在变为永久键
	if _, exists := d.volatileKeys[key]; exists {
		delete(d.volatileKeys, key)
	}
}

func (d *GoMapDB) Delete(key string) bool {
	// 删除数据
	_, exists := d.data[key]
	if exists {
		delete(d.data, key)
		// 同时删除过期时间记录
		delete(d.volatileKeys, key)
		return true
	}
	return false
}

func (d *GoMapDB) FlushDb() {
	d.data = make(map[string]interface{})
	d.volatileKeys = make(map[string]*datastruct.VolatileKey)
}

// RandomExpiredKeys 定期删除策略：随机检查并删除过期键
func (d *GoMapDB) RandomExpiredKeys() {
	// 将所有过期键放入一个切片
	expiredKeys := make([]string, 0)
	for key, volatileKey := range d.volatileKeys {
		if volatileKey.ExpiredAt.Before(time.Now()) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	// 随机打乱过期键列表
	rand.Shuffle(len(expiredKeys), func(i, j int) {
		expiredKeys[i], expiredKeys[j] = expiredKeys[j], expiredKeys[i]
	})

	// 最多删除20个过期键
	count := 0
	for _, key := range expiredKeys {
		if count >= 20 {
			break
		}
		delete(d.data, key)
		delete(d.volatileKeys, key)
		count++
	}
}

// DeleteExpiredKeys 定时删除策略：主动删除所有过期键
func (d *GoMapDB) DeleteExpiredKeys() {
	// 查找所有过期键
	var expiredKeys []string
	for key, volatileKey := range d.volatileKeys {
		if volatileKey.ExpiredAt.Before(time.Now()) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	// 删除所有过期键
	for _, key := range expiredKeys {
		delete(d.data, key)
		delete(d.volatileKeys, key)
	}
}

// StartExpiredKeysDeletion 启动定期删除过期键的goroutine
func (d *GoMapDB) StartExpiredKeysDeletion() {
	go func() {
		for {
			// 每隔一段时间随机检查过期键
			time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
			d.RandomExpiredKeys()
		}
	}()
}

func NewGoMapDB() *GoMapDB {
	db := &GoMapDB{
		index:        0,
		data:         make(map[string]interface{}),
		volatileKeys: make(map[string]*datastruct.VolatileKey),
	}
	db.StartExpiredKeysDeletion()
	return db
}
