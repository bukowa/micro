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
package crawler

import (
	"io"
	"log"
	"os"
)

// Logger is a logger used by Crawler.
type Logger interface {
	Print(v ...interface{})
	Fatal(v ...interface{})
	Printf(format string, v ...interface{})
	SetPrefix(prefix string)
	SetOutput(writer io.Writer)
}

// NewLogger creates new Logger.
func NewLogger() Logger {
	prefix := "crawler:"
	return &BaseLogger{Logger: log.New(os.Stderr, prefix, log.LstdFlags)}
}

// BaseLogger implements Logger.
type BaseLogger struct {
	*log.Logger
}

// PrintRequest logs Request.
func (l *BaseLogger) PrintRequest(r Request, i int) {
	l.Printf("%v:request:url:%s", i, r.Request().URL.String())
}

// PrintResponse logs Response.
func (l *BaseLogger) PrintResponse(r Response, i int) {
	l.Printf("%v:response:url:%s:err:%s", i, r.Request().URL.String(), r.Error())
}
