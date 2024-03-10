package main

import "bytes"

type cache struct {
}

var Cache = &cache{}

func (c *cache) Put() {

}

func (c *cache) Get(key string) (*bytes.Buffer, error) {
	return nil, nil
}
