package main

import (
	"fmt"
	"net/url"
)

const myurl string = "https://www.youtube.com/watch?v=cl7_ouTMFh0&list=PLRAV69dS1uWQGDQoBYMZWKjzuhCaOnBpa&index=26"

func main() {
	fmt.Println("Welcome to handling URLs in golang")
	fmt.Println(myurl)

	// parsing
	result, err := url.Parse(myurl)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Scheme)
	fmt.Println(result.Host)
	fmt.Println(result.Path)
	fmt.Println(result.RawQuery)

	qparams := result.Query()
	fmt.Printf("The type of query params are %T\n", qparams)

	fmt.Println(qparams["coursename"])

	for _, val := range qparams {
		fmt.Println("Param is: ", val)
	}

	partsOfUrl := &url.URL{
		Scheme:   "https",
		Host:     "youtube.com",
		Path:     "/watch",
		RawQuery: "v=cl7_ouTMFh0&list=PLRAV69dS1uWQGDQoBYMZWKjzuhCaOnBpa&index=26",
	}
	anoterUrl := partsOfUrl.String()
	fmt.Println(anoterUrl)
}
