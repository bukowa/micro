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
	"fmt"
	"io"
	"net/http"
	"time"
)

type Option = func(c *Crawler)

var WithClient = func(client *http.Client) Option {
	return func(c *Crawler) {
		c.client = client
	}
}

var WithLogger = func(logger Logger) Option {
	return func(c *Crawler) {
		c.Logger = logger
	}
}

var WithLoggerPrefix = func(prefix string) Option {
	return func(c *Crawler) {
		c.SetPrefix(prefix)
	}
}

var WithTracker = func(tracker Tracker) Option {
	return func(c *Crawler) {
		c.Tracker = tracker
	}
}

var WithQueue = func(queue Queue) Option {
	return func(c *Crawler) {
		c.Queue = queue
	}
}

var WithLoggerOutput = func(w io.Writer) Option {
	return func(c *Crawler) {
		c.Logger.SetOutput(w)
	}
}

var WithResponseFunc = func(f NewResponseFunc) Option {
	return func(c *Crawler) {
		c.newRespFunc = f
	}
}

var WithSleep = func(t time.Duration) Option {
	return func(c *Crawler) {
		c.sleep = t
	}
}

var WithRequestLog = func(f func(i int, c *Crawler, r Request) string) Option {
	return func(c *Crawler) {
		c.OnRequest(func(i int, c *Crawler, r Request) error {
			c.Print(f(i, c, r))
			return nil
		})
	}
}

var WithResponseLog = func(f func(i int, c *Crawler, r Response) string) Option {
	return func(c *Crawler) {
		c.OnResponse(func(i int, c *Crawler, r Response) error {
			c.Print(f(i, c, r))
			return nil
		})
	}
}


var WithDefaultRequestLog = func() Option {
	return WithRequestLog(func(i int, c *Crawler, r Request) string {
		return fmt.Sprintf("%v:request:%s", i, r.Request().URL.String())
	})
}


var WithDefaultResponseLog = func() Option {
	return WithResponseLog(func(i int, c *Crawler, r Response) string {
		return fmt.Sprintf("%v:response:%s:err:%s", i, r.Request().URL.String(), r.Error())
	})
}

var WithEventLog = func(e Event, v ...interface{}) Option {
	return func(c *Crawler) {
		c.OnEvent(e, func(e Event, c *Crawler) {
			c.Print(e, v)
		})
	}
}

var WithStartLog = func() Option {
	return WithEventLog(Start, Start)
}

var WithStopLog = func() Option {
	return WithEventLog(Stop, Stop)
}

var WithWaitLog = func() Option {
	return WithEventLog(Wait, Wait)
}

var WithStartedLog = func() Option {
	return WithEventLog(Started, Started)
}

var WithStoppedLog = func() Option {
	return WithEventLog(Stopped, Stopped)
}


var WithDefaultLog = func(c *Crawler) {
	WithDefaultRequestLog()(c)
	WithDefaultResponseLog()(c)
	WithStartLog()(c)
	WithStopLog()(c)
	WithStartedLog()(c)
	WithStoppedLog()(c)
	WithWaitLog()(c)
}