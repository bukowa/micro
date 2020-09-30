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
	"time"
)

// WaitUnknownTime is a function for crawlers that cannot determine
// amount of requests or time they will take to complete crawling.
// This can happen when all urls of a webpage should be crawled.
// If crawlers did not process any requests or responses for an exceeded
// period of time and it means that crawlers should be stopped.
// After each 'tick' WaitUnknownTime checks len of responses and requests that crawlers made.
// If these numbers didn't change between 'tick', n is increased by 1, otherwise n is zeroed.
// When n equals `count` and channels of Response and Request are of len 0 - crawlers are stopped.
// It mean's all Response object have to be taken out from the Queue before Crawler can stop.
var WaitUnknownTime = func(c *Crawler, count int, tick time.Duration) {

	var stopped = make(chan struct{}, 1)
	c.OnEvent(Stop, func(e Event, c *Crawler) {
		stopped <- struct{}{}
	})

	var collect = func() (req int, res int) {
		return c.Requests().Size(), c.Responses().Size()
	}

	var changed = func(req int, res int) bool {
		reqN, resN := collect()
		if reqN == req && resN == res {
			return false
		}
		return true
	}

	go func() {

		var n int
		for {
			select {
			case <-stopped:
				return
			default:
				req, res := collect()
				time.Sleep(tick)
				if !changed(req, res) {
					n++
				} else {
					n = 0
				}
				if n >= count && len(c.Response()) == 0 && len(c.Request()) == 0 {
					c.Stop()
					return
				}
			}
		}
	}()
	c.Wait()
}
