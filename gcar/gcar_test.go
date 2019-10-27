package gcar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializer(t *testing.T) {
	t.Run("should send cache object initialized into pipline", func(t *testing.T) {
		pipe := make(chan *gCache, 1)
		initializer(pipe)

		c := <-pipe

		assert.NotNil(t, c.items)
	})
}

func TestGet(t *testing.T) {
	t.Run("Should be return value from by key", func(t *testing.T) {
		c := &gCache{
			items: map[string]interface{}{"name": "AnuchitO"},
		}
		pipe := make(chan *gCache, 1)
		pipe <- c

		value, ok := get("name", pipe)

		assert.True(t, ok)
		assert.Equal(t, "AnuchitO", value)
	})
}
