package core

import (
	"log"
	"net"
	"strconv"
	"erproxy/conf"
)

//Socks5Server .
type Socks5Server struct {
	c conf.InBound
}

// Init for socks5server
func (ss *Socks5Server) Init(c conf.InBound) net.Listener {
	ss.c = c
	return InitServer(ss.c)
}

// Handle for socks5server
func (ss *Socks5Server) Handle(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	// log.Printf("[erproxy] Connected from %v ...", client.RemoteAddr())
	ret := Sock5HandShake(client, ss.c )
	if ret != true {
		log.Printf("handshake failed")
		return
	}

	if isAuth(ss.c) {
		ret = Sock5Auth(client,ss.c)
		if ret == false {
			log.Printf("authenticate failed")
			return 
		}
	}

	ret,out := Socks5Request(client)

	if ret == false {
		log.Println("cannot connect to outbound")
		return
	}

	defer out.close()
	out.loop(client)
}

// Sock5HandShake result, ip, port
func Sock5HandShake(client net.Conn,c conf.InBound) bool {
	var b [1024]byte

	// Read hand shake message
	_, err := client.Read(b[:])
    if err != nil {
        log.Println("Socks Server:",err)
        return false
	}
	
	if b[0] !=  0x05 {
		log.Println("Socks Server: Protocal error or version error")
		return false
	}

	nm := b[1]
	ms := b[2:2+nm]

	if isAuth(c) && selectMethod(ms, 0x02){
		client.Write([]byte{0x05,0x02})
		return true
	} else if !isAuth(c) && selectMethod(ms, 0x00) {
		client.Write([]byte{0x05,0x00})
		return true
	} else {
		client.Write([]byte{0x05,0xFF})
		return false
	}
}

// Sock5Auth auth function
func Sock5Auth(client net.Conn,c conf.InBound) bool {
	var b [1024]byte

	_, err := client.Read(b[:])
	if err != nil {
		log.Println("Socks Server:", err)
		client.Write([]byte{0x01,0x01})
        return false
	}

	v1 := b[0]
	if v1 != 0x01 {
		log.Println("Socks Server:","Autnenticate version error")
		client.Write([]byte{0x01,0x01})
		return false
	}
	un := b[1]
	u := b[2:2+un]
	pn := b[2+un]
	p:= b[2+un+1:2+un+1+pn]
	
	if authenticate(string(u), string(p),c) {
		client.Write([]byte{0x01,0x00})
		return true
	}

	client.Write([]byte{0x01,0x01})
	return false
}

// Socks5Request is main request: ret, cmd, host, port
func Socks5Request(client net.Conn) (bool, Outbound)  {
	var b [1024]byte
	s := []byte{0x05,0x00,0x00,0x01,0x00,0x00,0x00,0x00,0x00,0x00}
	f := []byte{0x05,0x01,0x00,0x01,0x00,0x00,0x00,0x00,0x00,0x00}

	n, err := client.Read(b[:])
	if err != nil {
		log.Println("Socks Server:",err)
		client.Write(f)
        return false, nil
	}

	v := b[0]
	if v != 0x05 {
		log.Println("Socks Server:","Socks version error")
		client.Write(f)
        return false, nil
	}

	cmd := b[1]
	atyp := b[3]
	var host string
	var port string

	port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

	switch(atyp) {
	case 0x01: host = net.IPv4(b[4],b[5],b[6],b[7]).String()
	case 0x03: host = string(b[5 : n-2])
	case 0x04: host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], 
			b[14], b[15], b[16], b[17], b[18], b[19]}.String()
	default:
		log.Println("Socks Server:","Socks version error")
		client.Write(f)
        return false, nil
	}

	var con Outbound
	var ret bool

	switch(cmd) {
	case 0x01: ret,con = Socks5Connect(host, port, atyp)
	}

	if ret {
		client.Write(s)
		return true, con
	}
	client.Write(f)
	return false,con
}

// Socks5Connect connect
func Socks5Connect(host, port string, atype byte) (bool, Outbound) {

	out := getOutBound(host,port,atype)
	if out == nil {
		return false,nil
	}
	ret := out.start(host,port, atype)
	if ret != true {
		return false,nil
	}
	return true,out
}