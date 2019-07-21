package core

import (
	"io"
	"net"
	"log"
	"erproxy/conf"
)

//Outbound is outboubd client 
type Outbound interface {
	getserver() net.Conn
	init(string, conf.OutBound)
	start(string, string, byte) bool
	loop(net.Conn)
	close()
}

type freebound struct {
	name string
	c conf.OutBound
	server net.Conn
}

func (fb *freebound) getserver() net.Conn{
	return fb.server
}

func (fb *freebound) init(name string, c conf.OutBound) {
	fb.name = name
	fb.c =  c
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

type blockbound struct {
	name string
	c conf.OutBound
}

func (bb *blockbound) getserver() net.Conn{
	return nil
}

func (bb *blockbound) init(name string, c conf.OutBound) {
	bb.name = name
	bb.c =  c
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