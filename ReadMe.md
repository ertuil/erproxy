# erproxy

基于 golang 的多级 socks-tls 代理。

## Usage

```
Usage of ./erproxy-darwin:
  -c string
        set configuration file (default "config.yml")
  -d=true   if erproxy needs to run in the background
  -l string
        set logging file (default "erproxy.log")
```

## Config

example 1 :

``` yaml
in:
  addr: "127.0.0.1"
  port: "8080"
  
out:
  type: "sock"
  tls: true
  port: "8081"
  addr: "127.0.0.1"
  auth:
    user: "xxxx"
    token: "xxxxxx"
```

example2:

``` yaml
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

* route: freebound or next hop
* socks5: udp associate \ bind command
* accounting
* http/https/vmess