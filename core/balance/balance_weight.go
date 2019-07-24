package balance

import (
	"erproxy/conf"
	"math/rand"
	"time"
)

type WeightBalance struct {
	outs []string
}

func (bl *WeightBalance) Init(b conf.Balance) {
	valueName := OutBoundNameCheck(b)
	bl.outs = make([]string, 0)
	for _, n := range valueName {
		tmp := make([]string, b.OutBound[n])
		for k, _ := range tmp {
			tmp[k] = n
		}
		bl.outs = append(bl.outs, tmp...)
	}
}

func (bl *WeightBalance) GetOutbound() string {
	rand.Seed(time.Now().Unix())
	return bl.outs[rand.Intn(len(bl.outs))]
}

func (bl *WeightBalance) Check() {
	return
}
