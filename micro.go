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
	sync.Mutex
	sync.WaitGroup
	task Task

	// when event has occurred its sent on this channel
	eventC chan eventMessage

	// each event can have multiple listeners
	eventListeners  map[Event][]eventListener

	// events from callers
	notifyC chan eventMessage

	// callers sending events
	callers map[Event][]*Micro

	// receivers receive events
	receivers map[Event][]*Micro

	size    int
	started bool
	stop      chan struct{}
	ctx     context.Context
}

type eventListener struct {
	sync.WaitGroup
	task     Task
	receiver *Micro
	requestC chan eventMessage
	responseC chan eventMessage
	stopC    chan struct{}
}

func (el *eventListener) start() {
	el.Add(1)
	go func() {
		defer el.Done()
		for {
			select {
			case <- el.stopC:
				return
			case msg := <- el.requestC:
				event := el.task(msg.event, msg.caller, el.receiver)
				el.responseC <- eventMessage{
					event:  event,
					caller: el.receiver,
				}
			}
		}
	}()
}

type eventMessage struct {
	event Event
	caller *Micro
}

// Hook is a group of functions executed on Event.
type Hook interface {
	Register(m *Micro) map[Event][]Task
}

// Task is a function executed by Micro.
// If task execution is a result of notified event,
// then parent is set to Micro that issued that event.
// In case Task execution is a result of normal flow, then Event e is set to Noop.
// Task can be set to nil - if that happens - Micro acts only as Event listener.
type Task func(e Event, caller *Micro, receiver *Micro) Event

var (
	StartReceiver = Task(func(e Event, c *Micro, r *Micro) Event {r.Start();return Noop})
	StopReceiver = Task(func(e Event, c *Micro, r *Micro) Event {r.Stop();return Noop})
	WaitReceiver = Task(func(e Event, c *Micro, r *Micro) Event {r.Wait();return Noop})
	StartCaller = Task(func(e Event, c *Micro, r *Micro) Event {c.Start();return Noop})
	StopCaller = Task(func(e Event, c *Micro, r *Micro) Event {c.Stop();return Noop})
	WaitCaller = Task(func(e Event, c *Micro, r *Micro) Event {c.Wait();return Noop})
)

// NewMicro creates new Micro.
// Size is equal to number of goroutines spawned when it's started.
// Each spawned goroutine runs provided task until Micro is stopped.
// Stop occurs when Stop is called or context is done.
func NewMicro(size int, task Task, hooks ...Hook) *Micro {
	m := &Micro{
		task:      task,
		size:      size,
		eventC: make(chan eventMessage),
		stop:      make(chan struct{}, size),
		receivers: map[Event][]*Micro{},
		ctx:       context.Background(),
	}
	m.registerHooks(hooks...)
	return m
}

// NewMicroWithContext creates new Micro with context ctx.
func NewMicroWithContext(ctx context.Context, size int, task Task, hooks ...Hook) *Micro {
	m := NewMicro(size, task, hooks...)
	m.ctx = ctx
	return m
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

	// send event
	m.eventC <- eventMessage{
		event:  BeforeStart,
		caller: m,
	}

	// start notifies listener
	m.WaitGroup.Add(1)
	go func() {
		defer m.WaitGroup.Done()
		for {
			select {
			case <- m.ctx.Done():
				return
			case _, ok := <- m.notifyC:
				if !ok {
					return
				}

			}
		}
	}()

	// spawn goroutines
	if m.task != nil{

		// channel that each goroutines notifies when it starts
		var started = make(chan struct{}, m.size)

		m.forSize(func(i int, m *Micro) {
			m.WaitGroup.Add(1)
			go func(i int) {
				defer m.WaitGroup.Done()

				// receivers channel about start
				started <- struct{}{}

				// run task in loop
				for {
					select {
					case <-m.stop:
						return
					case <-m.ctx.Done():
						return
					default:
						event := m.task(Noop, m, m)
						m.eventC <- eventMessage{
							event: event,
							caller: m,
						}
					}
				}
			}(i)
		})

		// wait for all goroutines
		m.forSize(func(i int, m *Micro) {
			<-started
		})
		close(started)

	}

	// mark as started
	m.started = true

	// send event
	m.eventC <- eventMessage{
		event:  AfterStart,
		caller: m,
	}
	return true
}

// Started returns boolean indicating if micro is started.
func (m *Micro) Started() bool {
	defer m.Unlock()
	m.Lock()
	return m.started
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
	m.runHooks(BeforeStop)

	// receivers each goroutine to stop
	m.forSize(func(i int, m *Micro) {
		m.stop <- struct{}{}
	})

	// close event receiver
	close(m.notifyC)

	// run hooks
	m.runHooks(AfterStop)

	// mark as not started
	m.started = false
	return true
}

// StopAfter waits d time.Duration before calling Stop.
func (m *Micro) StopAfter(d time.Duration) {
	t := time.NewTimer(d)
	go func() {
		for {
			select {
			case <-t.C:
				t.Stop()
				m.Stop()
			}
		}
	}()
}

// Wait waits for Micro to finish.
func (m *Micro) Wait() {
	defer m.Unlock()
	m.Lock()
	m.runHooks(BeforeWait)
	m.WaitGroup.Wait()
	m.runHooks(AfterWait)
}

// WaitFor waits d time.Duration before calling Stop.
func (m *Micro) WaitFor(d time.Duration) {
	m.StopAfter(d)
	m.Wait()
}
// OnEvent registers task that runs on event.
func (m *Micro) OnEvent(event Event, task ...Task) *Micro {
	defer m.Unlock()
	m.Lock()
	m.registerHooks(&hook{map[Event][]Task{event: task}})
	return m
}

// Notify registers event sent to micro when it happens.
func (m *Micro) Notify(micro *Micro, event Event) {
	defer m.Unlock()
	m.Lock()
	m.registerReceiver(micro, event)
}

// UNotify cancels event notification for micro.
func (m *Micro) UNotify(micro *Micro, event Event) {
	defer m.Unlock()
	m.Lock()
	m.unotify(micro, event)
}


func (m *Micro) registerReceiver(micro *Micro, event Event) {
	m.notifyC = make(chan Event, len(m.notifyC)+1)
	if receivers, ok := m.receivers[event]; !ok {
		m.receivers[event] = append(receivers, micro)
		return
	}
	m.receivers[event] = []*Micro{micro}
}

func (m *Micro) unotify(micro *Micro, event Event) {
	if receivers, ok := m.receivers[event]; ok {
		for i, receiver := range receivers {
			if receiver == micro {
				m.receivers[event] = append(receivers[:i], receivers[i+1:]...)
			}
		}
	}
}

func (m *Micro) runHooks(e Event) {
	if hooks, ok := m.hooks[e]; ok {
		for _, hook := range hooks {
			hook(e, m, m)
		}
	}
	if receivers, ok := m.receivers[e]; ok {
		for _, receiver := range receivers {
			if cap(receiver.notifyC) > 0 {
				receiver.notifyC <- e
			}
		}
	}
}

func (m *Micro) registerHooks(hooks ...Hook) {
	for _, hook := range hooks {
		for event, tasks := range hook.Register(m) {
			if h, ok := m.hooks[event]; ok {
				m.hooks[event] = append(h, tasks...)
				continue
			}
			m.hooks[event] = tasks
		}
	}
}

// forSize runs function f Micro.size number of times.
func (m *Micro) forSize(f func(int, *Micro)) {
	for i := 0; i < m.size; i++ {
		f(i, m)
	}
}

type hook struct {
	h map[Event][]Task
}

func (h *hook) Register(*Micro) map[Event][]Task {
	return h.h
}
