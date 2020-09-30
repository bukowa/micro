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

import "sync"

// Counter tracks number of times something has happened.
// It has to be safe to use by multiple goroutines.
type Counter interface {
	Add(int)
	Size() int
}

// NewCounter creates new Counter.
func NewCounter() Counter {
	return &BaseCounter{}
}

// BaseCounter implements Counter.
type BaseCounter struct {
	sync.RWMutex
	n int
}

// Add adds n to internal counter value.
func (c *BaseCounter) Add(n int) {
	defer c.Unlock()
	c.Lock()
	c.n += n
}

// Size returns internal counter value.
func (c *BaseCounter) Size() int {
	defer c.RUnlock()
	c.RLock()
	return c.n
}
