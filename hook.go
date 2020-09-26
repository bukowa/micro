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

// Hook is a group of functions executed on Event.
type Hook interface {
	Register(m *Micro) map[Event][]func()
}

// NewHook creates new Hook.
func NewHook(micro *Micro, hooks map[Event][]func(m *Micro) func()) Hook {
	h := map[Event][]func(){}
	for event, funcs := range hooks {
		for _, f := range funcs {
			h[event] = append(h[event], f(micro))
		}
	}
	return &hook{h:h}
}

type hook struct {
	h map[Event][]func()
}

func (h *hook) Register(m *Micro) map[Event][]func() {
	return h.h
}
