package nat

import (
	"crypto/md5"
)

func NatAuth(token string) []byte {
	m := md5.New()
	m.Write([]byte(token + "erproxy"))
	return m.Sum(nil)
}

func NatAuthVar(token string, auth []byte) bool {
	if len(auth) != 16 {
		return false
	}
	ntk := NatAuth(token)

	ret := true
	for i := range auth {
		if auth[i] != ntk[i] {
			ret = false
		}
	}
	return ret
}
