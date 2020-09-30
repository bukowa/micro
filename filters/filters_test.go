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
package filters_test

import (
	. "github.com/bukowa/micro/filters"
	"testing"
)

type testCase struct {
	filter Filter
	tests  []test
}

type test struct {
	value    string
	filtered bool
}

func runTestCase(t *testing.T, testCase *testCase) {
	if len(testCase.tests) == 0 {
		panic("")
	}
	for _, test := range testCase.tests {
		t.Run(test.value, func(t *testing.T) {
			if got := testCase.filter.Func(test.value); got != test.filtered {
				t.Errorf("%s: filtered: %v got: %v", test.value, test.filtered, got)
			}
		})
	}
}

func TestPrefix(t *testing.T) {
	test := &testCase{
		filter: Prefix("a", "b"),
		tests: []test{
			{
				value:    "abc",
				filtered: true,
			},
			{
				value:    "bca",
				filtered: true,
			},
			{
				value:    "cab",
				filtered: false,
			},
		},
	}
	runTestCase(t, test)
}

func TestSuffix(t *testing.T) {
	test := &testCase{
		filter: Suffix("a", "b"),
		tests: []test{
			{
				value:    "abca",
				filtered: true,
			},
			{
				value:    "abcab",
				filtered: true,
			},
			{
				value:    "abcabc",
				filtered: false,
			},
		}}
	runTestCase(t, test)
}

func TestEqual(t *testing.T) {
	test := &testCase{
		filter: Equal("a", "b"),
		tests: []test{
			{
				value:    "a",
				filtered: true,
			},
			{
				value:    "b",
				filtered: true,
			},
			{
				value:    "c",
				filtered: false,
			},
		}}
	runTestCase(t, test)
}

func TestEqualFold(t *testing.T) {
	test := &testCase{
		filter: EqualFold("a", "b"),
		tests: []test{
			{
				value:    "a",
				filtered: true,
			},
			{
				value:    "b",
				filtered: true,
			},
			{
				value:    "c",
				filtered: false,
			},
			{
				value:    "A",
				filtered: true,
			},
		}}
	runTestCase(t, test)
}

func TestNotEqual(t *testing.T) {
	runTestCase(t, &testCase{
		filter: NotEqual("a", "b"),
		tests: []test{
			{
				value:    "a",
				filtered: false,
			},
			{
				value:    "b",
				filtered: false,
			},
			{
				value:    "A",
				filtered: true,
			},
			{
				value:    "B",
				filtered: true,
			},
			{
				value:    "ab",
				filtered: true,
			},
			{
				value:    "ba",
				filtered: true,
			},
		},
	})
}
