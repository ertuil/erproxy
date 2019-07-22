# erproxy

基于 golang 的多级 socks-tls、http(s)、sutp（自有协议） 代理。

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

Example

``` yaml
log: "stdin"

in:  # Inbounds
  a:
    type: "http" # http or socks
    addr: "127.0.0.1"
    port: "8080"
  b:
    type: "socks"
    addr: "127.0.0.1"
    port: "8081"
    tls:
      cert: "fullchain.cer"
      key: "xxxx.key"
    auth:
      aaa: "bbb"
      c: "d"
  
out:
  c:
    type: "http"
    tls: true
    port: "29980"
    addr: "light.ustclug.org"
    auth:
      "xxxx": "xxxx"
  d:
    type: "socks"
    tls: true
    port: "465"
    addr: "japan.ertuil.top"
    auth:
      "xxxx": "xxxx"
  free:
    type: "free"
  block:
    type: "block"

routes:
  default: "c"
  route:
    baidu.com: "free"
    google: "d"
    twitter: "block"
```

## Route 

```
routes:
  default: "d" #
  route:
    www.baidu.com: "free" # daemon
    111.222.333.444: "free" # IPv4 or IPv6
    222.333.444.555/24: "block" # CIDR
    port:80: "c" # PORT
    github@a: "block" # from a to github via block, aka drop
    github@b: "c" # from b to github via c
```

## Todo

* socks5: udp associate \ bind command

## Notice
在 http/https 代理协议下访问 http 站点请使用强制 connect 方式连接。

## SUTP

一种自由代理协议，实现 1RTT 完成信道建立可能比 socks 具有更好性能和较高混淆特性。