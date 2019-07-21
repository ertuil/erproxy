package core

import (
	"io"
	"net"
	"log"
	"strconv"
	"erproxy/conf"
	"crypto/tls"
)

type sockbound struct {
	server net.Conn
}

func (sb *sockbound) getserver() net.Conn{
	return sb.server
}

func (sb *sockbound) start(host,  port string, atype byte) bool {
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

	sb.server = server
	ret := Sock5Client(sb.server, host, port, atype)
	return ret
}

func (sb *sockbound) loop(client net.Conn){
	go io.Copy(sb.server, client)
	io.Copy(client, sb.server)
}

// Sock5Client the client for sock5
func Sock5Client(server net.Conn, host, port string, atype byte) bool {
	ret := Socks5HandShake(server)
	if ret == false {
		log.Println("can not connect to the next hop")
		return false
	}

	if isOutAuth() {
		ret := Socks5ClientAuth(server)
		if ret == false {
			log.Println("cancel connection")
			return false
		}
	}

	ret = Socks5ClientConnect(server, host, port, atype)
	if ret == false {
		log.Println("cancel connection")
		return false
	}

	return true
}

// Socks5HandShake is handshake 
func Socks5HandShake(server net.Conn) bool {
	var b [1024]byte
	if isOutAuth()  {
		server.Write([]byte{0x05,0x01,0x02})

		_, err := server.Read(b[:])
		if err != nil {
			log.Println(err)
			return false
		}
		if b[0] != 0x05 && b[1] != 0x02{
			log.Println("can not read handshake response")
			return false
		}
		return true
	}

	server.Write([]byte{0x05,0x01,0x00})

	_, err := server.Read(b[:])
	if err != nil {
		log.Println(err)
		return false
	}
	if b[0] != 0x05 && b[1] != 0x00 {
		log.Println("can not read handshake response")
		return false
	}
	return true
}

// Socks5ClientAuth is auth client
func Socks5ClientAuth(server net.Conn) bool {
	b := make([]byte,0)
	var r [1024]byte 
	b = append(b,0x01)
	tu,tp := getOutAuth()
	user := []byte(tu)
	nu := byte(len(user))
	pass := []byte(tp)
	np := byte(len(pass))
	b = append(b,nu)
	b = append(b,user...)
	b = append(b,np)
	b = append(b,pass...)

	server.Write(b)
	_,err := server.Read(r[:])
	if err != nil {
		log.Println(err)
		return false
	}

	if r[0] != 0x01 {
		log.Println("cannot read auth response")
		return false
	}
	
	if r[1] == 0x01 {
		log.Println("username or password error")
		return false
	}

	if r[1] == 0x00 {
		return true
	}
	log.Println("cannot read auth response")
	return false
}

// Socks5ClientConnect lalala
func Socks5ClientConnect(server net.Conn, host, port string, atype byte) bool {
	b := make([]byte,0)
	log.Println("Trying to connect to", host, " : " ,port)
	var ip []byte
	var err error
	if atype == 0x01 || atype == 0x04{
		t := net.ParseIP(host)
		ip, err = t.MarshalText()
		if err != nil {
			log.Println("cannot marshal ip")
			return false
		}
	} else {
		nh := []byte{byte(len(host))}
		ip = []byte(host)
		ip = append(nh,ip...)
	}

	p,err := strconv.Atoi(port)
	if err != nil {
		log.Println("cannot marshal port")
		return false
	}
	var pp  [2]byte
	pp[0] = byte(p / 256)
	pp[1] = byte(p % 256)

	b = append(b,0x05, 0x01, 0x00, atype)
	b = append(b, ip... )
	b = append(b, pp[0],pp[1])

	server.Write(b)
	var rsp [1024]byte
	_,err = server.Read(rsp[:])

	if err != nil {
		log.Println("Error Cannot connect from server")
		return false
	}

	if rsp[0] == 0x05 && rsp[1] == 0x00 {
		return true
	}

	log.Println("Cannot connect from server")
	return false
}