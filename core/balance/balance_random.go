package balance

import (
	"erproxy/conf"
	"math/rand"
	"time"
)

type RandomBalance struct {
	outs []string
}

func (bl *RandomBalance) Init(balance conf.Balance) {
	bl.outs = OutBoundNameCheck(balance)
}

func (bl *RandomBalance) GetOutbound() string {
	rand.Seed(time.Now().Unix())
	return bl.outs[rand.Intn(len(bl.outs))]
}

func (bl *RandomBalance) Check() {
	return
}
