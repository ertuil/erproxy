package nat

import (
	"crypto/tls"
	"erproxy/conf"
	"io"
	"log"
	"net"
)

const (
	versionC     = 0x01
	defaultPassC = "erproxy"
)

type NatServer struct {
	name        string
	c           conf.NatServerConf
	conConn     net.Conn
	inConn      net.Conn
	InListener  net.Listener
	OutListener net.Listener
}

//  Init  for NatServer
func (ns *NatServer) Init(name string, c conf.NatServerConf) error {
	ns.name = name
	ns.c = c
	var err error
	ns.InListener, err = CreateNatListener(c.InAddr, c.InPort, c.InTLS.Cert, c.InTLS.Key)
	if err != nil {
		return err
	}
	ns.OutListener, err = CreateNatListener(c.OutAddr, c.OutPort, c.OutTLS.Cert, c.OutTLS.Key)
	if err != nil {
		return err
	}
	return nil
}

// HeartBeat

//func (ns *NatServer) HeartBeat() {
//	for {
//		if ns.conConn != nil {
//			_, err := ns.conConn.Write([]byte{0x01, 0xFF})
//			if err != nil {
//				if ns.conConn != nil {
//					ns.conConn.Close()
//				}
//				ns.conConn = nil
//			}
//		}
//		time.Sleep(time.Second * ns.c.Time)
//	}
//}

// OutBound Handle
func (ns *NatServer) OutLink() {
	for {
		oc, err := ns.OutListener.Accept()
		if err != nil {
			log.Println("NAT Server:", err)
			continue
		}
		go ns.OutHandle(oc)
	}
}

func (ns *NatServer) OutHandle(conn net.Conn) {

	var b [1024]byte

	n, err := conn.Read(b[:])
	if err != nil {
		log.Println("Nat Server:", err)
		_, _ = conn.Write([]byte{0x01, 0x02})
		return
	}

	log.Println("Nat Server: [debug]", b[:n])

	// 版本检查
	if n < 18 || b[0] != versionC {
		log.Println("Nat Server:", "error format")
		_, _ = conn.Write([]byte{0x01, 0x02})
		return
	}

	// 授权检查
	auth := b[2:18]

	token := ns.c.Auth
	if token == "" {
		token = defaultPassC
	}

	if !NatAuthVar(token, auth) {
		_, _ = conn.Write([]byte{0x01, 0x01})
		log.Println("Nat Server:", "error token")
		return
	}

	_, err = conn.Write([]byte{0x01, 0x00})
	if err != nil {
		log.Println("Nat Server:", err)
		return
	}

	if b[1] == 0x00 {
		ns.conConn = conn
	} else if b[1] == 0x01 {
		// Todo Connect with InBound
		go ns.Connect(conn)
	}
}

func (ns *NatServer) InLink() {
	for {
		if ns.conConn != nil {
			in, err := ns.InListener.Accept()
			if err != nil {
				log.Println("NAT Server:", err)
				continue
			}
			go ns.InHandle(in)
		}
	}
}

func (ns *NatServer) InHandle(in net.Conn) {
	ns.inConn = in

	_, err := ns.conConn.Write([]byte{0x01, 0x03})

	if err != nil {
		log.Println("NAT Server:", err)
		_ = in.Close()
		ns.inConn = nil
		return
	}

	var b [1024]byte

	for {
		n, err := ns.conConn.Read(b[:])
		if err != nil {
			log.Println("NAT Server:", err)
			_ = in.Close()
			ns.inConn = nil
			return
		}
		log.Println("NAT Server: [Debug]", b[:n])
		if b[0] != versionC || b[1] != 0xFF {
			ns.inConn = nil
			break
		}
	}

	if b[0] != versionC || b[1] != 0x02 {
		log.Println("NAT Server:", "Can not create normal tunnel")
		_ = in.Close()
		return
	}
}

func (ns *NatServer) Connect(out net.Conn) {
	log.Println("Nat Server:", ns.inConn)
	if ns.inConn == nil {
		_ = out.Close()
		return
	}

	in := ns.inConn
	ns.inConn = nil

	if in == nil {
		_ = out.Close()
		return
	}

	defer func() {
		if in != nil {
			_ = in.Close()
		}
		if out != nil {
			_ = out.Close()
		}
	}()

	go io.Copy(in, out)
	io.Copy(out, in)
}

