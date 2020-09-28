package micro_test

import (
	. "github.com/bukowa/micro"
	"testing"
)

func TestMicro_Notify(t *testing.T) {
	var n int

	var stop = Task(func(e Event, c *Micro, r *Micro) Event {
		if n == 100 {
			return Stop
		}
		return Noop
	})
	var increase = Task(func(e Event, c *Micro, r *Micro) Event {
		n ++
		return Noop
	})
	var request = Task(func(e Event, c *Micro, r *Micro) Event {
		return "add"
	})

	caller := NewMicro(1, request)
	receiver := NewMicro(1, stop)

	receiver.OnEvent("add", increase)
	receiver.OnEvent(Stop, StopReceiver, StopCaller)

	caller.Notify(receiver, "add")

	receiver.Start()
	caller.Start()

	receiver.Wait()
	caller.Wait()
}
func TestMicro_OnEventStopReceiver(t *testing.T) {
	var n int
	task := Task(func(e Event, c *Micro, r *Micro) Event {
		n ++
		return Stop
	})
	micro := NewMicro(1, task)
	micro.OnEvent(Stop, StopReceiver)
	micro.Start()
	micro.Wait()
	if n != 1 {
		t.Error(n)
	}
}

func TestMicro_TaskNil(t *testing.T) {
	micro := NewMicro(1, nil)
	micro.Start()
	micro.Stop()
	micro.Wait()
}

func TestMicro_TaskNilNotify(t *testing.T) {
	
}

//
//const (
//	SuccessEvent Event = "success"
//	FailedEvent Event = "failed"
//	DoneEvent Event = "done"
//)
//
//func BenchmarkNew(t *testing.B) {
//	micro := NewMicro(1, func(i int, m *Micro) Event {
//		return SuccessEvent
//	})
//	micro.Start()
//	defer testTime(t, time.Now(), time.Millisecond*100, 10)
//	micro.WaitFor(time.Millisecond*100)
//}
//
//func TestNewWithContext(t *testing.T) {
//	defer testTime(t, time.Now(), time.Second/2, 50)
//	ctx, canc := context.WithCancel(context.Background())
//	m := NewMicroWithContext(ctx, 1, func(i int, m *Micro) Event {
//		return SuccessEvent
//	})
//	m.Start()
//	m.StopAfter(time.Second)
//	time.Sleep(time.Second/2)
//	canc()
//	m.Wait()
//	m.Wait()
//}
//
//func TestMicro_Hooks(t *testing.T) {
//
//	micro1 := NewMicro(1, func(i int, m *Micro) Event {
//		log.Print("hello from micro1!")
//		return DoneEvent
//	})
//
//	micro1.RegisterHooks(map[Event][]func(*Micro) func(){
//		DoneEvent: {
//			func(micro *Micro) func() {
//				return func() {
//					log.Print("done, stop!")
//					micro.Stop()
//				}
//			},
//		},
//		BeforeStop: {
//			func(micro *Micro) func() {
//				return func() {
//					log.Print("stop called!")
//				}
//			},
//		},
//	})
//	log.Print(micro1.Start())
//	micro1.WaitFor(time.Millisecond)
//	log.Print(micro1.Start())
//	micro1.Wait()
//}
//
//func TestDoOnce(t *testing.T) {
//	micro := NewMicro(5, DoOnce(func(i int, m *Micro) {
//		log.Println("hello!")
//	}))
//	micro.Start()
//	micro.WaitFor(time.Second)
//}
//
//func TestStopWhen_Register(t *testing.T) {
//	const done Event = "done"
//	micro := NewMicro(5, func(i int, m *Micro) Event {
//		log.Println("done!")
//		return done
//	}, StopWhen{Event:done})
//	micro.Start()
//	micro.Wait()
//}
//
//func ExampleMicro_Start() {
//	const Success Event = "success"
//	const Failed Event = "failed"
//
//	micro1 := NewMicro(1, func(i int, m *Micro) Event {
//		return Success
//	})
//
//	micro2 := NewMicro(1, func(i int, m *Micro) Event {
//		fmt.Println("hello from micro2")
//		return Failed
//	})
//
//	micro1.RegisterHooks(map[Event][]func(m *Micro) func(){
//		AfterStart: {
//			func(m *Micro) func() {
//				return func() {
//					fmt.Println("hello from micro")
//					fmt.Println("starting micro2")
//					micro2.Start()
//				}
//			},
//		},
//	})
//
//	micro2.RegisterHooks(map[Event][]func(m *Micro) func(){
//		Failed: {
//			func(m *Micro) func() {
//				fmt.Println("hook registered")
//				return func() {
//					fmt.Println("stopping all micros")
//					micro1.Stop()
//					micro2.Stop()
//				}
//			},
//		},
//	})
//
//	micro1.Start()
//	micro1.Wait()
//	micro2.Wait()
//	// Output: hook registered
//	//hello from micro
//	//starting micro2
//	//hello from micro2
//	//stopping all micros
//}
//
//func testTime(t testing.TB, start time.Time, want, vary time.Duration) {
//	stop := time.Since(start)
//	max := want + want / vary
//	min := want - want/vary
//	if max < stop {
//		t.Error(max, stop)
//	} else {
//		t.Log(max, stop)
//	}
//	if min > stop {
//		t.Error(min, stop)
//	} else {
//		t.Log(min, stop)
//	}
//}
