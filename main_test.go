package main

import (
	"testing"

	"github.com/matryer/is"
) 

func TestMain(t *testing.T) {
	is := is.New(t)

	_, err := newServer()

	is.NoErr(err)
}
