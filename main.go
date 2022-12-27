package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"sync"
)

func checkAndSaveBody(url string, wg *sync.WaitGroup, chn chan string) {
	resp, err := http.Get(url)

	if err != nil {
		strError := fmt.Sprintf("%s is Down! \n", url)
		strError += fmt.Sprintf("Error: %v", err)

		chn <- strError
	} else {
		defer resp.Body.Close()

		strResponse := fmt.Sprintf("%s -> Status Code: %d \n", url, resp.StatusCode)

		if resp.StatusCode == 200 {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				strError := fmt.Sprintf("Read Body err: %v", err)
				chn <- strError
			}

			fileName := strings.Split(url, "//")[1] // this will take domain only except `http:`
			fileName += ".txt"

			strResponse += fmt.Sprintf("Writing response body to %s\n", fileName)

			err = ioutil.WriteFile(fileName, bodyBytes, 0664)
			if err != nil {
				strError := fmt.Sprintf("Write file err: %v", err)
				chn <- strError
			}
		}

		strResponse += fmt.Sprintf("%s is UP \n", url)

		chn <- strResponse
	}

	wg.Done()
}

func main() {
	urls := []string{"https://golang1.org", "https://www.google.com", "https://www.medium.com"}

	// Create wait group
	var wg sync.WaitGroup

	// Create channel
	chn := make(chan string)

	wg.Add(len(urls))

	for _, url := range urls {
		go checkAndSaveBody(url, &wg, chn)
	}

	fmt.Println("No. of Goroutines:", runtime.NumGoroutine())

	// Receive message from channel
	for i := 0; i < len(urls); i++ {
		fmt.Println(<-chn)
	}

	wg.Wait()
}
