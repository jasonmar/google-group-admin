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

import (
	"bufio"
	"os"
	"strings"
)

func readList(path string) ([]string, error) {
	out := make([]string, 0, 1024)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			out = append(out, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func normalize(s string) string {
	var out string
	out = s
	out = strings.Replace(out, " ", "", -1)
	out = strings.Replace(out, "\t", "", -1)
	out = strings.ToLower(out)
	return out
}

func isValidMember(s string) bool {
	if strings.HasPrefix(s, "#") {
		return false
	}
	if len(strings.Split(s, "@")) != 2 {
		return false
	}
	return true
}

func readPairs(path string) (map[string]map[string]bool, error) {
	out := make(map[string]map[string]bool)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var group, member string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := normalize(scanner.Text())
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			group = line
			group = strings.TrimPrefix(group, "[")
			group = strings.TrimSuffix(group, "]")
			_, exists := out[group]
			if !exists {
				out[group] = make(map[string]bool)
			}
		} else if isValidMember(line) {
			member = line
			m, exists := out[group]
			if !exists {
				m = make(map[string]bool)
				out[group] = m
			}
			m[member] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
