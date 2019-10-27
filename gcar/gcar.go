package gcar

import (
	"log"
	"time"
)

var pipeline = make(chan *gCache, 1)

type gCache struct {
	items map[string]interface{}
}

func (gc *gCache) update(key string, value interface{}) {
	gc.items[key] = value
}

func Get(key string) (value interface{}, isExist bool) {
	select {
	case c := <-pipeline:
		value, isExist = c.items[key]
		pipeline <- c
		return
	}
}

// TODO: add function set

// TODO: with context - timeout
// func GetWithContext(ctx context.Context) {}

func init() {
	gc := &gCache{
		items: map[string]interface{}{},
	}
	pipeline <- gc
}

type Source func() (interface{}, error)
type updater func(chan *gCache)

var ff = map[string]updater{}

func Add(key string, s Source) {
	f := wrapper(key, s)
	defer f(pipeline)
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
			c.update(key, value)
			pipe <- c
		}
	}
}

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
