package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)


func fuzz(url string, lines []string, client *http.Client){

	var fuzzUrl string

	for _, line := range lines {

		fuzzUrl = url + "/" + line

		resp, err := client.Get(fuzzUrl)
		if err != nil {
			log.Fatal(err)
		}

		numStatus, err := strconv.Atoi(resp.Status[:3])
		if err == nil {
			io.Copy(ioutil.Discard, resp.Body)
			if numStatus < 400 {

				fmt.Println("[" + strconv.Itoa(numStatus) + "] Found: " + fuzzUrl)
				pattern := ".php"
				matched, err := regexp.Match(pattern, []byte(fuzzUrl))
				if err != nil {
					log.Fatal(err)
				}
				if !matched {
					fuzz(fuzzUrl, lines, client)
				}

			} else {

				fmt.Print(fuzzUrl + strings.Repeat(" ", 50))
				fmt.Print("\r")

			}
		}
	}


}

func main() {

	url:= os.Args[len(os.Args) - 1]

	file, err := os.Open("common.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	tr := &http.Transport{
		MaxIdleConns:       100,
		IdleConnTimeout:    2 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}


	fuzz(url, lines, client)


}
