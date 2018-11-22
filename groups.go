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

func createGroup(ctx context.Context, srv *admin.Service, groupKey string) error {
	g := &admin.Group{
		Email: groupKey,
	}
	call := srv.Groups.Insert(g).Context(ctx)
	_, err := call.Do()
	if err != nil {
		return fmt.Errorf("Failed to create group %s: %v", groupKey, err)
	}
	log.Printf("Created group %s", groupKey)
	return nil
}

func createGroups(ctx context.Context, srv *admin.Service, groups map[string]bool) error {
	var err error
	for groupKey := range groups {
		err = createGroup(ctx, srv, groupKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func listGroups(ctx context.Context, srv *admin.Service, domain string) ([]*admin.Group, error) {
	var g []*admin.Group
	g = make([]*admin.Group, 0, 1000)
	addGroupsFunc := func(resp *admin.Groups) error {
		if len(resp.Groups) > 0 {
			g = append(g, resp.Groups...)
		}
		return nil
	}

	call := srv.Groups.List().Domain(domain).MaxResults(5).Context(ctx)
	call.Pages(ctx, addGroupsFunc)
	return g, nil
}

func mkGroup(srv *admin.Service, group string) error {
	g := &admin.Group{Name: group}
	req := srv.Groups.Insert(g)
	_, err := req.Do()
	return err
}

func mkGroups(srv *admin.Service, groups []string) error {
	for _, group := range groups {
		err := mkGroup(srv, group)
		if err != nil {
			return err
		}
	}
	return nil
}

func groupSet(groups []*admin.Group) map[string]bool {
	s := make(map[string]bool)
	for _, g := range groups {
		s[g.Email] = true
	}
	return s
}

func nameSet(names []string) map[string]bool {
	s := make(map[string]bool)
	for _, name := range names {
		s[name] = true
	}
	return s
}
