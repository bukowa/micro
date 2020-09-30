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

// Event represents something that can happen.
type Event string

const (
	// Start happens when Crawler.Start is called.
	Start Event = "start"
	// Stop happens when Crawler.Stop is called.
	Stop Event = "stop"
	// Wait happens when Crawler.Wait is called.
	Wait Event = "wait"
	// Started happens after Crawler.Start is called and all goroutines have started.
	Started Event = "started"
	// Stopped happens after Crawler.Stop is called and all goroutines have returned.
	Stopped Event = "stopped"

	// RequestEvent happens after Crawler received a Request from the Queue.
	RequestEvent Event = "request"
	// ResponseEvent happens just before Crawler sends Response to the Queue.
	ResponseEvent Event = "response"
)
