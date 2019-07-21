package core

import (
	"erproxy/conf"
	"encoding/base64"
)

func isTLS(c conf.InBound) bool {
	if c.TLS.Cert != "" && c.TLS.Key != "" {
		return true
	}
	return false
}

func isAuth(c conf.InBound) bool {
	if len(c.Auth) > 0 {
		return true
	}
	return false
}

func getInAddr(c conf.InBound) (string,string) {
	ad := c.Addr
	pt := c.Port
	if ad == "" {
		ad = "0.0.0.0"
	}
	if pt == "" {
		pt = "1080"
	}
	return ad,pt
}

func isOutAuth(c conf.OutBound) bool {
	if len(c.Auth)> 0 {
		return true
	}
	return false
}

func getOutAuth(c conf.OutBound) (string,string) {
	for k,v := range(c.Auth) {
		return k,v
	}
	return "",""
}

func selectMethod(a []byte,b byte) bool {
	for _,c := range(a) {
		if c == b {
			return true
		}
	}
	return false
}



func authenticate(username string, password string,c conf.InBound) bool {
	for k,v := range(c.Auth) {
		if username == k && password ==  v {
			return true
		}
	}
	return false
}

// HTTPAuth .
func HTTPAuth(token string,c conf.InBound) bool {
	for k,v := range(c.Auth) {
		if token == base64.URLEncoding.EncodeToString([]byte(k+":"+v)) {
			return true
		}
	}
	return false
}