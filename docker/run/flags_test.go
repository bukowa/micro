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
package run_test

import (
	. "github.com/bukowa/micro/docker/run"
	"strings"
	"testing"
)

func TestRun_Env(t *testing.T) {
	f := Flags{}
	f.Env("key", "value")
	if !eq(f, "-e key=value") {
		t.Error(f)
	}
}

func TestRun_Volume(t *testing.T) {
	f := Flags{}
	f.Volume("key", "value")
	if !eq(f, "--volume=key:value") {
		t.Error(f)
	}
}

func TestRun_Remove(t *testing.T) {
	f := Flags{}
	f.Remove()
	if !eq(f, "--rm") {
		t.Error(f)
	}
}

func eq(r Flags, s string) bool {
	return strings.Join(r, " ") == s
}

func TestRunEnv(t *testing.T) {
	r := NewFlags(Env("key", "value"))
	if !eq(r, "-e key=value") {
		t.Error(r)
	}
}

func TestRunVolume(t *testing.T) {
	r := NewFlags(Volume("key", "value"))
	if !eq(r, "--volume=key:value") {
		t.Error(r)
	}
}

func TestRunRemove(t *testing.T) {
	r := NewFlags(Remove)
	if !eq(r, "--rm") {
		t.Error(r)
	}
}

func TestRunOpts(t *testing.T) {
	r := NewFlags(
		Env("key", "value"),
		Volume("key", "value"),
		Remove,
	)

	if !eq(r, "-e key=value --volume=key:value --rm") {
		t.Error(r)
	}
}
