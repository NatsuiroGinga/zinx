package znet

import (
	"bytes"
	"encoding/binary"
	"zinx/config"
	errs "zinx/lib/enum/err"
	"zinx/lib/logger"
	"zinx/ziface"
)

var dataPack *DataPack

func init() {
	dataPack = &DataPack{}
}

// DataPack 封包、拆包 模块
type DataPack struct {
}

// NewDataPack 返回一个DataPack单例
func NewDataPack() *DataPack {
	return dataPack
}

func (dataPack *DataPack) HeadLen() uint32 {
	// id uint32(4字节) + dataLen uint32(4字节)
	return 4 + 4
}

func (dataPack *DataPack) Pack(msg ziface.IMessage) (data []byte, err error) {
	// 1. 创建一个存放bytes字节的缓冲
	buf := new(bytes.Buffer)
	// 2. 写dataLen
	if err = binary.Write(buf, binary.LittleEndian, msg.DataLen()); err != nil {
		logger.Error("binary.Write(dataLen) error: ", err)
		return
	}
	// 3. 写msgID
	if err = binary.Write(buf, binary.LittleEndian, msg.ID()); err != nil {
		logger.Error("binary.Write(msgID) error: ", err)
		return
	}
	// 4. 写data数据
	if err = binary.Write(buf, binary.LittleEndian, msg.Data()); err != nil {
		logger.Error("binary.Write(data) error: ", err)
		return
	}

	return buf.Bytes(), nil
}

func (dataPack *DataPack) Unpack(data []byte) (msg ziface.IMessage, err error) {
	// 1. 创建一个从输入二进制数据的ioReader
	reader := bytes.NewReader(data)
	// 2. 只解压head信息，得到dataLen和msgID
	var dataLen, msgId uint32
	if err = binary.Read(reader, binary.LittleEndian, &dataLen); err != nil {
		logger.Error("binary.Read(dataLen) error: ", err)
		return
	}
	if err = binary.Read(reader, binary.LittleEndian, &msgId); err != nil {
		logger.Error("binary.Read(msgId) error: ", err)
		return
	}
	// 3. 判断dataLen是否超出了我们允许的最大包长度
	if config.ZinxProperties.MaxPackageSize > 0 && dataLen > config.ZinxProperties.MaxPackageSize {
		logger.Error(errs.TOO_LARGE_PACKAGE)
		return nil, errs.TOO_LARGE_PACKAGE
	}
	// 4. 封装到Message中
	msg = &Message{
		id:      msgId,
		dataLen: dataLen,
	}

	return
}
