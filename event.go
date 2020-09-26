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
	// BeforeStart happens before micro starts.
	BeforeStart Event = "before_start"
	// BeforeStop happens before micro stops.
	BeforeStop Event = "before_stop"
	// BeforeWait happens before micro waits.
	BeforeWait Event = "before_wait"

	// AfterStart happens after micro starts.
	AfterStart Event = "after_start"
	// AfterStop happens after micro stops.
	AfterStop Event = "after_stop"
	// AfterWait happens after micro waits.
	AfterWait Event = "after_wait"
)
