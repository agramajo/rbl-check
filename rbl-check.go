package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	timeout int
	verbose bool
)

// reverse string
func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// worker lookup
func worker(fqdns chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	var timeout = time.Duration(timeout*1000) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	for fqdn := range fqdns {
		var r net.Resolver
		if verbose {
			fmt.Println("QUERY", fqdn)
		}
		ips, _ := r.LookupIPAddr(ctx, fqdn)
		if len(ips) > 0 {
			for _, ip := range ips {
				fmt.Println("FOUND", fqdn, ip.String())
			}
		}
	}
}

func main() {
	var (
		flRbl     = flag.String("list", "", "The RBL list.")
		flIP      = flag.String("ip", "", "The IP to check.")
		flTimeout = flag.Int("t", 5, "The DNS request timeout.")
		flWorkers = flag.Int("c", 16, "The amount of workers to use.")
		flVerbose = flag.Bool("v", false, "Verbose output.")
	)
	flag.Parse()

	if *flRbl == "" || *flIP == "" {
		fmt.Println("-list and -ip are required")
		os.Exit(1)
	}

	verbose = *flVerbose
	timeout = *flTimeout

	var wg sync.WaitGroup
	fqdns := make(chan string, *flWorkers)

	fh, err := os.Open(*flRbl)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)
	for i := 0; i < *flWorkers; i++ {
		wg.Add(1)
		go worker(fqdns, &wg)
	}

	// reverse ip
	ip := strings.Split(*flIP, ".")
	rev := strings.Join(reverse(ip), ".")

	for scanner.Scan() {
		txt := scanner.Text()
		// skip comments
		if strings.HasPrefix(txt, "//") || strings.HasPrefix(txt, "#") {
			continue
		}
		fqdns <- fmt.Sprintf("%s.%s", rev, txt)
	}

	close(fqdns)
	wg.Wait()
}
