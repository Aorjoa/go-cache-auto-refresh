package gcar

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
	set(key, value, pipeline)
}

func set(key string, value interface{}, pipe chan *gCache) {
	select {
	case c := <-pipe:
		c.set(key, value)
		pipe <- c
	}
}

// TODO: with context - timeout
// func GetWithContext(ctx context.Context) {}
