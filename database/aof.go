package database

import (
	"bufio"
	"github.com/noexcs/redis-go/log"
	"github.com/noexcs/redis-go/redis/parser/resp"
	"github.com/noexcs/redis-go/redis/parser/resp2"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type AofHandler struct {
	aofFile  *os.File
	aofMutex sync.Mutex

	// fsync相关
	fsyncStrategy string    // 同步策略: always, everysec, no
	lastSyncTime  time.Time // 上次同步时间
	db            DB
}

// NewAofHandler 创建一个新的AOF处理器
func NewAofHandler(db DB, filename string, fsyncStrategy string) (*AofHandler, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	handler := &AofHandler{
		aofFile:       file,
		db:            db,
		fsyncStrategy: fsyncStrategy,
		lastSyncTime:  time.Now(),
	}

	// 加载现有AOF文件
	handler.LoadAof()

	// 如果策略是everysec，启动后台同步goroutine
	if fsyncStrategy == "everysec" {
		go handler.fsyncEverySecond()
	}

	return handler, nil
}

// fsyncEverySecond 每秒执行一次fsync
func (handler *AofHandler) fsyncEverySecond() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		handler.aofMutex.Lock()
		// 检查是否需要同步(至少1秒间隔)
		if time.Since(handler.lastSyncTime) >= time.Second {
			handler.aofFile.Sync()
			handler.lastSyncTime = time.Now()
		}
		handler.aofMutex.Unlock()
	}
}

// LoadAof 从AOF文件加载数据
func (handler *AofHandler) LoadAof() {
	// 获取文件信息以检查文件大小
	fileInfo, err := handler.aofFile.Stat()
	if err != nil {
		log.Error("Failed to get AOF file info: %v", err)
		return
	}

	// 如果文件为空，直接返回
	if fileInfo.Size() == 0 {
		log.Info("AOF file is empty")
		return
	}

	// 读取整个文件内容
	content, err := io.ReadAll(handler.aofFile)
	if err != nil {
		log.Error("Failed to read AOF file: %v", err)
		return
	}

	// 创建RESP解析器来解析AOF内容
	reader := bufio.NewReader(strings.NewReader(string(content)))

	for {
		respParser := &resp2.Resp2Parser{}
		value, err := respParser.Parse(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Error("Error parsing AOF file: %v", err)
			break
		}

		// 执行命令
		if array, ok := value.(*resp2.Array); ok {
			handler.dbExec(array.Data)
		}
	}

	log.Info("AOF file loaded successfully")
}

// dbExec 执行数据库命令
func (handler *AofHandler) dbExec(args []resp.RespValue) {
	if len(args) < 1 {
		return
	}

	commandName := strings.ToUpper(args[0].String())

	// 根据命令名称调用相应的数据库操作
	switch commandName {
	case "SET":
		if len(args) >= 3 {
			key := args[1].String()
			value := args[2].String()
			handler.db.SetValue(key, value)
		}
	case "DEL":
		if len(args) >= 2 {
			key := args[1].String()
			handler.db.Delete(key)
		}
	}
}

// AddAof 将命令添加到AOF文件
func (handler *AofHandler) AddAof(args []resp.RespValue) {
	handler.aofMutex.Lock()
	defer handler.aofMutex.Unlock()

	// 构造RESP数组，确保Length字段正确设置
	array := &resp2.Array{
		Length: len(args),
		Data:   args,
	}

	// 将命令转换为RESP格式并写入文件
	_, err := handler.aofFile.Write(array.ToBytes())
	if err != nil {
		log.Error("Failed to write to AOF file: %v", err)
		return
	}

	// 根据同步策略决定是否同步
	switch handler.fsyncStrategy {
	case "always":
		// 立即同步到磁盘
		handler.aofFile.Sync()
		handler.lastSyncTime = time.Now()
	case "everysec":
		// everysec策略由后台goroutine处理
	case "no":
		// 不主动同步，由操作系统决定
	}
}

// Close 关闭AOF处理器
func (handler *AofHandler) Close() error {
	handler.aofMutex.Lock()
	defer handler.aofMutex.Unlock()

	// 关闭前最后一次同步
	handler.aofFile.Sync()

	return handler.aofFile.Close()
}
