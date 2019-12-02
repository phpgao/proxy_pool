# proxy_pool

一个简单的代理池工具

A simple proxy pool written in go

## 功能

 - 定时抓取互联网公开免费的代理
 - 定时验证可用代理
 - 支持动态代理(https仅支持connect)
 - 使用采集到的代理访问代理网站
 - 使用命令行环境变量进行配置
 - 当没有IP可用时使用本地转发 - new

## 依赖

 - redis

## 使用说明

### 编译

```bash
go build

cp config_example.json config.json
# 修改redis和端口配置

# 感谢ipip.net提供精准的IP数据(已内置)
./proxy_pool

# 帮助
./proxy_pool_linux_amd64 --help

./proxy_pool_linux_amd64 -ip 8.8.8.8 -port 6379 -auth laogao
```

### api

```bash
# 统计
curl 127.0.0.1:8088
# 随机
curl 127.0.0.1:8088/random
# 获取列表
curl 127.0.0.1:8088/get
```

### 动态代理

```bash
# http
curl http://cip.cc -x 127.0.0.1:8089
# https
curl https://cip.cc -x 127.0.0.1:8089
```

## todo

 - tcp池
 - go test
 - 更精细的超时控制
 - 主从模式
 - 代理认证

## 反馈

我们 issue 见
