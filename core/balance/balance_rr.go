package balance

import (
	"erproxy/conf"
)

type RRBalance struct {
	outs map[string]int
}

func (bl *RRBalance) Init(b conf.Balance) {
	valueName := OutBoundNameCheck(b)
	bl.outs = make(map[string]int, 0)

	for _, n := range valueName {
		bl.outs[n] = 0
	}
}

func (bl *RRBalance) GetOutbound() string {

	flag := true
	for _, c := range bl.outs {
		if c != 1 {
			flag = false
			break
		}
	}

	if flag {
		for n, _ := range bl.outs {
			bl.outs[n] = 0
		}
	}

	for n, c := range bl.outs {
		if c == 0 {
			bl.outs[n] = 1
			return n
		}
	}

	for n, _ := range bl.outs {
		return n
	}
	return ""
}

func (bl *RRBalance) Check() {
	return
}
