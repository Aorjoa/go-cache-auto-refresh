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
