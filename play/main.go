package main

import (
	"context"
	"log"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var c = make(chan struct{}, 1)
	var ctx = context.Background()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-c:
				return
			case <-ctx.Done():
				return
			default:
				log.Print("hello!")
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.NewTicker(1)
		for {
			select {
			case <-t.C:
				time.Sleep(time.Second)
				log.Print(c)
			}
		}
	}()

	wg.Wait()
}
