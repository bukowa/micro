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

// Queue represents communication betwee caller and the crawler.
// Crawler performs requests received from Request() method.
// Once the request is completed its send into Response() method.
type Queue interface {
	Response() chan Response
	Request() chan Request
}

// NewQueue creates new Queue.
func NewQueue(requestSize, responseSize int) Queue {
	return BaseQueue{
		results: make(chan Response, responseSize),
		request: make(chan Request, requestSize),
	}
}

// BaseQueue implements Queue.
type BaseQueue struct {
	results chan Response
	request chan Request
}

// Response returns underyling Response channel.
func (q BaseQueue) Response() chan Response {
	return q.results
}

// Request returns underlying Request channel.
func (q BaseQueue) Request() chan Request {
	return q.request
}
