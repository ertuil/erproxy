package core

import (
	"net"
	"log"
	"strings"
	"erproxy/conf"
)

// ResultStatus ,
type ResultStatus int

const (
	// Direct .
	Direct ResultStatus = iota 
	// Block .
	Block	
	// Proxy .				
	Proxy
)

func getOutBound(from, host, port string,atype byte) Outbound {
	var ob Outbound
	name,c := route(from,host, port,atype)
	switch(c.Type) {
	case "socks": ob = new(sockbound)
	case "http": ob =  new(httpbound)
	case "free": ob = new(freebound)
	case "block": ob = new(blockbound)
	default :  ob = new(blockbound);name="block";c=conf.OutBound{Type: "block"}
	}
	log.Println("Route: from",from,"to",net.JoinHostPort(host, port),"via",name)
	ob.init(name,c)
	return ob
}

func route(from, host, port string , atype byte) (string,conf.OutBound) {
	for rule,policy := range(conf.CC.Routes.Route) {
		if routeMatch(from, host, port,atype,rule) {
			return getPolicy(policy)
		}
	}
	return getPolicy(conf.CC.Routes.Default)
}

func routeMatch(from, host, port string , atype byte, rule string) bool {
	testips := make([]net.IP,0)

	if atype == 0x01 || atype == 0x04 {
		t := net.ParseIP(host)
		if t != nil {
			testips = append(testips,t)
		}
	} else {
		ts,err := net.LookupHost(host)
		if err == nil {
			for _,v := range(ts) {
				t := net.ParseIP(v)
				if t != nil {
					testips = append(testips,t)
				}
			}
		}
	}


	// Interface route
	ruleFrom := strings.Split(rule,"@")
	if len(ruleFrom) >= 2 {
		rule = ruleFrom[0]
	}
	if len(ruleFrom) >= 2 && ruleFrom[1] != from {
		return false
	}

	// IP route
	ruleIP := net.ParseIP(rule)
	if ruleIP != nil {
		for _,v := range(testips) {
			if v.Equal(ruleIP) {
				return true
			}
		}
		return false
	}

	// CIDR route
	ruleIP,ruleCIDR,err := net.ParseCIDR(rule)
	if err == nil {
		for _,v := range(testips) {
			if ruleCIDR.Contains(v) {
				return true
			}
			return false
		}
	}

	// port route
	if strings.Contains(rule,"port:") {
		if rule[5:] == port {
			return true
		}
	}

	// Daemon route
	return strings.Contains(host,rule)
}

func getPolicy(v string) (string,conf.OutBound) {

	for n,c := range(conf.CC.OutBound) {
		if n == v {
			return n,c
		}
	}
	return "block",conf.OutBound{Type: "block"}
}
