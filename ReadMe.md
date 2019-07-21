# erproxy

基于 golang 的多级 socks-tls 代理。

## Install

```
make build
```

or 

```
docker pull ertuil/erproxy:latest
docker run -it --rm --name erproxy -p  1080:1080 -v <data>:/app ertuil/erproxy
```
## Usage

```
Usage of ./erproxy-darwin:
  -c string
        set configuration file (default "config.yml")
  -d=true   if erproxy needs to run in the background
```

## Config

example 1 :

``` yaml
log: "log.txt"   # default: stdin
in:
  type: "http" #  or "socks"
  addr: "127.0.0.1" # default: 0.0.0.0
  port: "8080"  # default 1080
  
out:
  type: "socks"
  tls: true
  port: "8081"
  addr: "127.0.0.1"
  auth:
    username: "password"
```

example2:

``` yaml
log: "stdin"
in:
  addr: "127.0.0.1"
  port: "8081"
  tls:
    cert: "xxxxx.cer"
    key: "xxxxx.key"
  auth:
    a: "b"
    c: "d"
  
out:
  type: "free"

routes:
  default: "free"
  a.com: "block"
```

## route

``` yml
routes:
  default: "proxy" # or "direct" or "block"
  route:
    www.baidu.com: "block"
    111.222.333.444: "proxy"
    222.333.444.555/24: "free"
    port:80: "block"
```

## Todo

* socks5: udp associate \ bind command
* vmess

## Notice
在 http/https 代理协议下访问 http 站点请使用强制 connect 方式连接。