package gcar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitSource(t *testing.T) {
	t.Run("should send sources object into sourcePipeline channel", func(t *testing.T) {
		pipe := make(chan sources, 1)
		initSources(pipe)

		ss := <-pipe

		assert.NotNil(t, ss)
	})
}

func TestAdd(t *testing.T) {
	t.Run("should add source function into collection", func(t *testing.T) {
		key := "key"
		f := wrapper(key, func() (interface{}, error) {
			return "data", nil
		})
		dataPipe := make(chan *gCache, 1)
		sourcePipe := make(chan sources, 1)
		dataPipe <- &gCache{items: map[string]interface{}{}}
		sourcePipe <- sources{}

		add(key, f, dataPipe, sourcePipe)
		ss := <-sourcePipe

		assert.NotNil(t, ss[key])
	})
}
