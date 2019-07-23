package core

import (
	"crypto/tls"
	"erproxy/conf"
	"erproxy/header"
	"io"
	"log"
	"net"
	"strconv"
)

type sutpbound struct {
	name   string
	server net.Conn
	c      conf.OutBound
}

func (sb *sutpbound) getserver() net.Conn {
	return sb.server
}

func (sb *sutpbound) close() {
	sb.server.Close()
}

func (sb *sutpbound) init(name string, c conf.OutBound) {
	sb.name = name
	sb.c = c
}

func (sb *sutpbound) loop(client net.Conn) {
	go io.Copy(sb.server, client)
	io.Copy(client, sb.server)
}

func (sb *sutpbound) start(ad header.AddrInfo) bool {
	serverHost, serverPort := sb.c.Addr, sb.c.Port

	var server net.Conn
	var err error

	if sb.c.UseTLS == true {
		c := &tls.Config{
			InsecureSkipVerify: true,
		}
		server, err = tls.Dial("tcp", net.JoinHostPort(serverHost, serverPort), c)
	} else {
		server, err = net.Dial("tcp", net.JoinHostPort(serverHost, serverPort))
	}

	if err != nil {
		log.Println(err)
		return false
	}

	sb.server = server
	_, host, port, atype, cmd := ad.GetInfo()

	var ip []byte
	if atype == 0x01 {
		t := net.ParseIP(host)
		if t == nil {
			log.Println("SUTP Client: Cannot marshal ip")
			server.Close()
			return false
		}
		ip = t.To4()
		if ip == nil {
			log.Println("SUTP Client: Cannot marshal ip")
			server.Close()
			return false
		}
	} else if atype == 0x04 {
		t := net.ParseIP(host)
		if t == nil {
			log.Println("SUTP Client: Cannot marshal ip")
			server.Close()
			return false
		}
		ip = t.To16()
		if ip == nil {
			log.Println("SUTP Client: Cannot marshal ip")
			server.Close()
			return false
		}
	} else {
		nh := []byte{byte(len(host))}
		ip = []byte(host)
		ip = append(nh, ip...)
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		log.Println("SUTP Client: Can not marshal port")
		server.Close()
		return false
	}
	var pp [2]byte
	pp[0] = byte(p / 256)
	pp[1] = byte(p % 256)

	msg := []byte{0x01, cmd, 0x00, atype}
	msg = append(msg, ip...)
	msg = append(msg, pp[0], pp[1])

	_, tk := getOutAuth(sb.c)
	if tk == "" {
		tk = "erproxy"
	}
	fiv := initIV()
	key, iv := getKeyIV(tk, fiv)
	encrypted := encryptSession(key, iv, msg)
	b := append(fiv, iv[:2]...)
	b = append(b, encrypted...)

	server.Write(b)
	var resp [1024]byte
	_, err = server.Read(resp[:])
	if err != nil {
		log.Println("SUTP Client:", err)
	}

	if len(resp) > 2 && resp[0] == 0x01 && resp[1] == 0x00 {
		return true
	}
	log.Println("SUTP Client: Cannot connect to server")
	server.Close()
	return false
}
