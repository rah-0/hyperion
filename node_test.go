package main

import (
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	go func() {
		err := run()
		if err != nil {
			t.Error(err)
		}
	}()

	//TODO: check handleConnection
	
	time.Sleep(1 * time.Second)
}
