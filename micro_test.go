package micro_test

import (
	"context"
	. "github.com/bukowa/micro"
	"log"
	"testing"
	"time"
)

func BenchmarkNew(t *testing.B) {
	micro := New(1, func(i int, m *Micro) {})
	micro.Start()
	defer testTime(t, time.Now(), time.Millisecond*100, 10)
	micro.WaitFor(time.Millisecond*100)
}

func TestNewWithContext(t *testing.T) {
	defer testTime(t, time.Now(), time.Second/2, 50)
	ctx, canc := context.WithCancel(context.Background())
	m := NewWithContext(ctx, 1, func(i int, m *Micro) {
	})
	m.Start()
	m.StopAfter(time.Second)
	time.Sleep(time.Second/2)
	canc()
	m.Wait()
}

func TestMicro_Hooks(t *testing.T) {
	m := New(1, func(i int, m *Micro) {})
	m.RegisterHooks(map[Event][]func(m *Micro) func(){
		BeforeStart: {
			func(m *Micro) func() {
				return func() {
					log.Print("before start!")
				}
			},
			func(m *Micro) func() {
				return func() {
					log.Print("before start2!")
				}
			},
		},
	})
	m.RegisterHooks(map[Event][]func(m *Micro) func(){
		BeforeStart: {
			func(m *Micro) func() {
				return func() {
					log.Print("before start3!")
				}
			},
		},
	})
	m.Start()
	m.WaitFor(time.Millisecond*100)
}

func testTime(t testing.TB, start time.Time, want, vary time.Duration) {
	stop := time.Since(start)
	max := want + want / vary
	min := want - want/vary
	if max < stop {
		t.Error(max, stop)
	} else {
		t.Log(max, stop)
	}
	if min > stop {
		t.Error(min, stop)
	} else {
		t.Log(min, stop)
	}
}
