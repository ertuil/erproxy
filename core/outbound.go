package core

import (
	"io"
	"net"
	"log"
)

type outbound interface {
	getserver() net.Conn
	start(string, string, byte) bool
	loop(net.Conn)
	close()
}

type freebound struct {
	server net.Conn
}

func (fb *freebound) getserver() net.Conn{
	return fb.server
}

func (fb *freebound) start(host,  port string, atype byte) bool {
	server, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Println("Free Client:",err)
		return false
	}
	fb.server = server
	log.Println("Free Client: Try to connect to", host + ":" + port)
	return true
}

func (fb *freebound) loop(client net.Conn){
	go io.Copy(fb.server, client)
	io.Copy(client, fb.server)
}

func (fb *freebound) close(){
	if fb.server != nil {
		fb.server.Close()
	}
}

type blockbound struct {}

func (bb *blockbound) getserver() net.Conn{
	return nil
}

func (bb *blockbound) start(host,  port string, atype byte) bool {
	log.Println("Block Client: Block connection to", host + ":" + port)
	return false
}

func (bb *blockbound) loop(client net.Conn) {
	return 
}

func (bb *blockbound) close()  {
	return
}