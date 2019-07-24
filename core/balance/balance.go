package balance

import (
	"erproxy/conf"
)

var (
	bls map[string]Balance
)

type Balance interface {
	GetOutbound() string
	Init(balance conf.Balance)
	Check()
}

func OutBoundNameCheck(out conf.Balance) []string {
	ret := make([]string, 0)
	for n, _ := range out.OutBound {
		for x, _ := range conf.CC.OutBound {
			if x == n {
				ret = append(ret, x)
			}
		}
	}
	return ret
}

func InitBalance() {
	bls = make(map[string]Balance)
	var tmp Balance

	for n, b := range conf.CC.Balance {
		switch b.Type {
		case "random":
			tmp = new(RandomBalance)
		case "weight":
			tmp = new(WeightBalance)
		case "rr":
			tmp = new(RRBalance)
		case "ping":
			tmp = new(PingBalance)
		case "alive":
			tmp = new(WeightPingBalance)
		default:
			tmp = new(RandomBalance)
		}

		tmp.Init(conf.CC.Balance[n])
		bls[n] = tmp
	}
}

func CheckBalance() {
	for _, bl := range bls {
		go bl.Check()
	}
}

func GetBalance(name string) string {
	b := bls[name]
	if b == nil {
		return ""
	}
	return b.GetOutbound()
}
