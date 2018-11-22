/**
 * @license
 * Copyright Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import "fmt"

func diff(a, b map[string]bool) map[string]bool {
	var exists bool
	r := make(map[string]bool)
	for k := range a {
		_, exists = b[k]
		if !exists {
			r[k] = true
		}
	}
	return r
}

func intersect(a, b map[string]bool) map[string]bool {
	var exists bool
	r := make(map[string]bool)
	for k := range a {
		_, exists = b[k]
		if exists {
			r[k] = true
		}
	}
	return r
}

func union(a, b map[string]bool) map[string]bool {
	r := make(map[string]bool)
	for k := range a {
		r[k] = true
	}
	for k := range b {
		r[k] = true
	}
	return r
}

func printSet(s map[string]bool, title string) {
	if len(s) == 0 {
		return
	}
	fmt.Printf("%s\n", title)
	for k := range s {
		fmt.Printf("%s\n", k)
	}
}
