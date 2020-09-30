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

// Tracker tracks count of requests and responses.
type Tracker interface {
	Requests() Counter
	Responses() Counter
	Errors() Counter
}

// NewTracker creates new Tracker.
func NewTracker() Tracker {
	return &BaseTracker{
		requests:  NewCounter(),
		responses: NewCounter(),
		errors:    NewCounter(),
	}
}

// BaseTracker implements Tracker.
type BaseTracker struct {
	requests  Counter
	responses Counter
	errors    Counter
}

// Requests returns Counter.
func (t *BaseTracker) Requests() Counter {
	return t.requests
}

// Responses returns Counter.
func (t *BaseTracker) Responses() Counter {
	return t.responses
}

// Errors returns Counter.
func (t *BaseTracker) Errors() Counter {
	return t.errors
}
