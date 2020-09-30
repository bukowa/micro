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
package filters

// Name is the name of the filter.
type Name string

// Func - given value returns true - if the value is filtered.
type Func = func(value string) bool

// Filter is responsible for filtering values.
type Filter interface {
	Name() Name
	Func(value string) bool
}

// New creates new basic Filter implementation.
func New(name Name, f func(value string) bool) Filter {
	return &filter{
		name:  name,
		xfunc: f,
	}
}

// filter implements basic Filter.
type filter struct {
	name  Name
	xfunc Func
}

// Name returns FilterName for given filter.
func (f *filter) Name() Name {
	return f.name
}

// Func is a FilterFunc.
func (f *filter) Func(uri string) bool {
	return f.xfunc(uri)
}
