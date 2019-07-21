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

// HTTPServerHandle .
func HTTPServerHandle(client net.Conn) {
	var b [1024]byte

	_,err := client.Read(b[:])
	if err != nil {
		FakeHandle(client, "Cannot read from server")
		return
	}
	rawstr := string(b[:])
	strs := strings.Split(rawstr,"\n")
	if len(strs) <= 0 {
		FakeHandle(client, "Cannot read from server")
		return 
	}

	words := strings.Split(strs[0]," ")
	if len(words) < 2 {
		FakeHandle(client, "Cannot read from server")
		return 
	} 

	if words[0] == "connect" || words[0] == "CONNECT" {
		FakeHandle(client, "Method Error")
		return
	}

	var rawurl string
	if len(words[1]) <= 5 || words[1][0:4] != "http" {
		rawurl = "http://" + words[1]
	} else {
		rawurl = words[1]
	}

	u,err := url.Parse(rawurl)
	if err != nil {
		FakeHandle(client,"Cannot parse host")
	}

	host := u.Hostname()
	port := u.Port()
	var atype byte= 0x01
	ret := true
	if isAuth() {
		ret = false
		for _,v := range(strs) {
			if len(v) > 5 && strings.ToLower(v[0:5]) == "proxy" {
				token := strings.Split(v," ")
				if len(token) < 3 {
					FakeHandle(client,"Error proxy auth")
				}
				ret = HTTPAuth(token[2])
			}
		}
	}

	if ret == false {
		log.Println("HTTP Server: need to auth")
		client.Write([]byte("HTTP/1.1 401 Unauthorized\r\nProxy-Authenticate: Basic realm=erproxy\r\n"))
	}

	if net.ParseIP(host) != nil && len(host) > 15{
		atype = 0x04
	} else {
		atype = 0x03
	}
	ob := getOutBound(host,port,atype)
	ret = ob.start(host,port,atype)
	if ret == true {
		client.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		ob.loop(client)
	}
	client.Write([]byte("HTTP/1.1 404 Not found\r\n\r\n"))
}
