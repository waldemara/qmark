package main

import (
	"flag"
	"log"
)

var cli struct {
	build   bool
	debug   bool
	quiet   bool
	clients int
	servers int
	runs    int
}

func parse_cli(prog string) {

	log.SetFlags(0)

	flag.BoolVar(&cli.build, "b", false, "print build id and exit")
	flag.BoolVar(&cli.quiet, "q", false, "quiet mode, only print qmark value")
	//flag.BoolVar(&cli.debug, "debug", false, "print debug information")
	flag.IntVar(&cli.clients, "clients", 1705, "number of clients")
	flag.IntVar(&cli.servers, "servers", 347, "number of servers")
	flag.IntVar(&cli.runs, "runs", 7, "number of runs")
	flag.Usage = func() {

		log.Println("Simple cpu benchmark based on message passing between a large")
		log.Println("number of go routines. Default values are tuned to produce")
		log.Println("qmark value of about 500 on a six core 2.2 GHz Xeon.")
		log.Println("")
		log.Println(" ", prog, "[OPTIONS]")
		log.Println("")
		log.Println("options:")
		log.Println("")
		flag.PrintDefaults()
		log.Println("")
	}
	flag.Parse()

	if cli.clients < 1 {
		cli.clients = 1
	}
	if cli.servers < 1 {
		cli.servers = 1
	}
	if cli.runs < 1 {
		cli.runs = 1
	}
}
