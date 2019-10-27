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

func Update(c *Cache, key string, value interface{}) {
	c.item[key] = value
}

func update1(key string, pipe chan *Cache) {
	i := 0
	for {
		ball := <-pipe
		// key, value, error
		Update(ball, key, fmt.Sprintf("update key:%s=%d", key, i))
		i++
		pipe <- ball
	}
}

func update2(key string, pipe chan *Cache) {
	i := 111
	for {
		ball := <-pipe
		Update(ball, key, fmt.Sprintf("udpate2 key:%s=%d", key, i))
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

func wrapper(key string, f func() (interface{}, error)) Updater {
	i := 111
	return func(key string, pipe chan *Cache) {
		fmt.Println("update:", key)
		value, err := f()
		if err != nil {
			return
		}

		select {
		case ball := <-pipe:
			Update(ball, key, fmt.Sprintf("%s=%d", value, i))
			i++
			pipe <- ball
		}
	}
}

// executer
var ff = map[string]Updater{
	"name": update2,
}

// Add key func() interface error
func Add(key string, f func() (interface{}, error)) {
	ff[key] = wrapper(key, f)
}

type Updater func(key string, pipe chan *Cache)

var ip int

func main() {
	pipe := make(chan *Cache)

	Add("myip", func() (interface{}, error) {
		return fmt.Sprintf("myip key:%s=", "myip"), nil
	})

	// tricker controller
	go func(d time.Duration) {
		ticker := time.NewTicker(d)
		for {
			select {
			case <-ticker.C:
				for key, f := range ff {
					go f(key, pipe)
				}
			}
		}
	}(800 * time.Millisecond)

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
