# Group Admin Utility

## Purpose

This is a command-line utility which enables membership of Google Groups to be managed via a version controlled configuration file and a command-line utility scheduled to execute as part of a CI/CD pipeline. This provides an audit trail for all group membership changes and a source of record outside Google.

Modifications to the membership file can be approved via code review, with each commit linked to a change management ticket system.


## Usage

Execute the utility with `client_id.json`, `members.conf` and `groups.lst` in the working directory. Groups defined in the configuration files will be created. Members found in the groups but not defined in configuration will be removed.


## Setup

The first time you run the program, you will need to copy an authorization URI from the console and visit the URI in your browser to login and generate and OAuth 2.0 token.

You will see on the authorization screen
```
GoogleGroupAdmin wants to access your Google Account

user@example.com
This will allow GroupAdmin to:
View and manage group subscriptions on your domain
View and manage the provisioning of groups on your domain
```

Copy the string from your browser and paste it into the console then press enter. This will create `token.json` in the working directory for future use. If this token is invalidated you will have to log in again to obtain a new access token.


## Configuration


List of groups that should exist in a file named `groups.lst`

```
group1@example.com
group2@example.com
```

Add group email in brackets followed by a list of members in `members.conf`

```
[group1@example.com]
user1@example.com
user2@example.com

[group2@example.com]
user3@example.com
user4@example.com

[group3@example.com]
```

## Prerequisites

Create client credentials and access token
0. You'll need to enable the [Admin API](https://console.developers.google.com/apis/api/admin.googleapis.com) 
1. Visit [Admin API Credentials](https://console.developers.google.com/apis/api/admin.googleapis.com/credentials)
2. Click "Create Credential"
3. In the "Add credentials to your project" wizard, in "Which API are you using?" leave "Admin SDK" selected and in "Where will you be calling the API from?" select "Other UI (e.g. Windows, CLI, tool)"
4. For "What data will you be accessing?" select "User data"
5. Click "What credentials to I need?"
6. Under "Create an OAuth 2.0 client ID" enter a client name (you can use "GoogleGroupAdmin" or come up with your own name)
7. Click "Create OAuth client ID
8. Click Download (downloads as `client_id.json`)
9. Copy `client_id.json` to the location where you will execute the groupadmin utility

Example token creation:

```
$ groupadmin
Go to the following link in your browser then type the authorization code: 
https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=123-xyz.apps.googleusercontent.com&redirect_uri=urn%3Aietf%3Awg%3Aoauth%3A2.0%3Aoob&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fadmin.directory.group+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fadmin.directory.group.member&state=state-token
<token copied from browser>
Saving credential file to: token.json
...
```


## OAuth 2.0 API Scopes

This application requires a token with two scopes:
`https://www.googleapis.com/auth/admin.directory.group` to create groups
`https://www.googleapis.com/auth/admin.directory.group.member` to add and delete group members

[Admin Directory Oauth 2.0 API Scopes](https://developers.google.com/identity/protocols/googlescopes#admindirectory_v1)


## Documentation

[Admin SDK Group Management](https://developers.google.com/admin-sdk/directory/v1/guides/manage-groups)
[Admin SDK Group Membership](https://developers.google.com/admin-sdk/directory/v1/guides/manage-group-members)


## Example Output


```
$ groupadmin
2018/11/22 02:24:27 Reading groups from groups.lst
2018/11/22 02:24:27 Querying Directory API for groups
2018/11/22 02:24:27 Creating groups
2018/11/22 02:24:28 Created group groupadmin-test-2@example.com
2018/11/22 02:24:29 Created group groupadmin-test-3@example.com
2018/11/22 02:24:29 Created group groupadmin-test-1@example.com
2018/11/22 02:24:29 Reading membership from members.conf
2018/11/22 02:24:31 Added 'user1@example.com' to 'groupadmin-test-1@example.com'
2018/11/22 02:24:32 Added 'user2@example.com' to 'groupadmin-test-1@example.com'
2018/11/22 02:24:32 Refreshed group 'groupadmin-test-1@example.com'
Unexpected Members in groupadmin-test-1@example.com:
user3@example.com
2018/11/22 02:27:51 Deleted 'user3@example.com' from 'groupadmin-test-1@example.com'
2018/11/22 02:24:33 Refreshed group 'groupadmin-test-2@example.com'
2018/11/22 02:24:33 Refreshed group 'groupadmin-test-3@example.com'
```
