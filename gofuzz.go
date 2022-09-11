package gofuzz

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	util "github.com/Th3Beetle/thutils"
)

const (
	mark1 = "G0FUZZ"
	mark2 = "G1FUZZ"
)

func Fuzz(target string, request string, wordlistPath [2]string, resps chan string, errs chan error) {

	var wg sync.WaitGroup
	var scanners [2](*bufio.Scanner)

	for i, path := range wordlistPath {
		if path != "" {
			wordlist, err := os.Open(path)
			if err != nil {
				errs <- fmt.Errorf("failed to open file: %s", err)
			}
			scanners[i] = bufio.NewScanner(wordlist)
		}
	}

	raddr, err := net.ResolveTCPAddr("tcp", target)
	if err != nil {
		errs <- fmt.Errorf("failed to resolve remote addr: %s", err)
	}

	for scanners[0].Scan() {
		var payload [2]string
		payload[0] = scanners[0].Text()
		if scanners[1] != nil {
			for scanners[1].Scan() {
				payload[1] = scanners[1].Text()
				wg.Add(1)
				go sendRequest(request, payload, raddr, resps, &wg, errs)
			}
		} else {
			wg.Add(1)
			go sendRequest(request, payload, raddr, resps, &wg, errs)
		}
	}
	wg.Wait()
	resps <- "fin"

}

func sendRequest(request string, payload [2]string, raddr *net.TCPAddr, response chan string, wg *sync.WaitGroup, errs chan error) {
	defer wg.Done()
	newRequest := insertPayload(request, payload[0], mark1)
	if payload[1] != "" {
		newRequest = insertPayload(request, payload[1], mark2)
	}

	rconn, err := net.DialTCP("tcp", nil, raddr)

	if err != nil {
		errs <- fmt.Errorf("failed to establish connection: %s", err)
	}
	rconn.Write([]byte(newRequest))
	resp, err := util.ReadAll(rconn)

	if err != nil {
		errs <- fmt.Errorf("failed to read response: %s", err)
	}

	response <- string(resp)

}

func insertPayload(request string, payload string, mark string) string {
	parts := strings.Split(request, mark)
	requestWithPayload := strings.Join(parts, payload)
	return requestWithPayload
}
