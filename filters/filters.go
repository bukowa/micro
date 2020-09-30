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

import (
	"strings"
)

const (
	PrefixName    Name = "prefix"
	SuffixName    Name = "suffix"
	EqualName     Name = "equal"
	EqualFoldName Name = "equal_fold"
	NotEqualName  Name = "not_equal"
	ContainsName  Name = "contains"
)

var (
	// Prefix filters uri that starts with any of prefix...
	Prefix = func(prefix ...string) Filter {
		return New(PrefixName, func(value string) bool {
			for _, pre := range prefix {
				if strings.HasPrefix(value, pre) {
					return true
				}
			}
			return false
		})
	}
	// Suffix filters value that ends with any of suffix...
	Suffix = func(suffix ...string) Filter {
		return New(SuffixName, func(value string) bool {
			for _, suf := range suffix {
				if strings.HasSuffix(value, suf) {
					return true
				}
			}
			return false
		})
	}
	// Equal filters value that matches any of match...
	Equal = func(match ...string) Filter {
		return New(EqualName, func(value string) bool {
			for _, m := range match {
				if value == m {
					return true
				}
			}
			return false
		})
	}
	// EqualFold filters value that equal fold for any of match...
	EqualFold = func(match ...string) Filter {
		return New(EqualFoldName, func(value string) bool {
			for _, m := range match {
				if strings.EqualFold(value, m) {
					return true
				}
			}
			return false
		})
	}

	// NotEqual filters value that is not equal to any of match.
	NotEqual = func(match ...string) Filter {
		return New(NotEqualName, func(value string) bool {
			for _, m := range match {
				if m == value {
					return false
				}
			}
			return true
		})
	}

	// Contains filters value that contains any of match...
	Contains = func(match ...string) Filter {
		return New(ContainsName, func(value string) bool {
			for _, m := range match {
				if strings.Contains(value, m) {
					return true
				}
			}
			return false
		})
	}
)
