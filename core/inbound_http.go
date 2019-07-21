package core

import (
	"log"
	"net/url"
	"strings"
	"net"
)

// FakeHandle .
func FakeHandle(client net.Conn,info interface{}) {
	log.Println("HTTP server:", info)
	client.Write([]byte("HTTP/1.1 200 OK \r\n\r\nhello,world!"))
}

func parseInfos(rawurl string) (string,string,byte) {

	if len(rawurl) <= 7 || rawurl[0:7] != "http://" {
		rawurl = "http://" + rawurl
	}

	u,err := url.Parse(rawurl)
	if err != nil {
		return "","",0x00
	}

	host := u.Hostname()
	port := u.Port()
	var atype byte= 0x01
	if net.ParseIP(host) != nil && len(host) > 15{
		atype = 0x04
	} else {
		atype = 0x03
	}
	log.Println(host,port)
	return host,port,atype
}

// HTTPServerHandle .
func HTTPServerHandle(client net.Conn) {
	var b [1024]byte

	_,err := client.Read(b[:])
	if err != nil {
		FakeHandle(client, err)
		return
	}
	rawstr := string(b[:])
	strs := strings.Split(rawstr,"\n")
	if len(strs) <= 0 {
		FakeHandle(client, "2Cannot read from server")
		return 
	}

	words := strings.Split(strs[0]," ")
	if len(words) < 2 {
		FakeHandle(client, "Cannot read from server")
		return 
	} 

	log.Println(strs[0])

	host,port,atype := parseInfos(words[1])
	if atype == 0x00 {
		FakeHandle(client, "Cannot read url")
		return
	}

	if strings.ToLower(words[0]) != "connect" {
		FakeHandle(client,"error method")
		return
	}

	ret := true
	if isAuth() {
		ret = false
		for _,v := range(strs) {
			if len(v) > 8 && strings.ToLower(v[0:8]) == "proxy-au" {
				token := strings.Split(v," ")
				token[2] = strings.TrimSuffix(token[2],"\r")
				if len(token) < 3 {
					FakeHandle(client,"Error proxy auth")
					return
				}
				ret = HTTPAuth(token[2])
				break
			}
		}
	}

	if ret == false {
		log.Println("HTTP Server: need to auth")
		client.Write([]byte("HTTP/1.1 407 Proxy authentication required\r\nProxy-Authenticate: Basic realm=erproxy\r\n"))
	}

	ob := getOutBound(host,port,atype)
	ret = ob.start(host,port,atype)
	if ret == true {
		client.Write([]byte("HTTP/1.1 200 Connection established\r\nProxy-agent: erproxy\r\n\r\n"))
		ob.loop(client)
	}
	client.Write([]byte("HTTP/1.1 404 Not found\r\n\r\n"))
}
