package zinx

import (
	"time"
	"zinx/lib/logger"
)

func init() {
	// 初始化日志
	logger.Setup(&logger.Settings{
		Path:       "logs",        // 日志文件路径
		Name:       "zinx",        // 日志文件名称
		Ext:        "log",         // 日志文件后缀
		TimeFormat: time.DateOnly, // 时间格式
	})
}
