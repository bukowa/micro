package main

import (
	"context"
	"log"
	"net/http"
)

type (
	Event string
	Task func(e Event, c Client, r Server, payload interface{}) (event Event, p interface{})
)

type Client interface {
	Handle(event Event, payload interface{})
}

type Server interface {
	Client
	Method(event Event, payload interface{})
}

type Micro struct {
	Methods  []map[Event][]Task
	Handlers []map[Event][]Task
	Before   []map[Event]map[Event][]Task
	When     []map[Event]map[Event]map[Event][]Task
}

func (m *Micro) Method(e Event, payload interface{}) {}
func (m *Micro) Handle(e Event, payload interface{}) {}

const AuthRequest Event = "auth"
const LogRequest Event = "mses"
const ServerCreated Event = "server_created"
const ErrorBadPayload Event = ""
const Success Event = ""
const NOOP Event = "noop"
const Stop Event = "stop"
const Start Event = "start"
const Request Event = "message"
const Error Event = "error"
const ServerStarted Event = "server_started"
const ServerStopped Event = "server_stopped"
const Unauthorized Event = "unauthorized"
const Authorized Event = "authorized"
const OK Event = "success"

func main() {
	// CLIENT
	client := &Micro{}

	client.Methods = []map[Event][]Task{
		{Request: {SayHello}},
	}

	client.Handlers = []map[Event][]Task{
		{Error: {}},
		{ErrorBadPayload: {logPrint("bad payload"), logBadPayloadSlack("channel2")}},
		// how should we stop this chain without implementing built-in event type like "Discard"
		// each task can check for a custom event like "Skip"
		{Success: {CheckPayloadType, CheckPayloadType2, CheckPayloadType3}},
		{Unauthorized: {logPrint("unauthorized")}},
	}

	client.When = []map[Event]map[Event]map[Event][]Task{
		{Success: {
			// TaskOK
			OK: {TaskOK, "discard?"},
			Error: {"discard?"},
		}},
	}

	// SERVER
	server := &Micro{}

	server.Methods = []map[Event][]Task{
		{Start: {CreateServer, StartServer, LogStarted,}},
		{Stop: {StopServer}},
	}

	server.Handlers = []map[Event][]Task{
		{Request: {MessageHandler}},
	}

	server.Before = []map[Event]map[Event][]Task{
		{Request: {
			LogRequest:  {logPrint("new request!")},
			AuthRequest: {Authorize},
		}},
	}

	server.When = []map[Event]map[Event]map[Event][]Task{
		{Request: {
			AuthRequest: {
				// ???
				Unauthorized: {ReplyUnathorized?, "discard?"},
			},
		}},
	}

	server.Method(Start, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.Handle(Request, serverRequest{w, r})
	}))

	client.Method(Request, clientRequest{
		server: server,
		body: "hello world!",
	})
}

type serverRequest struct {
	w http.ResponseWriter
	r *http.Request
}

type clientRequest struct{
	server *Micro
	body string
}

var SayHello = Task(func(event Event, caller Client, receiver Server, payload interface{}) (e Event, p interface{}) {
	if v, ok := payload.(clientRequest); ok {
		v.server.Handle(event, v.body)
		return "??", nil
	}
	return ErrorBadPayload, payload
})

var CreateServer = Task(func(event Event, caller Client, receiver Server, payload interface{}) (e Event, p interface{}) {
	return ServerCreated, &http.Server{Handler: payload.(http.HandlerFunc)}
})

var StartServer = Task(func(event Event, caller Client, receiver Server, payload interface{}) (e Event, p interface{}) {
	if v, ok := payload.(*http.Server); ok {
		if err := v.ListenAndServe(); err != nil {
			return Error, err
		}
	}
	return ServerStarted, payload
})

var MessageHandler = Task(func(event Event, c Client, s Server, payload interface{}) (e Event, p interface{}) {
	if _, ok := payload.(clientRequest); ok {
		c.Handle(Success, nil)
		return Success, payload
	}
	c.Handle(ErrorBadPayload, payload)
	return ErrorBadPayload, payload
})

var logPrint = func(v ...interface{}) Task {
	return func(e Event, c Client, r Server, payload interface{}) (event Event, p interface{}) {
		log.Print(e, v)
		return e, p
	}
}

type user struct {
	login string
	password string
}

var Authorize = Task(func(event Event, caller Client, receiver Server, payload interface{}) (e Event, p interface{}) {
	if v, ok := payload.(*user); ok {
		if v.login == "admin" {
			return Authorized, v
		}
	}
	return Unauthorized, payload
})

var StopServer = Task(func(event Event, caller Client, receiver Server, payload interface{}) (e Event, p interface{}) {
	err := p.(*http.Server).Shutdown(context.Background())
	if err != nil {
		return Error, err
	}
	return ServerStopped, payload
})

var LogStarted = Task(func(event Event, caller Client, receiver Server, payload interface{}) (e Event, p interface{}) {
	log.Print("server started")
	return NOOP, payload
})
