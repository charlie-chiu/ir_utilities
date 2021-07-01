package main

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

func main() {
	targetURL := "http://go.dev3x.club"
	portStart := 9000
	portEnd := 10100

	portsAvailable := make([]int, 0, portEnd-portStart+1)
	portsUnavailable := make([]int, 0, portEnd-portStart+1)

	wg := sync.WaitGroup{}
	wg.Add(portEnd - portStart + 1)
	for port := portStart; port <= portEnd; port++ {
		go func(p int) {
			defer wg.Done()
			urlWithPort := targetURL + ":" + strconv.Itoa(p)
			resp, err := http.Get(urlWithPort)
			if err != nil {
				log.Printf("http get error, %v\n", err)
				portsUnavailable = append(portsUnavailable, p)
			} else {
				log.Printf("url with port: %s, status code: %d\n", urlWithPort, resp.StatusCode)
				portsAvailable = append(portsAvailable, p)
			}
		}(port)
	}
	wg.Wait()

	sort.Ints(portsAvailable)
	sort.Ints(portsUnavailable)
	log.Println("scan completed")
	log.Printf("available ports : %v\n", portsAvailable)
	log.Printf("unavailable ports : %v\n", portsUnavailable)

}
