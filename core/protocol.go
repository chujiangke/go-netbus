package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

const (
	// 协议长度=字段数量
	protocolByteMinLength = 15
	protocolByteMaxLength = 25

	// 协议-结果
	protocolResultSuccess       = 0 // 成功，默认值
	protocolResultFail          = 1 // 失败
	protocolResultFailToSend    = 2 // 发送失败
	protocolResultFailToReceive = 3 // 接收失败
	protocolResultFailToParse   = 4 // 解析失败
	protocolResultFailToAuth    = 5 // 鉴权失败
	protocolResultFailToListen  = 6 // 监听失败

	// Key 长度
	protocolKeyMinLength = 6  // Key 最小长度
	protocolKeyMaxLength = 16 // Key 最大长度
)

// 协议格式
// 结果|消息类型|原端口|访问端口|Key
// 0|0|3306|13306|winshu

// 协议
type Protocol struct {
	Result     byte   // 结果：0 失败，1 成功
	AccessPort uint32 // 访问端口
	Port       uint32 // 原端口
	Key        string // 身份验证
}

// 转字符串
func (p *Protocol) String() string {
	return fmt.Sprintf("%d|%d|%d|%s", p.Result, p.Port, p.AccessPort, p.Key)
}

// 返回一个新结果
func (p *Protocol) NewResult(newResult byte) Protocol {
	return Protocol{
		Result:     newResult,
		Port:       p.Port,
		AccessPort: p.AccessPort,
		Key:        p.Key,
	}
}

func (p *Protocol) Bytes() []byte {
	buffer := bytes.NewBuffer([]byte{})

	buffer.WriteByte(p.Result)
	_ = binary.Write(buffer, binary.BigEndian, p.Port)
	_ = binary.Write(buffer, binary.BigEndian, p.AccessPort)
	buffer.WriteString(p.Key)
	return buffer.Bytes()
}

func (p *Protocol) Len() byte {
	return byte(len(p.Bytes()))
}

// 解析协议
func _parseProtocol(body []byte) (Protocol, bool) {
	// 长度异常
	if len(body) < protocolByteMinLength || len(body) > protocolByteMaxLength {
		log.Println("Fail to parse protocol length")
		return Protocol{Result: protocolResultFailToParse}, false
	}

	protocol := Protocol{
		Result:     body[0],
		Port:       binary.BigEndian.Uint32(body[1:5]),
		AccessPort: binary.BigEndian.Uint32(body[5:9]),
		Key:        string(body[9:]),
	}
	log.Println("Parse Protocol", protocol.String())
	return protocol, true
}

// 发送协议
// 第一个字节为协议长度
// 协议长度只支持到255
func sendProtocol(conn net.Conn, req Protocol) bool {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteByte(req.Len())
	buffer.Write(req.Bytes())

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		log.Printf("Send protocol failed. [%s] %s\n", req.String(), err.Error())
		return false
	}
	//log.Println("Send protocol", req.String())
	return true
}

// 接收协议
// 第一个字节为协议长度
func receiveProtocol(conn net.Conn) (protocol Protocol, ok bool) {
	var err error
	var length byte

	if err = binary.Read(conn, binary.BigEndian, &length); err != nil {
		log.Println("Parse protocol length failed.", err.Error())
		return Protocol{Result: protocolResultFailToReceive}, false
	}
	// 读取协议内容
	body := make([]byte, length)
	if err = binary.Read(conn, binary.BigEndian, &body); err != nil {
		log.Println("Parse protocol body failed.", err.Error())
		return Protocol{Result: protocolResultFailToReceive}, false
	}
	return _parseProtocol(body)
}
