package core

import (
	"erproxy/conf"
)

func getOutBound() outbound {
	var out outbound
	if conf.CC.OutBound.Type == "free" {
		out = new(freebound)
		return out
	} else if conf.CC.OutBound.Type == "sock" {
		out = new(sockbound)
		return out
	}
	return nil
}