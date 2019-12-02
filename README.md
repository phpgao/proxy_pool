# proxy_pool

一个简单的代理池工具

A simple proxy pool written in go

## 功能

 - 定时抓取互联网公开免费的代理
 - 定时验证可用代理
 - 支持动态代理(https仅支持connect)
 - 使用采集到的代理访问代理网站
 - 使用命令行环境变量进行配置
 - 当没有IP可用时使用本地转发

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

# 打印可设定参数
./proxy_pool_linux_amd64 --help

# 命令行指定配置
./proxy_pool_linux_amd64 -host 8.8.8.8 -port 6379 -auth laogao

# 后台运行
nohup ./proxy_pool_linux_amd64 > /dev/null 2>&1 &
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

## 一些细节

### 流程图

```
                                                                 +-------------------------+
                                                                 |                         |
                                                                 |                         |
+-------------+      +------------+         +-------------+      |                         |
|             |      |            |         |             |      |                         |
|   source    +------> new proxy  +--------->  validator  +------>                         |
|             +------>            +--------->             +------>                         |
|             |      |            |         |             |      |         The Pool        |
+-------------+      +------------+         +-------------+      |                         |
                                                                 |                         |
                             +----------------------------------->                         |
                             |           +/- score               |                         |
                             |                                   |                         |
                             |                                   |                         |
                             |                                   +-------------+-----------+
                             |                                                 |
                     +-------+------+        +--------------+                  |
                     |              |        |              |                  |
                     |              <--------+              |      cron        |
                     |  old proxy   <--------+  validator   +<-----------------+
                     |              |        |              |
                     |              |        |              |
                     +--------------+        +--------------+

```

### 关于验证逻辑

 1. 代理检测采用打分机制，新代理默认60分，满分100，检测每失败一次扣30分，成功一次加10分，当分数小于等于0时，对应的代理地址将会被删除
 1. 新的代理入库前有三道检测(tcp,http,https)，只要通过了http测试，就会被添加到数据库中
 1. 定时检测只会测试tcp和https的connect方法，同时会把之前判定为http的代理修正为http，但是如果一个https被检测到错误，会扣20分
 1. 目前规则还不算很完善，欢迎大家一起讨论，提高代理的稳定性

### 关于添加采集源

 1. demo见源码source文件夹下
 1. 爬虫分为html、json、re还有text型
 1. 爬虫复用了Spider结构体，新爬虫必须实现的方法如下
 
   - Cron --> 定义了启动间隔
   - Name --> 定义了爬虫名
   - Run -->  通用方法，用来执行下载和解析，照抄即可
   - StartUrl --> 返回目标网站的入口页面
   - Parse --> 接收最终的html代理，返回[]model.HttpProxy实例的指针
 
### 关于动态代理

 1. 核心代码取自[HTTP(S) Proxy in Golang in less than 100 lines of code](https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c)，做了一些针对性的优化
 1. 如果当前没有代理可用，软件会把自身作为透明代理

### 关于传统API代理
 
 1. 接口分为统计和获取
 2. 查询支持schema=http(s)，source=spider.name，score=100，country=cn

## todo

 - tcp池
 - go test
 - 更精细的超时控制
 - 主从模式
 - 代理认证

## 反馈

期待大家的测试和反馈！

## 感谢
 
- https://www.ipip.net/
- https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c
- 感谢给我点star的朋友~