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
	"net/http"
	"testing"
	"time"

	. "github.com/bukowa/micro/crawler"
)

func TestNewResponse(t *testing.T) {
	request, err := NewRequest("GET", "", nil)
	if err != nil {
		t.Error(err)
	}
	response := NewResponse(nil, time.Second, request, &http.Response{Request: request.Request()}, nil)
	if response.Request().Method != "GET" {
		t.Error("request has wrong method")
	}
	if response.Response().Request.Method != "GET" {
		t.Error("response request has wrong method")
	}
	if response.Time() != time.Second {
		t.Errorf("response time invalid: %s", response.Time())
	}
}
