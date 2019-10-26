package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	c := New()

	u := c.Get("user")

	assert.Equal(t, u, "anuchito")
}

func TestGetNoFoundInItem(t *testing.T) {
	c := New()
	go PeriodicallyUpdate(c)
	time.Sleep(20 * time.Second)
	u, _ := c.Get("myip")

	assert.Equal(t, u, "fresh data")
	time.Sleep(20 * time.Second)
}
