/*
Copyright © 2020 Mateusz Kurowski

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package crawler

import (
	"net/http"
	"sync"
	"time"
)

// Crawler is responsible for sending http requests.
type Crawler struct {
	Tracker
	Queue
	Logger

	sync.Mutex
	sync.WaitGroup

	size   int
	stop   chan struct{}
	sleep  time.Duration
	client *http.Client

	newRespFunc NewResponseFunc

	onRequest []func(int, *Crawler, Request) error
	onResponse []func(int, *Crawler, Response) error
	events map[Event][]func(Event, *Crawler)
}

// NewResponseFunc is a function used by crawler to create new Response.
// It can be replaced by passing WithResponseFunc as opt to NewCrawler.
type NewResponseFunc func(crawler *Crawler, took time.Duration, req Request, res *http.Response, err error) Response

// NewCrawler returns Crawler of size n.
func NewCrawler(size int, opts ...Option) *Crawler {
	c := &Crawler{
		Tracker: NewTracker(),
		Queue:   NewQueue(size, size),
		Logger:  NewLogger(),
		sleep:   time.Millisecond,
		size:    size,
		// non-buffered channel
		stop:    make(chan struct{}),
		events:  map[Event][]func(Event, *Crawler){},
		// client is modified to avoid networking problems
		// while testing with default http client there are issues
		client: &http.Client{Transport: &http.Transport{
			TLSHandshakeTimeout: time.Second * 10,
			DisableKeepAlives:   false,
			MaxIdleConns:        size,
			MaxIdleConnsPerHost: size,
			MaxConnsPerHost:     size,
			IdleConnTimeout:     time.Second * 2,
		}},
		newRespFunc: NewResponse,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Start starts crawling goroutines.
func (c *Crawler) Start() {
	c.event(Start)

	// channel that each goroutine notify when it started
	started := make(chan struct{})

	for i := 0; i < c.size; i++ {

		c.Add(1)
		go func(i int) {
			defer func() {
				c.Done()
			}()
			// notify started
			started <- struct{}{}

			RequestLoop:
			for {
				select {
				case <-c.stop:
					return
				case request := <-c.Request():

					// execute OnRequest functions
					// if any of them returns error, cancel Request
					for _, f := range c.onRequest {
						if err := f(i, c, request); err != nil {
							goto RequestLoop
						}
					}

					// increment requests count
					c.Requests().Add(1)
					c.event(RequestEvent)

					// perform http request
					start := time.Now()
					responseHTTP, err := c.client.Do(request.Request())
					took := time.Since(start)

					// create new response
					response := c.newRespFunc(c, took, request, responseHTTP, err)

					// execute OnResponse functions
					// if any of them returns error, cancel Response
					for _, f := range c.onResponse {
						if err := f(i, c, response); err != nil {
							goto RequestLoop
						}
					}

					// increment responses && error count
					c.Responses().Add(1)
					if err != nil {
						c.Errors().Add(1)
					}
					c.event(ResponseEvent)

					// send response to Queue
					c.Response() <- response
				default:
					// sleep on default case - otherwise cpu usage can be very high
					// b-c network is always a bottleneck - this is not an issue
					// otherwise set it to 0 via Crawler option WithSleep
					time.Sleep(c.sleep)
				}
			}
		}(i)
	}

	// wait for startup
	for i := 0; i < c.size; i++ {
		<-started
	}

	// execute event
	c.event(Started)
}

// Stop sends signal to notify all goroutines to return.
func (c *Crawler) Stop() {
	c.event(Stop)
	for i := 0; i < c.size; i++ {
		// non-buffered channel guarantees all goroutines
		// stopped before event is executed
		c.stop <- struct{}{}
	}
	c.event(Stopped)
}

// Wait waits for all goroutines to finish.
func (c *Crawler) Wait() {
	c.event(Wait)
	c.WaitGroup.Wait()
}

// Event registers function f executed on Event e.
func (c *Crawler) OnEvent(e Event, f func(e Event, c *Crawler)) {
	defer c.Unlock()
	c.Lock()

	if funcs, ok := c.events[e]; ok {
		c.events[e] = append(funcs, f)
		return
	}
	c.events[e] = []func(Event, *Crawler){f}
}

// OnRequest registers function f executed when Crawler received Request from Queue.
// If this function returns an error then Request is abandoned and Crawler will continue.
func (c *Crawler) OnRequest(f func(i int, c *Crawler, request Request) error) {
	defer c.Unlock()
	c.Lock()
	c.onRequest = append(c.onRequest, f)
}

// OnResponse registers function f executed before Crawler sends Response to Queue.
// If this function returns an error then Response is abandoned (not sent to Queue).
func (c *Crawler) OnResponse(f func(i int, c *Crawler, response Response) error) {
	defer c.Unlock()
	c.Lock()
	c.onResponse = append(c.onResponse, f)
}

func (c *Crawler) event(e Event) {
	defer c.Unlock()
	c.Lock()
	if funcs, ok := c.events[e]; ok {
		for _, f := range funcs {
			f(e, c)
		}
	}
}
