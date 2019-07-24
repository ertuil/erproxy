package balance

import (
	"erproxy/conf"
	"github.com/sparrc/go-ping"
)

type PingBalance struct {
	outs map[string]int
}

func (bl *PingBalance) Init(balance conf.Balance) {
	n := OutBoundNameCheck(balance)
	bl.outs = make(map[string]int)
	for _, n := range n {
		bl.outs[n] = 0
	}
}

func (bl *PingBalance) GetOutbound() string {
	var ln string = ""
	var lt int = 1000000000
	for n, b := range bl.outs {
		if b < lt {
			lt = b
			ln = n
		}
	}
	return ln
}

func (bl *PingBalance) Check() {
	for n := range bl.outs {
		addr := conf.CC.OutBound[n].Addr

		if conf.CC.OutBound[n].Type == "free" {
			bl.outs[n] = 0
			continue
		} else if conf.CC.OutBound[n].Type == "block" {
			bl.outs[n] = 999999999
			continue
		}

		Pinger, err := ping.NewPinger(addr)
		if err != nil {
			bl.outs[n] = 999999999
			continue
		}

		Pinger.Count = 3
		Pinger.Run()
		stats := Pinger.Statistics()
		if stats.PacketLoss > 0.6 {
			bl.outs[n] = 999999999
			continue
		}
		bl.outs[n] = int(stats.AvgRtt)
	}

}
