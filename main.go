package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
	"time"
)

//go:generate ./set_build_id

func main() {

	toks := strings.Split(os.Args[0], "/")
	prog := toks[len(toks)-1]

	parse_cli(prog)

	if cli.build {
		log.Print(build_id)
		return
	}

	if !cli.quiet {
		log.Printf("qmark build: %v  go version: %v  num cpus: %v",
			build_id, runtime.Version(), runtime.NumCPU())

		log.Printf("calculating qmark... ")
	}

	var sum time.Duration

	results := run_qmark(cli.clients, cli.servers, cli.runs)

	for _, res := range results {
		sum += res
	}

	avg := (float64(sum)/float64(len(results)))/float64(time.Second) // [s]
	qmark := int(1000.0 / avg)

	if cli.quiet {

		log.Print(qmark)

	} else {

		sumsqr := 0.0
		for _, res := range results {
			diff := float64(res)/1000000000.0 - avg
			sumsqr += diff * diff
		}
		stdev := math.Sqrt(sumsqr / float64(len(results)))

		rstr := ""
		for _, res := range results {
			rstr += fmt.Sprintf("  %5.3f", float64(res)/1000000000.0)
		}
		log.Printf("results [s]:%v", rstr)
		//log.Printf("")
		log.Printf("average [s]:  %5.3f", avg)
		log.Printf("stdev [s]:    %5.3f", stdev)
		log.Printf("qmark:        %d", qmark)

	}
}
