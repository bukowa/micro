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
package run

import "fmt"

// Flags represents flags passed to `docker run` command
type Flags []string

// Flag represents opts for Flags
type Flag func(*Flags)

// NewFlags returns Flags modified by Flag
func NewFlags(flags ...Flag) Flags {
	r := Flags{}
	for _, opt := range flags {
		opt(&r)
	}
	return r
}

var Env = func(k, v string) Flag {
	return func(r *Flags) {
		r.Env(k, v)
	}
}

var Volume = func(k, v string) Flag {
	return func(r *Flags) {
		r.Volume(k, v)
	}
}

var Remove Flag = func(r *Flags) {
	r.Remove()
}

var TTY Flag = func(r *Flags) {
	r.TTY()
}

// Env '-e, --env list'
// Set environment variables
func (f *Flags) Env(k, v string) {
	f.add("-e", kv(k, v))
}

// Volume '-v, --volume list'
// Bind mount a volume
func (f *Flags) Volume(k, v string) {
	f.add(fmt.Sprintf("--volume=%s:%s", k, v))
}

// Remove `--rm`
// Automatically remove the container
func (f *Flags) Remove() {
	f.add("--rm")
}

// TTY `--tty` `-t`
// Allocate a pseduo-TTY
func (f *Flags) TTY() {
	f.add("--tty")
}

func (f *Flags) add(s ...string) {
	*f = append(*f, s...)
}

func kv(k, v string) string {
	return fmt.Sprintf("%s=%s", k, v)
}
