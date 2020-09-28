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

import (
	"context"
	"sync"
	"time"
)

// Micro is a micro thing.
// You can start it.
// You can stop it.
// You can wait for it to finish.
type Micro struct {
	sync.RWMutex
	sync.WaitGroup
	size  int
	task  Task
	stop  chan struct{}
	ctx   context.Context
	hooks map[Event][]func()

	started bool
}

// Task is a function executed by Micro.
type Task func(i int, m *Micro) (event Event)

// NewMicro creates new Micro.
// Size is equal to number of goroutines spawned when it's started.
// Each spawned goroutine runs provided task.
func NewMicro(size int, task Task, hooks ...Hook) *Micro {

	m := &Micro{
		size:  size,
		task:  task,
		stop:  make(chan struct{}, size),
		ctx:   context.Background(),
		hooks: make(map[Event][]func()),
	}
	m.registerHooks(hooks...)
	return m
}

// NewMicroWithContext creates new Micro with context.
func NewMicroWithContext(ctx context.Context, size int, task Task, hooks ...Hook) *Micro {
	m := NewMicro(size, task, hooks...)
	m.ctx = ctx
	return m
}

// Stop sends signals to stop all goroutines.
func (m *Micro) Stop() bool {

	// obtain lock
	defer m.Unlock()
	m.Lock()

	// do not stop if not started
	if !m.started {
		return false
	}

	// run hooks
	m.RunHooks(BeforeStop)

	// notify each goroutine to stop
	m.ForSize(func(i int, m *Micro) {
		m.stop <- struct{}{}
	})

	// run hooks
	m.RunHooks(AfterStop)

	// mark as not started
	m.started = false
	return true
}

// Started returns boolean indicating if micro is started.
func (m *Micro) Started() bool {
	defer m.Unlock()
	m.Lock()
	return m.started
}

// Start starts Micro.size of goroutines.
func (m *Micro) Start() bool {

	// obtain lock
	defer m.Unlock()
	m.Lock()

	// skip if started
	if m.started {
		return false
	}

	// run hooks
	m.RunHooks(BeforeStart)

	// spawn goroutines
	if m.task != nil {

		// channel that each goroutines notifies when it starts
		var started = make(chan struct{}, m.size)

		m.ForSize(func(i int, m *Micro) {
			m.WaitGroup.Add(1)
			go func(i int) {
				defer m.WaitGroup.Done()

				// notify channel about start
				started <- struct{}{}

				// run task in loop
				for {
					select {
					case <-m.stop:
						return
					case <-m.Context().Done():
						return
					default:
						m.RunHooks(m.task(i, m))
					}
				}
			}(i)
		})

		// wait for all goroutines
		m.ForSize(func(i int, m *Micro) {
			<-started
		})

	}
	// mark as started
	m.started = true

	// run hooks
	m.RunHooks(AfterStart)
	return true
}

// Wait waits for Micro to finish.
func (m *Micro) Wait() {
	m.RunHooks(BeforeWait)
	m.WaitGroup.Wait()
	m.RunHooks(AfterWait)
}

// WaitFor waits d time.Duration before calling Stop.
func (m *Micro) WaitFor(d time.Duration) {
	t := time.NewTimer(d)
	defer func() {
		t.Stop()
		m.Stop()
		m.Wait()
	}()
	for {
		select {
		case <-t.C:
			return
		}
	}
}

// Context returns Micro context.
func (m *Micro) Context() context.Context {
	return m.ctx
}

// ForSize runs function f Micro.size number of times.
func (m *Micro) ForSize(f func(int, *Micro)) {
	for i := 0; i < m.size; i++ {
		f(i, m)
	}
}

func (m *Micro) RegisterHooks(hooks map[Event][]func(m *Micro) func()) {
	h := map[Event][]func(){}
	for event, funcs := range hooks {
		for _, f := range funcs {
			h[event] = append(h[event], f(m))
		}
	}
	m.registerHooks(&hook{h: h})
}

// RunHooks runs Event hooks registered on Micro.
func (m *Micro) RunHooks(e Event) {
	if hooks, ok := m.hooks[e]; ok {
		for _, hook := range hooks {
			hook()
		}
	}
}

// OnEvent registers func executed on Event.
func (m *Micro) OnEvent(event Event, f func(m *Micro) func()) {
	m.registerHooks(m.createHook(map[Event][]func(m *Micro) func(){
		event: {f},
	}))
}

func (m *Micro) registerHooks(hooks ...Hook) {
	defer m.Unlock()
	m.Lock()
	for _, hook := range hooks {
		for event, funcs := range hook.Register(m) {
			if h, ok := m.hooks[event]; ok {
				m.hooks[event] = append(h, funcs...)
				continue
			}
			m.hooks[event] = funcs
		}
	}
}


func (m *Micro) createHook(hooks map[Event][]func(m *Micro) func()) Hook {
	h := map[Event][]func(){}
	for event, funcs := range hooks {
		for _, f := range funcs {
			h[event] = append(h[event], f(m))
		}
	}
	return &hook{h:h}
}