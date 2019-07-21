package core

import (
	"io"
	"net"
	"log"
	"strings"
	"erproxy/conf"
	"crypto/tls"
	"encoding/base64"
)

type httpbound struct {
	server net.Conn
}

func (hb *httpbound) getserver() net.Conn{
	return hb.server
}

func (hb *httpbound) start(host,  port string, atype byte) bool {

	var server net.Conn
	var err error

	serverHost := conf.CC.OutBound.Addr
	serverPort := conf.CC.OutBound.Port

	if conf.CC.OutBound.UseTLS == true {
		c := &tls.Config{
			InsecureSkipVerify: true,
		}
		server, err = tls.Dial("tcp", net.JoinHostPort(serverHost,serverPort), c)
	} else {
		server, err = net.Dial("tcp", net.JoinHostPort(serverHost,serverPort))
	}

	if err != nil {
		log.Println(err)
		return false
	}

	hb.server = server

	str := "CONNECT " + net.JoinHostPort(host,port)  + "\r\nUser-agent: erproxy\\0.0.4\r\n"
	if isOutAuth() {
		user,token := getOutAuth()
		str += "Proxy-authorization: Basic " + base64.URLEncoding.EncodeToString([]byte(user+":"+token))+"\r\n\r\n"
	}
	log.Println("HTTP Client:",str)
	server.Write([]byte(str))
	
	var b [1024]byte
	_,err = server.Read(b[:])
	if err != nil  {
		log.Println("HTTP client: Cannot read from server")
		return false
	}
	strs := strings.Split(string(b[:]),"\r\n")
	if len(strs) <= 0 {
		log.Println("HTTP client: Cannot read from server")
		return false
	}

	words := strings.Split(strs[0]," ")
	if len(words) < 2 {
		log.Println("HTTP client: Cannot read status code")
		return false
	} 
	if words[1] == "200"{
		log.Println("HTTP client: Try to connect to ",net.JoinHostPort(host,port))
		return true
	} else if words[1] == "407" {
		log.Println("HTTP client: Proxy Authentication Required")
		return false
	}
	log.Println("HTTP Client: error:",words[1:])
	return false
}

func (hb *httpbound) loop(client net.Conn){
	go io.Copy(hb.server, client)
	io.Copy(client, hb.server)
}

func (hb *httpbound) close(){
	if hb.server != nil {
		hb.server.Close()
	}
}