//
//
//func (ns *NatServer) Live() {
//
//	for {
//		if ns.conConn == nil {
//			ns.GetConConn()
//			log.Println(ns.conConn)
//		} else {
//			log.Println(ns.conConn)
//			conn, err := ns.InListener.Accept()
//			if err != nil {
//				log.Println("Nat Server", err)
//				continue
//			}
//			if ns.conConn != nil {
//				go ns.NetServerHandle(conn)
//			} else {
//				conn.Close()
//			}
//		}
//	}
//}
//
//
//func (ns *NatServer) GetConConn() {
//
//	log.Println("NAT Server: [debug] conConn:", ns.conConn)
//	for ns.conConn == nil {
//		log.Println("NAT Server: [debug] conConn:", ns.conConn)
//		mnc, err := ns.OutListener.Accept()
//		if err != nil {
//			log.Println("NAT Server:", err)
//			continue
//		}
//		log.Println(1)
//		ret, isc := ns.OutHandle(mnc)
//		if ret == false {
//			log.Println("NAT Server: Auth error")
//			continue
//		}
//		if isc {
//			ns.conConn = mnc
//		} else {
//			err := mnc.Close()
//			if err != nil {
//				continue
//			}
//		}
//	}
//	log.Println("NAT Server: Get control connection", ns.conConn)
//}
//
//func (ns *NatServer) OutHandle(conn net.Conn) (bool, bool) {
//	isControl := false
//	var b [1024]byte
//
//	n, err := conn.Read(b[:])
//	log.Println("NAT Server: [debug] Out Handle:", b[:n])
//	if err != nil {
//		log.Println("Nat Server", err)
//		conn.Write([]byte{0x01, 0x02})
//		return false, false
//	}
//	log.Println("Nat Server: [debug]", b[:n])
//
//	if n < 18 || b[0] != 0x01 {
//		log.Println("Nat Server:", "error format")
//		conn.Write([]byte{0x01, 0x02})
//		return false, false
//	}
//
//	if b[1] == 0x00 {
//		isControl = true
//	}
//
//	auth := b[2:18]
//
//	token := ns.c.Auth
//	if token == "" {
//		token = "erproxy"
//	}
//
//	if !NatAuthVar(token, auth) {
//		conn.Write([]byte{0x01, 0x01})
//		log.Println("Nat Server:", "error token")
//		return false, false
//	}
//
//	conn.Write([]byte{0x01, 0x00})
//	return true, isControl
//}
//
//func (ns *NatServer) NetServerHandle(conn net.Conn) {
//
//	defer conn.Close()
//
//	ret, nc := ns.Core()
//	if ret == false {
//		conn.Close()
//		return
//	}
//
//	defer nc.Close()
//
//	go io.Copy(nc, conn)
//	io.Copy(conn, nc)
//}
//
//func (ns *NatServer) Core() (bool, net.Conn) {
//	if ns.conConn == nil {
//		return false, nil
//	}
//	ns.conConn.Write([]byte{0x01, 0x03})
//	nc, err := ns.OutListener.Accept()
//	if err != nil {
//		log.Println("123")
//		return false, nil
//	}
//	ret, isc := ns.OutHandle(nc)
//
//	if ret == false || isc == true {
//		return false, nil
//	}
//
//	return true, nc
//}

func CreateNatListener(host, port, certfile, keyfile string) (net.Listener, error) {
	var l net.Listener
	var err error

	if certfile != "" && keyfile != "" {
		cert, err := tls.LoadX509KeyPair(certfile, keyfile)
		if err != nil {
			return nil, err
		}

		config := &tls.Config{Certificates: []tls.Certificate{cert}}

		l, err = tls.Listen("tcp", net.JoinHostPort(host, port), config)
		if err != nil {
			return nil, err
		}
	} else {
		l, err = net.Listen("tcp", net.JoinHostPort(host, port))
		if err != nil {
			return nil, err
		}
	}
	log.Println("NAT Server: Start Nat Server at ", net.JoinHostPort(host, port))
	return l, nil
}
