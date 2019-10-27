package gcar

import (
	"log"
	"time"
)

var pipeline = make(chan *gCache, 1)

func init() {
	initializer(pipeline)
}

func initializer(pipe chan *gCache) {
	gc := &gCache{
		items: map[string]interface{}{},
	}
	pipe <- gc
}

func Get(key string) (value interface{}, isExist bool) {
	return get(key, pipeline)
}

func get(key string, pipe chan *gCache) (value interface{}, isExist bool) {
	select {
	case c := <-pipe:
		value, isExist = c.items[key]
		pipe <- c
		return
	}
}

func Set(key string, value interface{}) {
	select {
	case c := <-pipeline:
		c.set(key, value)
		pipeline <- c
	}
}

// TODO: with context - timeout
// func GetWithContext(ctx context.Context) {}

type Source func() (interface{}, error)
type updater func(chan *gCache)

type gCache struct {
	items map[string]interface{}
}

func (gc *gCache) set(key string, value interface{}) {
	gc.items[key] = value
}

var ff = map[string]updater{}

// Add is register a function get Source of value
// when calll Add it will excute Source function to get value for first time.
func Add(key string, s Source) {
	f := wrapper(key, s)
	add(key, f, pipeline)
}

func add(key string, f updater, pipe chan *gCache) {
	defer f(pipe)
	ff[key] = f
}

func wrapper(key string, s Source) updater {
	return func(pipe chan *gCache) {
		value, err := s()
		if err != nil {
			// TODO: refactor handler error
			log.Println("not update date: because retrive date error key:", key, ",error:", err)
			return
		}

		select {
		case c := <-pipe:
			c.set(key, value)
			pipe <- c
		}
	}
}

// UpdateTick is update new data from Source every duration it take.
func UpdateTick(d time.Duration) {
	// TODO: extract to private function
	go func(d time.Duration) {
		ticker := time.NewTicker(d)
		for {
			select {
			case <-ticker.C:
				for _, f := range ff {
					go f(pipeline)
				}
			}
		}
	}(d)
}
