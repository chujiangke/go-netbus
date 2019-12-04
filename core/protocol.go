package core

import (
	"../util"
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	// 协议长度=字段数量
	protocolLength = 3

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
// 结果|消息类型|访问端口|Key
// 0|0|7001|customabcd8000

// 协议
type Protocol struct {
	Result int // 结果：0 失败，1 成功
	Port   int // 端口
	Key    string
}

// 转字符串
func (p *Protocol) String() string {
	return fmt.Sprintf("%d|%d|%s", p.Result, p.Port, p.Key)
}

func (p *Protocol) NewResult(newResult int) Protocol {
	return Protocol{
		Result: newResult,
		Port:   p.Port,
		Key:    p.Key,
	}
}

// 解析协议
func _parseProtocol(body string) (Protocol, bool) {
	// 拆解字符
	arr := strings.Split(body, "|")
	if len(arr) != protocolLength {
		return Protocol{Result: protocolResultFailToParse}, false
	}
	// 前两个字段
	params, err := util.AtoInt(arr[0:2])
	if err != nil {
		return Protocol{Result: protocolResultFailToParse}, false
	}
	return Protocol{
		Result: params[0],
		Port:   params[1],
		Key:    arr[2],
	}, true
}

// 发送协议
// 第一个字节为协议长度
// 协议长度只支持到255
func sendProtocol(conn net.Conn, req Protocol) bool {
	buffer := bytes.NewBuffer([]byte{})
	length := byte(len(req.String()))

	// 协议长度
	buffer.WriteByte(length)
	// 协议内容
	buffer.WriteString(req.String())

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		log.Printf("Send protocol failed. [%s] %s\n", req.String(), err.Error())
		return false
	}
	//log.Println("Send protocol", req.String())
	return true
}

// 接收协议
// 第一个字节为协议长度
func receiveProtocol(conn net.Conn) (Protocol, bool) {
	var err error

	// 读取协议长度
	buffer := make([]byte, 1)
	if _, err := conn.Read(buffer); err != nil {
		log.Println("Parse protocol length failed.", err.Error())
		return Protocol{Result: protocolResultFailToReceive}, false
	}
	// 读取协议内容
	buffer = make([]byte, buffer[0])
	if _, err = conn.Read(buffer); err != nil {
		log.Println("Parse protocol body failed.", err.Error())
		return Protocol{Result: protocolResultFailToReceive}, false
	}
	// 解析消息
	body := strings.TrimSpace(string(buffer))
	//log.Println("----------> Receive protocol", body)

	return _parseProtocol(body)
}
