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
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/api/admin/directory/v1"
)

func listMembers(ctx context.Context, srv *admin.Service, group string) (map[string]bool, error) {
	out := make(map[string]bool)
	call := srv.Members.List(group).Context(ctx)
	err := call.Pages(ctx, func(resp *admin.Members) error {
		if len(resp.Members) > 0 {
			for _, m := range resp.Members {
				out[m.Email] = true
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func addMembers(ctx context.Context, srv *admin.Service, group string, members map[string]bool) error {
	for email := range members {
		call := srv.Members.Insert(group, &admin.Member{Email: email}).Context(ctx)
		_, err := call.Do()
		if err != nil {
			return fmt.Errorf("Failed to add '%s' to '%s': %v", email, group, err)
		}
		log.Printf("Added '%s' to '%s'", email, group)
	}
	return nil
}

func deleteMembers(ctx context.Context, srv *admin.Service, group string, members map[string]bool) error {
	for email := range members {
		call := srv.Members.Delete(group, email).Context(ctx)
		err := call.Do()
		if err != nil {
			return fmt.Errorf("Failed to delete '%s' from '%s': %v", email, group, err)
		}
		log.Printf("Deleted '%s' from '%s'", email, group)
	}
	return nil
}
