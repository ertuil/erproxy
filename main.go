package main

import (
	"erproxy/conf"
	"erproxy/core"
	"erproxy/core/balance"
	"erproxy/core/nat"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

var (
	logfile  string
	conffile string
	back     bool
	sw       sync.WaitGroup
)

func setFlag() {
	// flag.StringVar(&logfile, "l", "erproxy.log", "set logging file")
	flag.StringVar(&conffile, "c", "config.yml", "set configuration file")
	flag.BoolVar(&back, "d", false, "if erproxy needs to run in the background")
	if !flag.Parsed() {
		flag.Parse()
	}
}
func setLog(logfile string) {

	if logfile != "stdin" {
		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalln("Can not open log file.")
		}
		log.SetOutput(f)
	}
	log.SetPrefix("[erproxy]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	setFlag()
	conf.GetConfig(conffile)

	logfile = conf.CC.Log
	if logfile == "" {
		logfile = "stdin"
	}

	if back == true {
		st := " -c " + conffile
		cmd := exec.Command(os.Args[0], st)
		err := cmd.Start()
		fmt.Println(os.Args[0] + st)
		if err != nil {
			fmt.Println(err)
		}
		back = false
		os.Exit(0)
	}

	setLog(logfile)

	log.Printf("Erproxy start, config file: %v", conffile)

	// Init Balances

	balance.InitBalance()
	balance.CheckBalance()
	ti := time.NewTicker(300 * time.Second)

	go func() {
		for range ti.C {
			balance.CheckBalance()
		}
	}()

	// NAT
	nat.InitNat(conf.CC.Nat)
	for _, s := range nat.Servers {
		sw.Add(1)
		go s.OutLink()
		go s.InLink()
	}

	for _, c := range nat.Clients {
		sw.Add(1)
		go c.InLink()
		go c.HeartBeat()
	}
	// Inbound and OutBound

	for n, c := range conf.CC.InBound {
		var ib core.Inbound
		if c.Type == "socks" {
			ib = new(core.Socks5Server)
		} else if c.Type == "http" {
			ib = new(core.HTTPServer)
		} else {
			ib = new(core.SUTPServer)
		}
		sw.Add(1)
		go core.InBoundServerRun(n, ib, c)
	}
	sw.Wait()
	for {
		time.Sleep(1 * time.Second)
	}
}
