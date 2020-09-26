# micro

```go
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
```