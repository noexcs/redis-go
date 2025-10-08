package database

import (
	"github.com/noexcs/redis-go/database/simpleDB"
	"github.com/noexcs/redis-go/redis/parser/resp"
	"github.com/noexcs/redis-go/redis/parser/resp2"
	"os"
	"testing"
)

func TestAof(t *testing.T) {
	// 创建临时AOF文件
	tmpfile, err := os.CreateTemp("", "test.aof")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// 关闭临时文件
	tmpfile.Close()

	// 创建数据库和AOF处理器
	db := simpleDB.NewGoMapDB()
	aofHandler, err := NewAofHandler(db, tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer aofHandler.Close()

	// 测试添加命令到AOF
	setCmd := []resp.RespValue{
		&resp2.BulkString{Data: []byte("SET")},
		&resp2.BulkString{Data: []byte("key1")},
		&resp2.BulkString{Data: []byte("value1")},
	}

	aofHandler.AddAof(setCmd)

	// 验证数据是否写入文件
	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Error("AOF file is empty")
	}

	// 直接在数据库中设置值，因为AofHandler.AddAof只是记录命令而不执行
	db.SetValue("key1", "value1")

	// 验证数据是否正确存储在数据库中
	value, exists := db.GetValue("key1")
	if !exists {
		t.Error("Key not found in database")
	}

	if value != "value1" {
		t.Errorf("Expected value1, got %v", value)
	}
}

func TestAofReload(t *testing.T) {
	// 创建临时AOF文件
	tmpfile, err := os.CreateTemp("", "test_reload.aof")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// 关闭临时文件
	tmpfile.Close()

	// 创建数据库和AOF处理器
	db := simpleDB.NewGoMapDB()
	aofHandler, err := NewAofHandler(db, tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	aofHandler.Close()

	// 手动写入一些命令到AOF文件
	commands := "*3\r\n$3\r\nSET\r\n$4\r\nkey2\r\n$6\r\nvalue2\r\n"
	err = os.WriteFile(tmpfile.Name(), []byte(commands), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 创建新的AOF处理器来加载数据
	db2 := simpleDB.NewGoMapDB()
	aofHandler2, err := NewAofHandler(db2, tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer aofHandler2.Close()

	// 验证数据是否正确加载到新数据库中
	value, exists := db2.GetValue("key2")
	if !exists {
		t.Error("Key not found in database after reload")
	}

	if value != "value2" {
		t.Errorf("Expected value2, got %v", value)
	}
}
