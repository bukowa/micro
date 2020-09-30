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
package html_test

import (
	. "github.com/bukowa/micro/html"
	"io"
	"os"
	"path"
	"runtime"
	"testing"
)

func TestCollectSelectorAttributes(t *testing.T) {
	var reader, writer = testDataReader("docker.com.html"), &testWriter{}
	var collector = NewGoQueryCollector(CollectSelectorAttributes(map[string][]string{
		"a":      {"href"},
		"img":    {"src"},
		"script": {"src"},
		"link":   {"href"},
	}))
	if err := collector.Collect(reader, writer); err != nil {
		t.Error(err)
	}
	if len(writer.written) != 103 {
		t.Error("collector not collected all attributes")
	}
}

type testWriter struct {
	written [][]byte
}

func (w *testWriter) Write(b []byte) (n int, err error) {
	w.written = append(w.written, b)
	return 1, nil
}

func testDataReader(name string) io.Reader {
	_, file, _, _ := runtime.Caller(0)
	filedir, _ := path.Split(file)
	filepath := path.Join(filedir, "testdata", name)
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	return f
}
