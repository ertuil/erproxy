package nat

import (
	"erproxy/conf"
	"log"
)

var (
	Servers map[string]NatServer
	Clients map[string]NatClient
)

func InitNat(c conf.Nat) {
	Servers = make(map[string]NatServer)
	Clients = make(map[string]NatClient)

	if len(c.Server) > 0 {
		for n, c := range c.Server {
			ns := NatServer{}
			err := ns.Init(n, c)
			if err == nil {
				log.Println("123")
				Servers[n] = ns
			}
		}
	}

	if len(c.Client) > 0 {
		for n, c := range c.Client {
			ns := NatClient{}
			ns.Init(n, c)
			Clients[n] = ns
		}
	}
}
