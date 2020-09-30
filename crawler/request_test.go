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
	"testing"
)

func TestNewRequest(t *testing.T) {
	request, err := NewRequest("GET", "", nil)
	if err != nil {
		t.Error(err)
	}
	if request.Request() == nil {
		t.Error("request http is empty")
	}
}
