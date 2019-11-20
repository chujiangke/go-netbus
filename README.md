# 项目介绍

**go-netbus**【网络直通车】是为了解决内网穿透问题而建立的项目。

## 功能列表

- 基于 TCP 协议
- 支持无限重连机制、有限次数重连机制
- 支持随机端口、偏移端口、同名端口三种代理模式

## 启动方式

> 启动前，注意`config.ini`在不同平台下**换行符**的问题

支持两种方式启动:

- 纯命令行启动（不需要配置文件）
- 配置文件启动（需要配置文件）**【推荐】**

### 命令行启动

```bash
# 启动服务端

$ netbus -server <port> [port-mode]

# 注释
# port      服务端端口，不要使用保留端口，必填
# port-mode 代理端口模式，支持三种方式，可选
#      0    使用随机端口(60000+) 默认方式
#      1    使用同名端口
#      2    偏移端口(比如被代理端口是 3000，代理端口就是 4000)

```

```bash
# 启动客户端

$ netbus -client <server:port> <local:port>

# 注释
# server:port  服务端地址，格式如：45.32.78.129:6666
# local:port   被代理服务地址，多个以逗号隔开，格式如：127.0.0.1:8080,127.0.0.1:9200
```

### 配置文件启动

配置文件`config.ini`需与启动文件置于同一目录。

通常情况下，服务端配置与客户端配置是分开的。

**服务端配置**
```ini
[server]
# 桥接端口
port = 6666
# 端口模式：0=使用随机端口(60000+)，1=同名端口，2=偏移端口(比如被代理端口是 3000，代理端口就是 4000)
port-mode = 1
```

**客户端配置**
```ini
[client]
# 远程代理地址，端口必须一致
server-host = 127.0.0.1:6666
# 内网被代理服务地址
local-host = 127.0.0.1:3306
# 最大重试次数，-1表示一直重试
max-redial-times = 20
```

**启动命令**
```bash
# 启动服务端
$ netbus -server

# 启动客户端
$ netbus -client
```