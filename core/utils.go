package core

import (
	"erproxy/conf"
)

func isTLS() bool {
	if conf.CC.InBound.TLS.Cert != "" && conf.CC.InBound.TLS.Key != "" {
		return true
	}
	return false
}

func isAuth() bool {
	if conf.CC.InBound.Auth.User != "" && conf.CC.InBound.Auth.Token != "" {
		return true
	}
	return false
}

func isOutAuth() bool {
	if conf.CC.OutBound.Auth.User != "" && conf.CC.OutBound.Auth.Token != "" {
		return true
	}
	return false
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
	if username == conf.CC.InBound.Auth.User && password == conf.CC.InBound.Auth.Token {
		return true
	}
	return false
}