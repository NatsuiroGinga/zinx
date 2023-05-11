package znet

import (
	"log"
	"net"
	"sync"
	"testing"
)

func TestDataPack_Pack(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	// 1. TCP
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 7777,
	})
	if err != nil {
		t.Error(err)
		return
	}
	// 2. 模拟服务器
	go func() {
		defer func() {
			wg.Done()
			_ = listener.Close()
		}()

		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			panic(err)
		}
		// 读取客户端的请求数据
		if err != nil {
			panic(err)
		}

		go func() {
			defer func() {
				wg.Done()
				_ = tcpConn.Close()
			}()
			for {
				header := make([]byte, dataPack.HeadLen())
				n, err := tcpConn.Read(header)
				if err != nil || n == 0 {
					log.Println(err)
					break
				}
				msg, err := dataPack.Unpack(header)
				if err != nil {
					panic(err)
				}
				if msg.DataLen() > 0 {
					data := make([]byte, msg.DataLen())
					_, err := tcpConn.Read(data)
					if err != nil {
						panic(err)
					}
					msg.SetData(data)
					t.Log("===> Recv Msg: ID=", msg.ID(), ", len=", msg.DataLen(), ", data=", string(msg.Data()))
				}
			}
		}()
	}()
	// 3. 模拟客户端
	tcpConn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 7777,
	})
	defer tcpConn.Close()
	if err != nil {
		t.Error(err)
		return
	}
	// 4. 封包
	bytes := []byte("Hello Zinx V0.1")
	msg1 := NewMessage(1, bytes)
	msg2 := NewMessage(2, []byte("Hello Zinx V0.2"))

	bytes, err = dataPack.Pack(msg1)
	data, err := dataPack.Pack(msg2)
	// 5. 发送
	if _, err := tcpConn.Write(append(bytes, data...)); err != nil {
		t.Error(err)
		return
	}

	wg.Wait()
}
