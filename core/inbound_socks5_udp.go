package core

import (
	"erproxy/conf"
	"erproxy/header"
	"log"
	"net"
	"strconv"
	"strings"
)

type Socks5UDPServer struct {
	name  string
	c     conf.InBound
	allow []header.AddrInfo
}

func (ss *Socks5UDPServer) Init(name string, c conf.InBound) {
	ss.name = name
	ss.c = c
	port := c.UDPPort
	addr := c.Addr
	log.Println("[erproxy] starting Socks UDP Server in", net.JoinHostPort(addr, port))
	go ss.Run(addr, port)
}

func (ss *Socks5UDPServer) Run(addr, port string) {
	for {
		client, err := net.ListenPacket("udp", net.JoinHostPort(addr, port))
		if err != nil {
			log.Println("Socks5 UDP Server:", err)
		} else {
			ss.Handle(client)
		}

	}
}

func (ss *Socks5UDPServer) addAllow(addr string, ad header.AddrInfo) {
	ss.allow = append(ss.allow, ad)
}

func (ss *Socks5UDPServer) delAllow(ad header.AddrInfo) {
	_, host, port, _, _ := ad.GetInfo()
	udpFrom := ad.UDPFrom
	newallow := make([]header.AddrInfo, len(ss.allow))

	for _, a := range ss.allow {
		_, rh, rp, _, _ := a.GetInfo()
		if udpFrom != a.UDPFrom || rh != host || rp != port {
			newallow = append(newallow, a)
		}
	}

	ss.allow = newallow
}

func checkAllow(allow []header.AddrInfo, udpfrom, host, port string) bool {
	for _, a := range allow {
		_, rh, rp, _, _ := a.GetInfo()
		if rh == host && rp == port && udpfrom == a.UDPFrom {
			return true
		}
	}
	return false
}

func (ss *Socks5UDPServer) Handle(client net.PacketConn) {

	defer client.Close()

	var b [65535]byte
	n, udpFroms, err := client.ReadFrom(b[:])
	udpFrom := getUDPAddrIP(udpFroms.String())

	if err != nil {
		log.Println("Socks UDP Server:", err)
	}

	atype := b[3]

	var host string
	var port string
	var l int
	//port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

	switch atype {
	case 0x01:
		if n < 10 {
			log.Println("Socks UDP Server: Too Short")
			return
		}
		host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		port = strconv.Itoa(int(b[8])<<8 | int(b[9]))
		l = 10
	case 0x03:
		l = int(b[4])
		if n < 7+l {
			log.Println("Socks UDP Server: Too Short")
			return
		}
		host = string(b[5 : 5+l])
		port = strconv.Itoa(int(b[5+l])<<8 | int(b[6+l]))
		l = 7 + l
	case 0x04:
		if n < 22 {
			log.Println("Socks UDP Server: Too Short")
			return
		}
		host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13],
			b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		port = strconv.Itoa(int(b[20])<<8 | int(b[21]))
		l = 22
	default:
		log.Println("Socks UDP Server:", "Socks version error")
		return
	}

	log.Println(host, port, atype)
	ret := checkAllow(ss.allow, udpFrom, host, port)
	if ret == false {
		log.Println("Socks UDP Server:", "Dest not allow")
		return
	}

	ad := header.AddrInfo{}
	ad.SetInfo(ss.name, host, port, atype, 0x02)

	ob := getOutBound(ad)
	log.Println("Socks UDP Server:Trys to connect to", ad.Host, ad.Port)
	ob.start(ad)
	defer ob.close()

	if ob == nil {
		log.Println("Socks UDP Server:", "Can not connect to outbound")
		return
	}
	ret = ob.start(ad)
	if ret != true {
		log.Println("Socks UDP Server:", "Can not connect to outbound")
		return
	}

	server := ob.getserver()
	server.Write(b[l:n])
	var resp [65535]byte
	n, err = server.Read(resp[:])
	if err != nil {
		return
	}
	log.Println("Socks UDP Server:resv", resp[:n])
	_, err = client.WriteTo(resp[:n], udpFroms)
	if err != nil {
		log.Println("Socks UDP Server:", "Can not send to client")
	}
	log.Println("Socks UDP Server:send to client")
}

func getUDPAddrIP(addr string) string {
	tmp := strings.Split(addr, ":")
	tmp = tmp[:len(tmp)-1]
	udpFrom := strings.Join(tmp, ":")
	return udpFrom
}
