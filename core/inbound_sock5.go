package core

import (
	"erproxy/conf"
	"erproxy/header"
	"log"
	"net"
	"strconv"
)

//Socks5Server .
type Socks5Server struct {
	c         conf.InBound
	name      string
	udpServer Socks5UDPServer
}

// Init for socks5server
func (ss *Socks5Server) Init(name string, c conf.InBound) net.Listener {
	ss.c = c
	ss.name = name
	if ss.c.UDPPort != "" {
		us := new(Socks5UDPServer)
		us.Init(name, c)
	}
	return InitServer(ss.c)
}

// Handle for socks5server
func (ss *Socks5Server) Handle(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	// log.Printf("[erproxy] Connected from %v ...", client.RemoteAddr())
	ret := Sock5HandShake(client, ss.c)
	if ret != true {
		log.Printf("handshake failed")
		return
	}

	if isAuth(ss.c) {
		ret = Sock5Auth(client, ss.c)
		if ret == false {
			log.Printf("authenticate failed")
			return
		}
	}

	Socks5Request(ss.name, client, ss)

}

// Sock5HandShake result, ip, port
func Sock5HandShake(client net.Conn, c conf.InBound) bool {
	var b [1024]byte

	// Read hand shake message
	_, err := client.Read(b[:])
	if err != nil {
		log.Println("Socks Server:", err)
		return false
	}

	if b[0] != 0x05 {
		log.Println("Socks Server: Protocal error or version error")
		return false
	}

	nm := b[1]
	ms := b[2 : 2+nm]

	if isAuth(c) && selectMethod(ms, 0x02) {
		client.Write([]byte{0x05, 0x02})
		return true
	} else if !isAuth(c) && selectMethod(ms, 0x00) {
		client.Write([]byte{0x05, 0x00})
		return true
	} else {
		client.Write([]byte{0x05, 0xFF})
		return false
	}
}

// Sock5Auth auth function
func Sock5Auth(client net.Conn, c conf.InBound) bool {
	var b [1024]byte

	_, err := client.Read(b[:])
	if err != nil {
		log.Println("Socks Server:", err)
		client.Write([]byte{0x01, 0x01})
		return false
	}

	v1 := b[0]
	if v1 != 0x01 {
		log.Println("Socks Server:", "Autnenticate version error")
		client.Write([]byte{0x01, 0x01})
		return false
	}
	un := b[1]
	u := b[2 : 2+un]
	pn := b[2+un]
	p := b[2+un+1 : 2+un+1+pn]

	if authenticate(string(u), string(p), c) {
		client.Write([]byte{0x01, 0x00})
		return true
	}

	client.Write([]byte{0x01, 0x01})
	return false
}

// Socks5Request is main request: ret, cmd, host, port
func Socks5Request(name string, client net.Conn, ss *Socks5Server) {
	var b [1024]byte
	f := []byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	n, err := client.Read(b[:])
	if err != nil {
		log.Println("Socks Server:", err)
		client.Write(f)
		return
	}

	v := b[0]
	if v != 0x05 {
		log.Println("Socks Server:", "Socks version error")
		client.Write(f)
		return
	}
	cmd := b[1]
	atyp := b[3]
	var host string
	var port string

	port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

	switch atyp {
	case 0x01:
		host = net.IPv4(b[4], b[5], b[6], b[7]).String()
	case 0x03:
		host = string(b[5 : n-2])
		atyp = header.HostCheck(host)
	case 0x04:
		host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13],
			b[14], b[15], b[16], b[17], b[18], b[19]}.String()
	default:
		log.Println("Socks Server:", "Socks version error")
		client.Write(f)
		return
	}

	ad := header.AddrInfo{}
	ad.SetInfo(name, host, port, atyp, 0x01)
	switch cmd {
	case 0x01:
		Socks5Connect(client, ad)
	case 0x03:
		Socks5UDP(client, ad, ss)
	default:
		log.Println("Socks Server:", "cmd not understand")
		client.Write(f)
		return
	}

}

// Socks5Connect connect
func Socks5Connect(client net.Conn, ad header.AddrInfo) {

	s := []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	f := []byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	out := getOutBound(ad)
	if out == nil {
		client.Write(f)
		log.Println("Socks Server: Cannot connect to outbound")
		return
	}
	ret := out.start(ad)
	if ret != true {
		client.Write(f)
		log.Println("Socks Server: Cannot connect to outbound")
		return
	}

	client.Write(s)

	defer out.close()
	out.loop(client)
}

func Socks5UDP(client net.Conn, ad header.AddrInfo, ss *Socks5Server) {

	ss.udpServer.addAllow(getUDPAddrIP(client.RemoteAddr().String()), ad)

	b := []byte{0x05, 0x00, 0x00}

	var atype byte
	port := ss.c.UDPPort
	host := ss.c.Addr

	t := net.ParseIP(host)
	ip, err := t.MarshalText()
	if err != nil {
		log.Println("Socks Server: Cannot marshal ip")
		return
	}

	if len(ip) == 4 {
		atype = 0x01
	} else if len(ip) == 16 {
		atype = 0x03
	} else {
		log.Println("Socks Server: Cannot marshal ip")
		return
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		log.Println("Socks Server: Can not marshal port")
		return
	}
	var pp [2]byte
	pp[0] = byte(p / 256)
	pp[1] = byte(p % 256)

	b = append(b, atype)
	b = append(b, ip...)
	b = append(b, port[0], port[1])

	client.Write(b)
}
