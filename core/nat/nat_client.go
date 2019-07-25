package nat

import (
	"crypto/tls"
	"erproxy/conf"
	"errors"
	"io"
	"log"
	"net"
	"time"
)

type NatClient struct {
	name    string
	c       conf.NatClientConf
	OutConn net.Conn
	conConn net.Conn
}

func (nc *NatClient) Init(name string, c conf.NatClientConf) {
	nc.c = c
	nc.name = name
}

func (nc *NatClient) InLink() {
	for {
		if nc.conConn == nil {
			tc, err := nc.NatClientGetConn(nc.c.InAddr, nc.c.InPort, nc.c.InTLS)
			if err != nil {
				log.Println("Nat Client:", err)
			} else {

				ret := nc.NatClientAuth(tc, 0x00)
				if ret == true {
					nc.conConn = tc
				}
				nc.Handle()
			}
		}
	}
}

func (nc *NatClient) NatClientGetConn(host, port string, istls bool) (net.Conn, error) {
	var server net.Conn
	var err error
	serverHost := host
	serverPort := port
	log.Println(host, port)
	if istls == true {
		c := &tls.Config{
			InsecureSkipVerify: true,
		}
		server, err = tls.Dial("tcp", net.JoinHostPort(serverHost, serverPort), c)
	} else {
		server, err = net.Dial("tcp", net.JoinHostPort(serverHost, serverPort))
	}

	if err != nil {
		return nil, err
	}
	return server, nil
}

func (nc *NatClient) NatClientAuth(con net.Conn, isC byte) bool {
	s := []byte{0x01, isC}
	s = append(s, NatAuth(nc.c.Auth)...)
	con.Write(s)

	var b [1024]byte
	n, err := con.Read(b[:])
	log.Println("NAT Client: [debug]", b[:n])
	if err != nil {
		log.Println("Nat Client:", err)
		return false
	}

	if n < 2 || b[0] != 0x01 {
		log.Println("Nat Client: Can not Read Protocal")
		return false
	}
	if b[1] == 0x01 {
		log.Println("Nat Client: Token error")
		return false
	} else if b[1] == 0x02 {
		log.Println("Nat Client: Common Failed")
		return false
	}
	return true
}

func (nc *NatClient) Handle() {

	var b [1024]byte

	for {
		if nc.conConn == nil {
			break
		}

		n, err := nc.conConn.Read(b[:])
		if err != nil {
			log.Println("Nat Client:", err)
			break
		}
		log.Println("NAT Client: [debug]", b[:n])

		if n < 2 || b[0] != versionC {
			log.Println("Nat Client", "Not nat protocal")
		}
		if b[1] == 0x03 {
			inb, outb, err := nc.NewTunnel()
			if err != nil {
				nc.conConn.Write([]byte{0x01, 0x03})
				log.Println("Nat Client:", err)
			} else {
				nc.conConn.Write([]byte{0x01, 0x02})
				go nc.HandleTunnel(inb, outb)
			}
		}
	}
}

func (nc *NatClient) HeartBeat() {
	for {
		if nc.conConn != nil {
			_, ret := nc.conConn.Write([]byte{0x01, 0xFF})
			if ret != nil {
				if nc.conConn != nil {
					nc.conConn.Close()
				}
				nc.conConn = nil
			}
		}
		tm := nc.c.Beat
		if tm == 0 {
			tm = 30
		}
		time.Sleep(time.Second * time.Duration(tm))
	}
}

func (nc *NatClient) NewTunnel() (net.Conn, net.Conn, error) {
	inb, err := nc.NatClientGetConn(nc.c.InAddr, nc.c.InPort, nc.c.InTLS)
	if err != nil {
		return nil, nil, err
	}
	ret := nc.NatClientAuth(inb, 0x01)
	if ret == false {
		inb.Close()
		return nil, nil, errors.New("NAT Client: authenticate failed")
	}

	outb, err := nc.NatClientGetConn(nc.c.OutAddr, nc.c.OutPort, nc.c.OutTLS)
	if err != nil {
		inb.Close()
		return nil, nil, err
	}
	return inb, outb, nil
}

func (nc *NatClient) HandleTunnel(inb, outb net.Conn) {
	defer inb.Close()
	defer outb.Close()

	go io.Copy(inb, outb)
	io.Copy(outb, inb)
}
