package main

import (
	"os"
	"os/exec"
	"fmt"
	"log"
	"flag"
	"erproxy/conf"
	"erproxy/core"
)

var (
	logfile string
	conffile string
	back bool
)

func setFlag() {
	// flag.StringVar(&logfile, "l", "erproxy.log", "set logging file")
	flag.StringVar(&conffile, "c", "config.yml", "set configuration file")
	flag.BoolVar(&back, "d",false,"if erproxy needs to run in the background")
	if ! flag.Parsed() {
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
	log.SetFlags(log.LstdFlags)
}

func main(){
	setFlag()
	conf.GetConfig(conffile);

	logfile = conf.CC.Log
	if logfile == ""  {
		logfile = "stdin"
	}

	if back == true {
		st := " -c " +  conffile
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

	l := core.InitServer()
	
	for {
		client, err := l.Accept()
        if err != nil {
            log.Panic(err)
		}
		if conf.CC.InBound.Type == "http" {
			go core.HTTPServerHandle(client)
		}
        go core.Sock5ServerHandle(client)
    }
}