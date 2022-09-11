package main

import (
	"fmt"

	"github.com/Th3Beetle/gofuzz"
)

func main() {
	wordlistPath := [2]string{"/home/greey/wordlist.txt", ""}
	request := "GET /G0FUZZ HTTP/1.1\r\nUser-Agent: Mozilla/4.0 (compatible; MSIE5.01; Windows NT)\r\nHost: www.tutorialspoint.com\r\nAccept-Language: en-us\r\nAccept-Encoding: gzip, deflate\r\nConnection: Keep-Alive\r\n\r\n"
	remoteAddr := "127.0.0.1:80"

	resps := make(chan string)
	errs := make(chan error)
	go readErrors(errs)
	go gofuzz.Fuzz(remoteAddr, request, wordlistPath, resps, errs)

	for {
		response := <-resps
		if response == "fin" {
			break
		}
		fmt.Println("New response:")
		fmt.Println(response)
	}
}

func readErrors(errs chan error) {
	for {
		err := <-errs
		fmt.Println(err)
	}
}
