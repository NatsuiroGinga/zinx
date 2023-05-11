package ziface

// IMessage 将请求的一个消息封装到一个Message中，定义抽象的接口
type IMessage interface {
	ID() uint32      // 获取消息的ID
	Data() []byte    // 获取数据的内容
	DataLen() uint32 // 获取数据的长度

	SetID(uint32)      // 设置消息的ID
	SetData([]byte)    // 设置数据的内容
	SetDataLen(uint32) // 设置数据的长度
}
