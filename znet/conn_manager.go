package znet

import (
	"sync"
	"sync/atomic"
	errs "zinx/lib/enum/err"
	"zinx/lib/logger"
	"zinx/ziface"
)

// ConnManager 连接管理模块
type ConnManager struct {
	connections sync.Map     // 管理的连接集合
	connSize    atomic.Int32 // 连接数量
}

// Add 添加连接, 返回修改后的manager
func (manager *ConnManager) Add(conn ziface.IConnection) ziface.IConnManager {
	manager.connections.Store(conn.ConnID(), conn)
	manager.connSize.Add(1)
	logger.Info("connection add to ConnManager successfully: conn num = ", manager.Len())

	return manager
}

// Remove 删除连接
func (manager *ConnManager) Remove(conn ziface.IConnection) error {
	if manager.Len() == 0 {
		return errs.CONN_NOT_FOUND
	}
	manager.connections.Delete(conn.ConnID())
	manager.connSize.Add(-1)
	logger.Info("connection remove from ConnManager successfully: conn num = ", manager.Len())
	return nil
}

// Get 根据connID获取连接
func (manager *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	if manager.Len() == 0 {
		return nil, errs.CONN_NOT_FOUND
	}

	conn, ok := manager.connections.Load(connID)
	if !ok {
		return nil, errs.CONN_NOT_FOUND
	}

	return conn.(ziface.IConnection), nil
}

// Len 获取当前连接数
func (manager *ConnManager) Len() int {
	return int(manager.connSize.Load())
}

// ClearAndStop 删除并停止所有连接
func (manager *ConnManager) ClearAndStop() {
	manager.connections.Range(func(_, value any) bool {
		// 停止
		connection := value.(ziface.IConnection)
		connection.Stop()
		// 删除
		if err := manager.Remove(connection); err != nil {
			return false
		}
		return true
	})
	logger.Info("clear all connections successfully: conn num = ", manager.Len())
}

func NewConnManager() *ConnManager {
	return &ConnManager{}
}
