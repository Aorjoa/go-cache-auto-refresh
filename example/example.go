package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	cache.CallFunctionThenCache("keyAPI", caller)
	val, ok = cache.Get("keyAPI")
	if !ok {
		log.Print("something went wrong")
	}
	log.Printf("try to add cache [keyAPI] : %v", val)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown
	log.Println("good bye.")
}

func caller() (interface{}, error) {
	response, err := http.Get("https://httpbin.org/ip")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, err := ioutil.ReadAll(response.Body)
	return string(data), err
}
