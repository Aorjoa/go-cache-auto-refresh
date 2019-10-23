package main

import (
	"log"

	"github.com/aorjoa/go-cache-auto-refresh/gcar"
)

func main() {
	cache := gcar.New()
	cache.Set("key", "value")
	val, ok := cache.Get("key")
	if !ok {
		log.Print("something went wrong")
	}
	log.Printf("try to add cache [key] : %v", val)
}
