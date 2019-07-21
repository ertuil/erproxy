package core

import (
	"erproxy/conf"
	"encoding/base64"
)

func isTLS() bool {
	if conf.CC.InBound.TLS.Cert != "" && conf.CC.InBound.TLS.Key != "" {
		return true
	}
	return false
}

func isAuth() bool {
	if len(conf.CC.InBound.Auth) > 0 {
		return true
	}
	return false
}

func getInAddr() (string,string) {
	ad := conf.CC.InBound.Addr
	pt := conf.CC.InBound.Port
	if ad == "" {
		ad = "0.0.0.0"
	}
	if pt == "" {
		pt = "1080"
	}
	return ad,pt
}

func isOutAuth() bool {
	if len(conf.CC.OutBound.Auth)> 0 {
		return true
	}
	return false
}

func getOutAuth() (string,string) {
	for k,v := range(conf.CC.OutBound.Auth) {
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



func authenticate(username string, password string) bool {
	for k,v := range(conf.CC.InBound.Auth) {
		if username == k && password ==  v {
			return true
		}
	}
	return false
}

// HTTPAuth .
func HTTPAuth(token string) bool {
	for k,v := range(conf.CC.InBound.Auth) {
		if token == base64.StdEncoding.EncodeToString([]byte(k+":"+v)) {
			return true
		}
	}
	return false
}