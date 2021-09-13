/* Simple CPU benchmark test

Go servers read messages from their queues, update trace information, and send
them back to the originators.

Go clients create a single message and pass it through all servers
sequentially.  Clients exit when the message passes through all
servers.  The test ends when all clients complete.

 Message format

	CMD:TRACE

	CMD     - command
	            exit    - exit
	            queue   - update trace and continue test

	TRACE   - list of visited clients and servers, the last item is the
	          originator of the message

 For example:

	server 1:	queue:client(1)
	client 1:	queue:client(1)-server(1)
	server 2:	queue:client(1)-server(1)-client(1)
	client 1:	queue:client(1)-server(1)-client(1)-server(2)
	...
*/

package main

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func client(cid int, clientqs, serverqs []chan string, client_exit chan<- int) {

	num_servers := len(serverqs)

	count := num_servers
	dstid := cid % num_servers
	msgout := fmt.Sprintf("queue:client(%d)", cid)
	serverqs[dstid] <- msgout
	for msg := range clientqs[cid] {
		if cli.debug {
			log.Printf("client(%d):  %s", cid, msg)
		}
		if count--; count < 1 {
			break
		}
		dstid = (dstid + 1) % num_servers
		msgout = fmt.Sprintf("%s-client(%d)", msg, cid)
		serverqs[dstid] <- msgout
		runtime.Gosched()
	}
	if cli.debug {
		log.Printf("client(%d):  exit", cid)
	}
	client_exit <- cid
}

func extract_srcid(trace string) int {

	src := trace[strings.LastIndex(trace, "(")+1 : len(trace)-1]
	srcid, _ := strconv.Atoi(src)
	return srcid
}

func server(sid int, clientqs, serverqs []chan string, server_exit chan<- int) {

	for msg := range serverqs[sid] {
		if cli.debug {
			log.Printf("server(%d):  %s", sid, msg)
		}
		toks := strings.Split(msg, ":")
		if toks[0] == "exit" {
			break
		}
		dstid := extract_srcid(toks[1])
		msgout := fmt.Sprintf("%s:%s-server(%d)", toks[0], toks[1], sid)
		clientqs[dstid] <- msgout
		runtime.Gosched()
	}
	if cli.debug {
		log.Printf("server(%d):  exit\n", sid)
	}
	server_exit <- sid
}

func run_qmark(num_clients, num_servers, num_runs int) []time.Duration {

	runs := make([]time.Duration, num_runs) // List of bench mark results

	for ix := range runs {

		clientqs := make([]chan string, num_clients)
		serverqs := make([]chan string, num_servers)
		client_exit := make(chan int)
		server_exit := make(chan int)

		for ii := 0; ii < num_clients; ii++ {
			clientqs[ii] = make(chan string)
		}
		for ii := 0; ii < num_servers; ii++ {
			serverqs[ii] = make(chan string, num_clients)
		}

		start_time := time.Now()

		// start the test

		for ii := 0; ii < num_servers; ii++ {
			go server(ii, clientqs, serverqs, server_exit)
		}
		for ii := 0; ii < num_clients; ii++ {
			go client(ii, clientqs, serverqs, client_exit)
		}

		// wait for clients to complete

		for ii := 0; ii < num_clients; ii++ {
			cid := <-client_exit
			if cli.debug {
				log.Printf("exit client(%d)", cid)
			}
		}
		for ii := 0; ii < num_servers; ii++ {
			serverqs[ii] <- "exit:adm"
		}

		// wait for servers to complete

		for ii := 0; ii < num_servers; ii++ {
			<-server_exit
		}

		// save run result

		result := time.Since(start_time)

		runs[ix] = result
	}

	return runs
}
