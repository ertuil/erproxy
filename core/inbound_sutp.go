package core

import (
	"erproxy/conf"
	"erproxy/header"
	"log"
	"net"
	"strconv"
)

// SUTPServer is sutp server
type SUTPServer struct {
	c    conf.InBound
	name string
}

// Init for SUTPServer
func (ss *SUTPServer) Init(name string, c conf.InBound) net.Listener {
	ss.name = name
	ss.c = c
	return InitServer(ss.c)
}

// Handle for SUTPServer
func (ss *SUTPServer) Handle(client net.Conn) {
	defer client.Close()
	ret, ad := ss.HandShake(client)
	if ret == false {
		return
	}

	out := getOutBound(ad)
	ret = out.start(ad)
	if ret == false {
		log.Println("SUTP Server: Cannot connect to next hop")
		client.Write([]byte{0x01, 0x01})
		return
	}
	defer out.close()
	out.loop(client)
}

// HandShake for SUTPServer
func (ss *SUTPServer) HandShake(client net.Conn) (bool, header.AddrInfo) {
	successcontent := []byte{0x01, 0x00}
	failedcontent := []byte{0x01, 0x00}

	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println("SUTP Server:", err)
		client.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return false, header.AddrInfo{}
	}

	ret, key, iv := ss.CheckToken(b[:])
	if ret == false {
		log.Println("SUTP Server: auth failed")
		client.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return false, header.AddrInfo{}
	}
	if len(b) <= 18 {
		return false, header.AddrInfo{}
	}
	msg, err := decryptSession(key, iv, b[18:n])
	if err != nil {
		log.Println("SUTP Server:", err)
		client.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return false, header.AddrInfo{}
	}

	n = len(msg)
	if n <= 6 || msg[0] != 0x01 {
		log.Println("SUTP Server: version error")
		client.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return false, header.AddrInfo{}
	}

	cmd := msg[1]
	atyp := msg[3]
	var host string
	var port string

	port = strconv.Itoa(int(msg[n-2])<<8 | int(msg[n-1]))

	switch atyp {
	case 0x01:
		host = net.IPv4(msg[4], msg[5], msg[6], msg[7]).String()
	case 0x03:
		host = string(msg[5 : n-2])
		atyp = header.HostCheck(host)
	case 0x04:
		host = net.IP{msg[4], msg[5], msg[6], msg[7], msg[8], msg[9], msg[10], msg[11], msg[12], msg[13],
			msg[14], msg[15], msg[16], msg[17], msg[18], msg[19]}.String()
	default:
		client.Write(failedcontent)
		log.Println("SUTP Server: Parse addr error")
		return false, header.AddrInfo{}
	}

	ad := header.AddrInfo{}
	ad.SetInfo(ss.name, host, port, atyp, cmd)

	client.Write(successcontent)
	return true, ad
}

func (ss *SUTPServer) CheckToken(b []byte) (bool, []byte, []byte) {
	fiv := b[:16]
	for _, tk := range ss.c.Auth {
		key, iv := getKeyIV(tk, fiv)
		if iv[0] == b[16] && iv[1] == b[17] {
			return true, key, iv
		}
	}
	key, iv := getKeyIV("erproxy", fiv)
	return true, key, iv
}
