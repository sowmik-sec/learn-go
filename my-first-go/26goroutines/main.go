package main

import (
	"fmt"
	"net/http"
	"sync"
	// "time"
)

var wg sync.WaitGroup // pointers

var signals = []string{"test"}
var mut sync.Mutex // pointer

func main() {
	// greeter("Hello")
	// greeter("World")

	websiteList := []string{
		// "https://lco.dev",
		"https://go.dev",
		"https://google.com",
		"https://github.com",
	}

	for _, web := range websiteList {
		go getStatusCode(web)
		wg.Add(1)
	}
	wg.Wait()
	fmt.Println(signals)
}

// func greeter(s string) {
// 	for i := 0; i < 5; i++ {
// 		time.Sleep(3 * time.Millisecond)
// 		fmt.Println(s)
// 	}
// }

func getStatusCode(endpoint string) {
	defer wg.Done()
	res, err := http.Get(endpoint)

	if err != nil {
		fmt.Println("OOPS in endpoint, ", err)

	} else {
		mut.Lock()
		signals = append(signals, endpoint)
		mut.Unlock()
		fmt.Printf("%d status code for website %s\n", res.StatusCode, endpoint)
	}
}
