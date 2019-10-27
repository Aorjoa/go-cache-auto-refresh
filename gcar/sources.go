package gcar

import (
	"log"
	"time"
)

var sourcePipeline = make(chan sources, 1)

func init() {
	initSources(sourcePipeline)
}

func initSources(sourcePipe chan sources) {
	ss := sources{}
	sourcePipe <- ss
}

type Source func() (interface{}, error)
type updater func(chan *gCache)

type gCache struct {
	items map[string]interface{}
}

func (gc *gCache) set(key string, value interface{}) {
	gc.items[key] = value
}

type sources map[string]updater

func (s sources) set(key string, f updater) {
	s[key] = f
}

// Add is register a function get Source of value
// when calll Add it will excute Source function to get value for first time.
func Add(key string, s Source) {
	f := wrapper(key, s)
	add(key, f, pipeline, sourcePipeline)
}

func add(key string, f updater, dataPipe chan *gCache, sourcePipe chan sources) {
	defer f(dataPipe)
	select {
	case ss := <-sourcePipe:
		ss[key] = f
		sourcePipe <- ss
	}
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
				select {
				case ss := <-sourcePipeline:
					for _, f := range ss {
						go f(pipeline)
					}
					sourcePipeline <- ss
				}
			}
		}
	}(d)
}
