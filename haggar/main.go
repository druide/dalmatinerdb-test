package main

import (
	"flag"

	"bytes"
	"encoding/binary"

	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"
)

// config vars, to be manipulated via command line flags
var (
	carbon        string
	prefix        string
	flushInterval time.Duration
	spawnInterval time.Duration
	metrics       int
	jitter        time.Duration
	agents        int
)

type Agent struct {
	ID            int
	FlushInterval time.Duration
	Conn          net.Conn
}

var events []string = []string{
	"request",
	"impression",
	"creativeView",
	"start",
	"firstQuartile",
	"midpoint",
	"thirdQuartile",
  "complete",
}
var elen int = len(events)

func (a *Agent) Start() {
	for {
		select {
		case <-time.NewTicker(a.FlushInterval).C:
			err := a.flush()
			if err != nil {
				log.Printf("agent %d: %s\n", a.ID, err)
			}
		}
	}
}

func (a *Agent) flush() error {
  conn := a.Conn
	epoch := time.Now().Unix()
	n := metrics

	var (
		domain string
		aid string
		cmp string
		event string
		agent string = fmt.Sprintf("a%3.3d", a.ID)
	)
	//for i := 0; i < n; i++ {
		//domain = fmt.Sprintf("domain%3.3d", i)
		for j := 0; j < n; j++ {
			aid = fmt.Sprintf("%3.3d", j)
			//for k := 0; k < n; k++ {
			//	cmp = fmt.Sprintf("%3.3d", k)
				for h := 0; h < elen; h++ {
					event = events[h]
					err := carbonate(conn, epoch, agent, domain, aid, cmp, event)
					if err != nil {
						return err
					}
				}
			//}
		}
	//}

  buf := new(bytes.Buffer)
  binary.Write(buf, binary.BigEndian, uint8(6))
  conn.Write(buf.Bytes())

	log.Printf("agent %d: flushed %d metrics\n", a.ID, n * /*n * n * */ elen)
	return nil
}

func launchAgent(id, n int, flush time.Duration, addr, prefix string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
        buf := new(bytes.Buffer)
        blen := len(prefix)
        size := 1 +    // prefix
                1 +    // max delta
                1 +    // Bucket size
                blen   // bucket string
        binary.Write(buf, binary.BigEndian, uint32(size))
        binary.Write(buf, binary.BigEndian, uint8(4))
        binary.Write(buf, binary.BigEndian, uint8(2))
        binary.Write(buf, binary.BigEndian, uint8(blen))
        conn.Write(buf.Bytes())
        conn.Write([]byte(prefix))

	a := &Agent{
		ID:            id,
		FlushInterval: time.Duration(flush),
		Conn:          conn}
	a.Start()
	return nil
}

func init() {
	flag.StringVar(&carbon, "carbon", "localhost:2003", "address of carbon host")
	flag.StringVar(&prefix, "prefix", "haggar", "prefix for metrics")
	flag.DurationVar(&flushInterval, "flush-interval", 10*time.Second, "how often to flush metrics")
	flag.DurationVar(&spawnInterval, "spawn-interval", 10*time.Second, "how often to gen new agents")
	flag.IntVar(&metrics, "metrics", 10000, "number of metrics for each agent to hold")
	flag.DurationVar(&jitter, "jitter", 10*time.Second, "max amount of jitter to introduce in between agent launches")
	flag.IntVar(&agents, "agents", 100, "max number of agents to run concurrently")
}

func main() {
	flag.Parse()

	spawnAgents := true
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGUSR1)

	// start timer at 1 milli so we launch an agent right away
	timer := time.NewTimer(1 * time.Millisecond)

	curID := 0

	log.Printf("master: pid %d\n", os.Getpid())

	for {
		select {
		case <-sigChan:
			spawnAgents = !spawnAgents
			log.Printf("master: spawn_agents=%t\n", spawnAgents)
		case <-timer.C:
			if curID < agents {
				if spawnAgents {
					go launchAgent(curID, metrics, flushInterval, carbon, prefix)
					log.Printf("agent %d: launched\n", curID)
					curID++

					timer = time.NewTimer(spawnInterval + (time.Duration(rand.Int63n(jitter.Nanoseconds()))))
				}
			}
		}
	}
}
