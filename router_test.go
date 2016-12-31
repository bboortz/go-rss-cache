package main

import (
	"testing"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)



func TestNewRouter(t *testing.T) {
	assert := assert.New(t)
	router := NewRouter()
	spew.Dump(router)
	assert.NotNil(router)
}

