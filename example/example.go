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
)

// var loaderFunc = map[string]func(pipeCh chan *Cache){
// 	"myip": func(pipeCh chan *Cache) {
// 		fmt.Println("take out pipeCh")
// 		cc := <-pipeCh
// 		cc.Update("myip", "it's works!!!!! myip")
// 		pipeCh <- cc
// 		fmt.Println("put in pipeCh")
// 	},
// }

type Cache struct{ item map[string]interface{} }

func update1(key string, pipe chan *Cache) {
	i := 0
	for {
		ball := <-pipe
		ball.item[key] = fmt.Sprintf("update key:%s=%d", key, i)
		i++
		pipe <- ball
	}
}

func update2(key string, pipe chan *Cache) {
	i := 111
	for {
		ball := <-pipe
		ball.item[key] = fmt.Sprintf("udpate2 key:%s=%d", key, i)
		i++
		pipe <- ball
	}
}

func getCache(key string, pipe chan *Cache) {
	select {
	case ball := <-pipe:
		fmt.Println("get ball:", ball.item[key])
		pipe <- ball
	}
}

func main() {
	pipe := make(chan *Cache)
	go update1("myip", pipe)
	go update2("name", pipe)

	pipe <- &Cache{
		item: map[string]interface{}{"name": "anuchito"},
	}

	for i := 0; i < 10; i++ {
		getCache("myip", pipe)
		getCache("name", pipe)
		getCache("myip", pipe)
		getCache("name", pipe)
		getCache("myip", pipe)
		getCache("name", pipe)
		getCache("myip", pipe)
		getCache("name", pipe)
		getCache("myip", pipe)
		getCache("name", pipe)
		time.Sleep(1 * time.Second)
	}

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
