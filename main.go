package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
)

func checkAndSaveBody(url string, wg *sync.WaitGroup) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("%s is DOWN!\n", url)
	} else {
		defer resp.Body.Close()
		fmt.Printf("%s -> Status Code: %d \n", url, resp.StatusCode)

		if resp.StatusCode == 200 {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("Read Body err: %v", err)
			}

			fileName := strings.Split(url, "//")[1] // this will take domain only except `http:`
			fileName += ".txt"

			fmt.Printf("Writing response body to %s\n", fileName)

			err = ioutil.WriteFile(fileName, bodyBytes, 0664)
			if err != nil {
				log.Fatalf("Write file err: %v", err)
			}
		}
	}

	wg.Done()
}

func main() {
	urls := []string{"https://golang1.org", "https://www.google.com", "https://www.medium.com"}

	var wg sync.WaitGroup

	wg.Add(len(urls))

	for _, url := range urls {
		go checkAndSaveBody(url, &wg)
		fmt.Println(strings.Repeat("#", 20))
	}

	fmt.Println("No. of Goroutines:", runtime.NumGoroutine())

	wg.Wait()
}
