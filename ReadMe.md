# erproxy

基于 golang 的多级 socks-tls 代理。

## Usage

```
  -c string
        set configuration file (default "config.yml")
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
    user: "xxxxxx"
    token: "xxxxxx"
  
out:
  type: "free"
```

## Todo

* route: freebound or next hop
* socks5: udp associate \ bind command
* accounting
* http/https/vmess