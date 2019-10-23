package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aorjoa/go-cache-auto-refresh/gcar"
)

func main() {
	cache := gcar.New()
	// simple
	cache.Set("key", "value")
	val, ok := cache.Get("key")
	if !ok {
		log.Print("something went wrong")
	}
	log.Printf("try to add cache [key] : %v", val)

	// call function then cache
	cache.CallFunctionThenCache("keyAPI", caller())
	val, ok = cache.Get("keyAPI")
	if !ok {
		log.Print("something went wrong")
	}
	log.Printf("try to add cache [keyAPI] : %v", val)

}

func caller() func() (interface{}, error) {
	return func() (interface{}, error) {
		response, err := http.Get("https://httpbin.org/ip")
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}
		data, err := ioutil.ReadAll(response.Body)
		return string(data), err
	}
}
