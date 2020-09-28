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
package micro

// Event represents named thing that happens.
type Event string

const (
	// Noop event tells micro nothing.
	Noop Event = "event_noop"
	// Start event tells micro to start.
	Start Event = "event_start"
	// Stop event tells micro to stop.
	Stop Event = "event_stop"
	// Wait event tells micro to wait.
	Wait Event = "event_wait"
)

const (
	// BeforeStart happens before micro starts.
	BeforeStart Event = "event_before_start"
	// BeforeStop happens before micro stops.
	BeforeStop Event = "event_before_stop"
	// BeforeWait happens before micro waits.
	BeforeWait Event = "event_before_wait"

	// AfterStart happens after micro starts.
	AfterStart Event = "event_after_start"
	// AfterStop happens after micro stops.
	AfterStop Event = "event_after_stop"
	// AfterWait happens after micro waits.
	AfterWait Event = "event_after_wait"
)
