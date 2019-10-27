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

	t.Run("Should be return nil, false for not found item", func(t *testing.T) {
		c := &gCache{
			items: map[string]interface{}{},
		}
		pipe := make(chan *gCache, 1)
		pipe <- c

		value, ok := get("nokey", pipe)

		assert.Nil(t, value)
		assert.False(t, ok)
	})
}

func TestSet(t *testing.T) {
	t.Run("should set value into cache object by using a key", func(t *testing.T) {
		old := &gCache{
			items: map[string]interface{}{},
		}
		pipe := make(chan *gCache, 1)
		pipe <- old

		set("name", "AnuchitO", pipe)
		fresh := <-pipe

		assert.Equal(t, "AnuchitO", fresh.items["name"])
	})
}
