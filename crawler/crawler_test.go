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
package crawler_test

import (
	"bufio"
	"bytes"
	. "github.com/bukowa/micro/crawler"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCrawler(t *testing.T) {
	var requestCounter = NewCounter()
	var responseCounter = NewCounter()
	var serverCounter = NewCounter()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCounter.Add(1)
		w.WriteHeader(200)
	}))
	defer ts.Close()

	opts := []Option{

	}
	c := NewCrawler(5, opts...)

	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				return
			default:
				r, err := NewRequest("GET", ts.URL, nil)
				if err != nil {
					t.Error(err)
				}
				c.Request() <- r
				requestCounter.Add(1)
			}
		}
	}()

	var e = make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-c.Response():
				responseCounter.Add(1)
			case <-e:
				return
			}
		}
	}()

	c.Start()
	WaitUnknownTime(c, 1, time.Second)
	e <- struct{}{}

	var srvN = serverCounter.Size()
	var reqN = requestCounter.Size()
	var resN = responseCounter.Size()

	if srvN == 0 || reqN == 0 || resN == 0 {
		t.Error()
	}
	if srvN != reqN {
		t.Error()
	}
	if srvN != resN {
		t.Error()
	}
}


func TestCrawlerWithDefaultLog(t *testing.T) {
	var output = bytes.NewBuffer(nil)
	var expected = map[int]string{
		0: "start[start]",
		1: "started[started]",
		2: "0:request:invalid",
		3: "0:response:invalid:err:Get \"invalid\": unsupported protocol scheme \"\"",
		4: "stop[stop]",
		5: "stopped[stopped]",
		6: "wait[wait]",
	}

	var crawler = NewCrawler(1, WithDefaultLog, WithLoggerOutput(output))
	req, _ := NewRequest("GET", "invalid", nil)

	crawler.Start()
	crawler.Request() <- req
	<- crawler.Response()
	crawler.Stop()
	crawler.Wait()

	logs := gatherLines(output)
	if len(logs) != len(expected) {
		t.Errorf("want len: %v, got len: %v", len(expected), len(logs))
	}
	for i, log := range logs {
		if !strings.Contains(log, expected[i]) {
			t.Errorf("could not find: %s in log: %s", expected[i], log)
		}
	}
}


func gatherLines(r io.Reader) []string {
	var s []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s = append(s, scanner.Text())
	}
	return s
}