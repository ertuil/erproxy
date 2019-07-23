package core

import (
	"crypto/tls"
	"erproxy/conf"
	"log"
	"net"
)

// Inbound is inbound server
type Inbound interface {
	Init(string, conf.InBound) net.Listener
	Handle(client net.Conn)
}

// InitServer is socks5 server
func InitServer(c conf.InBound) net.Listener {

	// Read from configuration
	var istls bool
	var certfile, keyfile, ip, port string

	if !isTLS(c) {
		istls = false
	} else {
		istls = true
		certfile = c.TLS.Cert
		keyfile = c.TLS.Key
	}

	ip, port = getInAddr(c)

	// TLS
	if istls {
		cert, err := tls.LoadX509KeyPair(certfile, keyfile)
		if err != nil {
			log.Fatalln(err)
		}

		config := &tls.Config{Certificates: []tls.Certificate{cert}}

		l, err := tls.Listen("tcp", ip+":"+port, config)
		log.Println("Starting "+c.Type+" server(tls) in ", l.Addr())
		if err != nil {
			log.Fatalln(err)
		}
		return l
	}

	// TCP with out TLS
	l, err := net.Listen("tcp", ip+":"+port)
	log.Println("[erproxy] starting "+c.Type+" server in ", l.Addr())
	if err != nil {
		log.Fatalln(err)
	}
	return l
}

// InBoundServerRun run inbound servers
func InBoundServerRun(name string, ib Inbound, c conf.InBound) {
	l := ib.Init(name, c)
	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		go ib.Handle(client)
	}
}
