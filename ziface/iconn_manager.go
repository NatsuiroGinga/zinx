package ziface

// IConnManager 连接管理模块抽象层
type IConnManager interface {
	Add(conn IConnection) IConnManager      // 添加连接
	Remove(conn IConnection) error          // 删除连接
	Get(connID uint32) (IConnection, error) // 根据connID获取连接
	Len() int                               // 获取当前连接数
	ClearAndStop()                          // 删除并停止所有连接
}
