package main

import "github.com/iddxc/memodb/storage"

func New(filename string, timePeriod int) *storage.Store {
	return storage.InitStore(filename, timePeriod)
}
