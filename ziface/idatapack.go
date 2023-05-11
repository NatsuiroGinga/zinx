package ziface

// IDataPack 封包、拆包 模块
type IDataPack interface {
	HeadLen() uint32                              // 获取包头长度方法
	Pack(msg IMessage) (data []byte, err error)   // 封包方法
	Unpack(data []byte) (msg IMessage, err error) // 拆包方法
}
