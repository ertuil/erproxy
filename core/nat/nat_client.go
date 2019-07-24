package nat

import (
	"crypto/tls"
	"erproxy/conf"
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

func (nc *NatClient) Link() {
	for {
		tc, err := nc.NatClientGetConn(nc.c.InAddr, nc.c.InPort, nc.c.InTLS)
		if err != nil {
			log.Println("Nat Client:", err)
		} else {
			log.Println("Nat Client: [debug]", tc)
			defer func() {
				if nc.conConn != nil {
					nc.conConn.Close()
				}
				nc.conConn = nil
			}()
			ret := nc.NatClientAuth(tc, 0x00)
			if ret == true {
				nc.conConn = tc
			}
			nc.Handle()
		}
		time.Sleep(3 * time.Second)
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
	tc := nc.conConn
	if tc == nil {
		return
	}
	var b [1024]byte

	for {
		if tc == nil {
			break
		}
		n, err := tc.Read(b[:])
		if err != nil {
			log.Println("Nat Client:", err)
			break
		}

		if n < 2 || b[0] != 0x01 {
			log.Println("Nat Client", "Not nat protocal")
		}
		if b[1] == 0x03 {
			go nc.NewTunnel()
		}
	}
}

func (nc *NatClient) HeartBeat() {
	for {
		if nc.conConn != nil {
			nc.conConn.Write([]byte{0x01, 0x04})
		}
		time.Sleep(time.Second * 8)
	}
}

func (nc *NatClient) NewTunnel() {
	a, err := nc.NatClientGetConn(nc.c.InAddr, nc.c.InPort, nc.c.InTLS)
	if err != nil {
		log.Println("Nat Client:", err)
		return
	}
	b, err := nc.NatClientGetConn(nc.c.OutAddr, nc.c.OutPort, nc.c.OutTLS)
	if err != nil {
		log.Println("Nat Client:", err)
		return
	}
	nc.NatClientAuth(a, 0x01)

	defer a.Close()
	defer b.Close()

	go io.Copy(a, b)
	io.Copy(b, a)
}

//func (nc *NatClient) Init() bool {
//	tc,err :=  net.Dial("tcp",net.JoinHostPort(nc.c.InAddr, nc.c.InPort))
//	if err != nil {
//		log.Println("Nat Client:",err)
//		return false
//	}
//
//	s := []byte{0x01,0x00}
//	s = append(s,NatAuth(nc.c.Auth)...)
//	tc.Write(s)
//
//	var b [1024]byte
//	n,err := tc.Read(b[:])
//	if err != nil  {
//		log.Println("Nat Client:",err)
//		return false
//	}
//
//	if n < 2 || b[0] != 0x01 {
//		log.Println("Nat Client: Can not Read Protocal")
//		return false
//	}
//	if b[1] == 0x01 {
//		log.Println("Nat Client: Token error")
//		return false
//	} else if b[1] == 0x02{
//		log.Println("Nat Client: Common Failed")
//		return false
//	} else {
//		nc.conConn = tc
//		return true
//	}
//
//	return true
//}
//
//func (nc *NatClient) Handle() {
//	var b [1024]byte
//	_,err := tc.Read(b[:])
//	if err != nil  {
//		log.Println("Nat Client:",err)
//	}
//
//}
