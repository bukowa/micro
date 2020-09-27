package main

import (
	"context"
	"log"
	"net/http"
	"sync"
)

type Client interface {
	Handle(event Event, payload interface{})
}

type Server interface {
	Client
	Method(event Event, payload interface{})
}

type Micro struct {
	eventC chan Event
	x interface{}

	Methods  interface{}
	Handlers interface{}
	Before   []map[Event]map[Event][]Task
	When     []map[Event]map[Event]map[Event][]Task
	Do       interface{}
}

func (m *Micro) Method(e Event, payload interface{}) {

}

func (m *Micro) Handle(e Event, payload interface{}) {

}
type Event string
type Task func(e Event, c Client, r Server, payload interface{}) (event Event, p interface{})

type message struct {
	event Event
	caller *Micro
	payload interface{}
}

type eventListener struct {
	sync.WaitGroup
	task     Task
	receiver *Micro
	requestC chan message
	responseC chan message
	stopC    chan struct{}
}
//
//func (el *eventListener) start() {
//	el.Add(1)
//	go func() {
//		defer el.Done()
//		for {
//			select {
//			case <- el.stopC:
//				return
//			case msg := <- el.requestC:
//				//event := el.task(msg.event, msg.caller, el.receiver)
//				el.responseC <- message{
//					event:  event,
//					caller: el.receiver,
//				}
//			}
//		}
//	}()
//}

const Stopped Event = "stopped"

func (m *Micro) Stop() Task {
	return func(event Event, caller *Micro, receiver *Micro, payload interface{}) (e Event, p interface{}) {
		return Stopped, payload
	}
}

type AuthEvent Event
const AuthRequest Event = "auth"
const AuthLog Event = "log"
const LogRequest Event = "mses"
const Response Event = "response"
const R200 Event = "valid"
const R404 Event = "invalid"

func main() {
	client := &Micro{}
	server := &Micro{}

	client.Methods = []map[Event][]Task{
		{Request: {SayHello}},
	}

	client.Handlers = []map[Event][]Task{
		{Request: {}},
	}

	client.When = []map[Event]map[Event]map[Event][]Task{
		{Request: {
			R200: {},
			R404: {},
		}},
	}
	server.Methods = []map[Event][]Task{
		{Start: {
			CreateServer,
			StartServer,
			LogStarted,
		}},
		{Stop: {
			StopServer,
		}},
	}

	server.Handlers = []map[Event][]Task{
		{Request: {MessageHandler}},
	}

	server.Before = []map[Event]map[Event][]Task{
		{Request: {
			LogRequest:  {logMessage},
			AuthRequest: {Authorize},
		}},
	}

	server.When = []map[Event]map[Event]map[Event][]Task{
		{Request: {
			AuthRequest: {
				Unauthorized: {ResponseUnauthorized, Reply},
			},
		}},
	}

	server.Method(Start, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.Handle(Request, request{w, r})
	}))

	client.Method(Request, clientRequest{
		server: server,
		body: "hello world!",
	})
}


const NOOP Event = "noop"
const Stop Event = "stop"
const Start Event = "start"
const StartError Event = "start_error"
const StopError Event = "stop_error"
const EventHello Event = "hello"
const EventWorld Event = "world"
const Request Event = "message"
const Error Event = "error"
const ServerStarted Event = "server_started"
const ServerStopped Event = "server_stopped"
const Unauthorized Event = "unathorized"
const Authorized Event = "authorized"
const Cancel Event = "cancel"

type clientRequest struct{
	server *Micro
	body string
}

var SayHello = Task(func(event Event, caller *Micro, receiver *Micro, payload interface{}) (e Event, p interface{}) {
	if v, ok := payload.(clientRequest); ok {
		v.server.Handle(event, v.body)
	}
	return NOOP, nil
})

var Reply = Task(func(event Event, caller *Micro, receiver *Micro, payload interface{}) (e Event, p interface{}) {
	caller.Handle(event, payload)
	return NOOP, nil
})

var ResponseUnauthorized = Task(func(event Event, caller *Micro, receiver *Micro, payload interface{}) (e Event, p interface{}) {
	v, ok := payload.(*user); ok {
	}
})


var CreateServer = Task(func(event Event, caller *Micro, receiver *Micro, payload interface{}) (e Event, p interface{}) {
	return NOOP, &http.Server{Handler: payload.(http.HandlerFunc)}
})

type request struct {
	w http.ResponseWriter
	r *http.Request
}

var StartServer = Task(func(event Event, caller *Micro, receiver *Micro, payload interface{}) (e Event, p interface{}) {
	if v, ok := payload.(*http.Server); ok {
		v.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receiver.Handle(Request, request{w, r})
		})
		if err := v.ListenAndServe(); err != nil {
			return Error, err
		}
	}
	return ServerStarted, payload
})

var SendResponse = Task(func(event Event, caller *Micro, receiver *Micro, payload interface{}) (e Event, p interface{}) {
	return NOOP, nil
})

var logMessage = Task(func(event Event, caller *Micro, receiver *Micro, payload interface{}) (e Event, p interface{}) {
	return NOOP, payload
})

var MessageHandler = Task(func(event Event, caller *Micro, receiver *Micro, payload interface{}) (e Event, p interface{}) {
	return NOOP, payload
})

type user struct {
	login string
	password string
}

var Authorize = Task(func(event Event, caller *Micro, receiver *Micro, payload interface{}) (e Event, p interface{}) {
	if v, ok := payload.(*user); ok {
		if v.login == "admin" {
			return Authorized, v
		}
	}
	return Unauthorized, payload
})

var StopServer = Task(func(event Event, caller *Micro, receiver *Micro, payload interface{}) (e Event, p interface{}) {
	err := p.(*http.Server).Shutdown(context.Background())
	if err != nil {
		return Error, err
	}
	return ServerStopped, payload
})

var LogStopped = Task(func(event Event, caller *Micro, receiver *Micro, payload interface{}) (e Event, p interface{}) {
	log.Print("server stopped")
	return NOOP, payload
})

var LogStarted = Task(func(event Event, caller *Micro, receiver *Micro, payload interface{}) (e Event, p interface{}) {
	log.Print("server started")
	return NOOP, payload
})
