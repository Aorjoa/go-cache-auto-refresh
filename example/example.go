package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aorjoa/go-cache-auto-refresh/gcar"
)

// Example struct call get

func main() {
	gcar.Add("nong", func() (interface{}, error) {
		return "anuchitO", nil
	})
	gcar.Add("myip", caller)
	gcar.Add("myip1", caller)
	gcar.Add("myip2", caller)
	gcar.Set("myip4", "dummyIP")

	v, ok := gcar.Get("nong")
	fmt.Println("nong:", v, "found:", ok)
	v, ok = gcar.Get("myip")
	fmt.Println("myip:", v, "found:", ok)

	gcar.UpdateTick(800 * time.Millisecond)
	time.Sleep(5 * time.Second)
	v, ok = gcar.Get("myip")
	fmt.Println("myip:", v, "found:", ok)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown
	log.Println("good bye.")
}

func caller() (interface{}, error) {
	fmt.Println("caller")
	response, err := http.Get("https://httpbin.org/ip")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, err := ioutil.ReadAll(response.Body)
	return string(data), err
}
