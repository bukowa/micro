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
	. "github.com/bukowa/micro/crawler"
	"io/ioutil"
	"testing"
	"time"
)

func TestWaiter(t *testing.T) {
	type args struct {
		count int
		ticks time.Duration
	}
	tests := []struct {
		name      string
		args      args
		wantSleep time.Duration
	}{
		{
			name: "1",
			args: args{
				count: 2,
				ticks: time.Second,
			},
			wantSleep: time.Second * 2,
		},
		{
			name: "2",
			args: args{
				count: 2,
				ticks: time.Second * 5,
			},
			wantSleep: time.Second * 10,
		},
		{
			name: "3",
			args: args{
				count: 5,
				ticks: time.Second,
			},
			wantSleep: time.Second * 5,
		},
		{
			name: "4",
			args: args{
				count: 10,
				ticks: time.Millisecond * 100,
			},
			wantSleep: time.Second,
		},
	}
	for _, tt := range tests {
		// add some time in margin of error
		wantSleep2 := tt.wantSleep - time.Millisecond*300
		wantSleep := tt.wantSleep + time.Millisecond*300

		t.Run(tt.name, func(t *testing.T) {
			crawler := NewCrawler(10,
				WithLoggerOutput(ioutil.Discard),
			)

			crawler.Start()
			start := time.Now()
			WaitUnknownTime(crawler, tt.args.count, tt.args.ticks)

			took := time.Since(start)
			if took > wantSleep {
				t.Errorf("took: %s should: %s", took, "took: %s should: %s")
			}
			if took < wantSleep2 {
				t.Errorf("took: %s should: %s", took, wantSleep2)
			}
		})
	}

}

// Time increase is 300ms because this test can fail on slow machines.
// For example Gitlab-Runner have issues with this.
func TestWaiter2(t *testing.T) {
	var count = 10
	var ticks = time.Second
	var wantLess = time.Second*15 - time.Millisecond*300
	var wantMore = time.Second*15 + time.Millisecond*300
	var ticker = time.NewTicker(time.Second * 5)

	crawler := NewCrawler(10,
		WithLoggerOutput(ioutil.Discard),
	)

	crawler.Start()
	start := time.Now()
	go func() {
		for {
			select {
			case <-ticker.C:
				return
			default:
				req, _ := NewRequest("GET", "http://bad.url.domain", nil)
				crawler.Request() <- req
				<-crawler.Response()
				time.Sleep(time.Millisecond * 100)
			}

		}
	}()
	WaitUnknownTime(crawler, count, ticks)
	took := time.Since(start)

	if took < wantLess {
		t.Errorf("took: %s should: %s", took, wantLess)
	}
	if took > wantMore {
		t.Errorf("took: %s should: %s", took, wantMore)
	}
}
