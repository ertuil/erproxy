package balance

import (
	"erproxy/conf"
	"github.com/sparrc/go-ping"
	"log"
	"math/rand"
	"time"
)

type WeightPingBalance struct {
	b    conf.Balance
	loss map[string]float64
}

func (bl *WeightPingBalance) Init(b conf.Balance) {
	bl.b = b
	n := OutBoundNameCheck(b)
	bl.loss = make(map[string]float64)
	for _, n := range n {
		bl.loss[n] = 0
	}
}

func (bl *WeightPingBalance) GetOutbound() string {
	log.Println(bl.loss)
	vn := make([]string, 0)

	for n, v := range bl.loss {
		if v > 0.6 {
			continue
		}

		tmp := make([]string, bl.b.OutBound[n])
		for k, _ := range tmp {
			tmp[k] = n
		}

		vn = append(vn, tmp...)
	}
	log.Println(bl.loss, vn)
	rand.Seed(time.Now().Unix())
	return vn[rand.Intn(len(vn))]
}

func (bl *WeightPingBalance) Check() {

	for n := range bl.loss {
		addr := conf.CC.OutBound[n].Addr
		if conf.CC.OutBound[n].Type == "free" {
			bl.loss[n] = 0
			continue
		} else if conf.CC.OutBound[n].Type == "block" {
			bl.loss[n] = 1
			continue
		}

		Pinger, err := ping.NewPinger(addr)
		if err != nil {
			bl.loss[n] = 1
			continue
		}

		Pinger.Count = 3
		Pinger.Run()
		stats := Pinger.Statistics()
		bl.loss[n] = stats.PacketLoss
	}

}
