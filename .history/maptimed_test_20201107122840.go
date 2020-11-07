package maptimed_test

import (
	"maptimed"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMapTimed(t *testing.T) {
	assert := assert.New(t)
	timeout := time.Duration(1)
	value := 1

	var m *maptimed.MapTimed
	var err error

	m, err = maptimed.NewMapTimed(0)
	assert.Nil(m)
	assert.NotNil(err)
	m, err = maptimed.NewMapTimed(timeout)
	assert.NotNil(m)
	assert.Nil(err)

	m.Set("hello", value)
	assert.Equal(value, m.Get("hello"))
	time.Sleep(2 * time.Duration(timeout) * time.Second)
	assert.Nil(m.Get("hello"))

	m.Set("hello", value)
	assert.Equal(value, m.Get("hello"))
	time.Sleep(2 * time.Duration(timeout) * time.Second)
	assert.Nil(m.Get("hello"))

}
