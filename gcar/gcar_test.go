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
