package micro_test

import (
	"context"
	. "github.com/bukowa/micro"
	"log"
	"testing"
	"time"
)

const (
	SuccessEvent Event = "success"
	FailedEvent Event = "failed"
	DoneEvent Event = "done"
)

func BenchmarkNew(t *testing.B) {
	micro := New(1, func(i int, m *Micro) Event {
		return SuccessEvent
	})
	micro.Start()
	defer testTime(t, time.Now(), time.Millisecond*100, 10)
	micro.WaitFor(time.Millisecond*100)
}

func TestNewWithContext(t *testing.T) {
	defer testTime(t, time.Now(), time.Second/2, 50)
	ctx, canc := context.WithCancel(context.Background())
	m := NewWithContext(ctx, 1, func(i int, m *Micro) Event {
		return SuccessEvent
	})
	m.Start()
	m.StopAfter(time.Second)
	time.Sleep(time.Second/2)
	canc()
	m.Wait()
}

func TestMicro_Hooks(t *testing.T) {

	micro1 := New(1, func(i int, m *Micro) Event {
		log.Print("hello from micro1!")
		return DoneEvent
	})

	micro1.RegisterHooks(map[Event][]func(*Micro) func(){
		DoneEvent: {
			func(micro *Micro) func() {
				return func() {
					log.Print("done, stop!")
					micro.Stop()
				}
			},
		},
		BeforeStop: {
			func(micro *Micro) func() {
				return func() {
					log.Print("stop called!")
				}
			},
		},
	})
	log.Print(micro1.Start())
	micro1.WaitFor(time.Millisecond)
	log.Print(micro1.Start())
	micro1.Wait()
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
