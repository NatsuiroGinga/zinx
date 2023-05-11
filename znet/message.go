package znet

// Message 将请求的一个消息封装到一个Message中
type Message struct {
	id      uint32 // 消息的ID
	dataLen uint32 // 消息中数据部分的长度
	data    []byte // 消息的内容
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{id: id, data: data, dataLen: uint32(len(data))}
}

func (msg *Message) ID() uint32 {
	return msg.id
}

func (msg *Message) DataLen() uint32 {
	return msg.dataLen
}

func (msg *Message) Data() []byte {
	return msg.data
}

func (msg *Message) SetID(id uint32) {
	msg.id = id
}

func (msg *Message) SetData(data []byte) {
	msg.data = data
}

func (msg *Message) SetDataLen(dataLen uint32) {
	msg.dataLen = dataLen
}
