package core

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"strconv"
	"strings"
)

type ServerConfig struct {
	Port       int
	RandomPort bool
}

type ClientConfig struct {
	ServerAddr     string
	LocalAddr      string
	MaxRedialTimes int
}

func (t ClientConfig) GetLocalAddr() []string {
	str := strings.ReplaceAll(t.LocalAddr, " ", "")
	return strings.Split(str, ",")
}

type NetAddress struct {
	IP   string
	Port int
}

func (t NetAddress) String() string {
	return fmt.Sprintf("%s:%d", t.IP, t.Port)
}

// 解析地址
func ParseNetAddress(host string) NetAddress {
	arr := strings.Split(host, ":")
	if len(arr) != 2 {
		return NetAddress{}
	}
	port, err := strconv.Atoi(strings.TrimSpace(arr[1]))
	if err != nil {
		port = 0
	}
	return NetAddress{strings.TrimSpace(arr[0]), port}
}

var (
	serverConfig ServerConfig
	clientConfig ClientConfig
)

func loadConfig() *ini.File {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("Fail to load config", err.Error())
	}
	return cfg
}

func InitServerConfig() ServerConfig {
	cfg := loadConfig()
	server := func(key string) *ini.Key {
		return cfg.Section("server").Key(key)
	}
	port, _ := server("port").Int()
	randomPort, _ := server("random-port").Bool()
	serverConfig = ServerConfig{Port: port, RandomPort: randomPort}

	log.Println("Init server config finished", serverConfig)
	return serverConfig
}

func InitClientConfig() ClientConfig {
	cfg := loadConfig()
	client := func(key string) *ini.Key {
		return cfg.Section("client").Key(key)
	}
	serverAddr := client("server-host").String()
	localAddr := client("local-host").String()
	maxRedialTimes, err := client("max-redial-times").Int()
	if err != nil {
		maxRedialTimes = 20
	}
	clientConfig = ClientConfig{ServerAddr: serverAddr, LocalAddr: localAddr, MaxRedialTimes: maxRedialTimes}
	log.Println("Init client config finished", clientConfig)
	return clientConfig
}
