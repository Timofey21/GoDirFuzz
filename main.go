package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strconv"
	
)

func main() {


	url:= os.Args[len(os.Args) - 1]
	println(url)

	file, err := os.Open("fuzz.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var fuzzUrl string

	for scanner.Scan() {

		fuzzUrl = url + "/" + scanner.Text()
		resp, err := http.Get(fuzzUrl)
		if err != nil {
			log.Fatal(err)


		}

		numStatus, err := strconv.Atoi(resp.Status[:3])
		if err == nil {
			if numStatus < 400 {
				println("[+] Found: " + fuzzUrl)
			} else {
				println("[-] Not found")
			}
		}
	}
}
