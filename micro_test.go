package micro_test

import (
	"context"
	"fmt"
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
	micro := NewMicro(1, func(i int, m *Micro) Event {
		return SuccessEvent
	})
	micro.Start()
	defer testTime(t, time.Now(), time.Millisecond*100, 10)
	micro.WaitFor(time.Millisecond*100)
}

func TestNewWithContext(t *testing.T) {
	defer testTime(t, time.Now(), time.Second/2, 50)
	ctx, canc := context.WithCancel(context.Background())
	m := NewMicroWithContext(ctx, 1, func(i int, m *Micro) Event {
		return SuccessEvent
	})
	m.Start()
	time.Sleep(time.Second/2)
	canc()
	m.Wait()
}

func TestMicro_Hooks(t *testing.T) {

	micro1 := NewMicro(1, func(i int, m *Micro) Event {
		log.Print("hello from micro1!")
		return DoneEvent
	})

	micro1.RegisterHooks(Hooks{
		DoneEvent: {
			func(micro *Micro) func() {
				return func() {
					log.Print("done, stop!")
					log.Print(micro.Stop())
					log.Print("stopped")
				}
			},
		},
		BeforeWait: {
			func(micro *Micro) func() {
				return func() {
					log.Print("before wait")
				}
			},
		},
		AfterWait: {
			func(micro *Micro) func() {
				return func() {
					log.Print("after wait")
				}
			},
		},
		BeforeStop: {
			func(micro *Micro) func() {
				return func() {
					log.Print("before stop")
				}
			},
		},
	})
	micro1.Start()
	micro1.WaitFor(time.Second)
}

func TestDoOnce(t *testing.T) {
	micro := NewMicro(5, DoOnce(func(i int, m *Micro) {
		log.Println("hello!")
	}))
	micro.Start()
	micro.WaitFor(time.Second)
}

func ExampleMicro_Start() {
	const Success Event = "success"
	const Failed Event = "failed"

	micro1 := NewMicro(1, func(i int, m *Micro) Event {
		return Success
	})

	micro2 := NewMicro(1, func(i int, m *Micro) Event {
		fmt.Println("hello from micro2")
		return Failed
	})

	micro1.RegisterHooks(map[Event][]func(m *Micro) func(){
		AfterStart: {
			func(m *Micro) func() {
				return func() {
					fmt.Println("hello from micro")
					fmt.Println("starting micro2")
					micro2.Start()
				}
			},
		},
	})

	micro2.RegisterHooks(map[Event][]func(m *Micro) func(){
		Failed: {
			func(m *Micro) func() {
				fmt.Println("hook registered")
				return func() {
					fmt.Println("stopping all micros")
					micro1.Stop()
					micro2.Stop()
				}
			},
		},
	})

	micro1.Start()
	micro1.Wait()
	micro2.Wait()
	// Output: hook registered
	//hello from micro
	//starting micro2
	//hello from micro2
	//stopping all micros
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
