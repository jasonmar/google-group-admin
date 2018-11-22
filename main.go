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
	"io/ioutil"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/admin/directory/v1"
)

var err error

func main() {
	credentialsFile := "client_id.json"
	tokenFile := "token.json"
	groupsFile := "groups.lst"
	membersFile := "members.conf"
	domain := "phoogle.net"
	deleteExtra := true

	scopes := []string{
		admin.AdminDirectoryGroupScope,
		admin.AdminDirectoryGroupMemberScope,
	}
	b, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		log.Fatalf("Unable to read credentials from '%s': %v", credentialsFile, err)
	}

	config, err := google.ConfigFromJSON(b, scopes...)
	if err != nil {
		log.Fatalf("Unable to parse client secret from '%s': %v", credentialsFile, err)
	}
	client := getClient(config, tokenFile)

	srv, err := admin.New(client)
	if err != nil {
		log.Fatalf("Unable to create Directory client: %v", err)
	}

	ctx := context.Background()

	log.Printf("Reading groups from %s", groupsFile)
	grps, err := readList(groupsFile)
	if err != nil {
		log.Fatalf("Unable to read groups file: %v", err)
	}

	log.Printf("Querying Directory API for groups")
	r, err := listGroups(ctx, srv, domain)
	if err != nil {
		log.Fatalf("Unable to retrieve groups: %v", err)
	}

	want := nameSet(grps)
	exist := groupSet(r)
	need := diff(want, exist)
	alreadyExist := intersect(want, exist)

	if len(need) > 0 {
		log.Printf("Creating groups")
		err = createGroups(ctx, srv, need)
		if err != nil {
			log.Fatalf("Unable to create groups: %v", err)
		}

		printSet(need, "Created:")
	}
	allGroups := union(alreadyExist, need)

	log.Printf("Reading membership from %s", membersFile)
	dir, err := readPairs(membersFile)
	if err != nil {
		log.Fatalf("Failed to list members: %v", err)
	}

	for groupKey, wantMembers := range dir {
		var members map[string]bool
		var err error
		_, exists := allGroups[groupKey]
		if !exists {
			createGroup(ctx, srv, groupKey)
		}
		members, err = listMembers(ctx, srv, groupKey)
		if err != nil {
			log.Printf("Failed to list members for '%s': %v", groupKey, err)
			members = make(map[string]bool)
		}

		extraMembers := diff(members, wantMembers)
		needMembers := diff(wantMembers, members)
		err = addMembers(ctx, srv, groupKey, needMembers)
		if err != nil {
			log.Printf("Failed to add members: %v", err)
		}
		printSet(needMembers, fmt.Sprintf("Members added to %s:", groupKey))
		printSet(extraMembers, fmt.Sprintf("Unexpected Members in %s:", groupKey))
		if deleteExtra {
			deleteMembers(ctx, srv, groupKey, extraMembers)
		}
		log.Printf("Refreshed group '%s'", groupKey)
	}
}
