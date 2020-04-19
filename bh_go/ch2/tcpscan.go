package main

import (
	"flag"
	"fmt"
	"net"
	"sort"
)

func main() {
	target := flag.String("target", "127.0.0.1", "Target you want to port scan")
	portRange := flag.Int("portRange", 1024, "Range for Ports to Scan on Target")
	flag.Parse()
	openPorts := TCPScan(*target, *portRange)
	fmt.Printf("%s : %d\n", *target, openPorts)
}

func worker(target string, portRange, results chan int) {
	for port := range portRange {
		address := fmt.Sprintf("%s:%d", target, port)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- port
	}
}

// TCPScan - Scan for open ports on target
func TCPScan(target string, portRange int) []int {
	var openPorts []int
	results := make(chan int)
	ports := make(chan int, 100)
	defer close(results)
	defer close(ports)
	for i := 1; i < cap(ports); i++ {
		go worker(target, ports, results)
	}

	go func(portRange int) {
		for i := 1; i <= portRange; i++ {
			ports <- i
		}
	}(portRange)

	for i := 0; i < portRange; i++ {
		port := <-results
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}
	sort.Ints(openPorts)
	fmt.Printf("Found %d open ports!\n", len(openPorts))
	return openPorts
}
