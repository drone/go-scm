// Copyright (c) 2017 Drone.io Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fixtures

import (
	"io"
	"net/http"
	"net/http/httptest"
)

// NewServer starts a new mock http.Server using the test data.
func NewServer() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(router),
	)
}

func router(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.RawPath) != 0 {
		r.URL.Path = r.URL.RawPath
	}
	if len(r.URL.RawQuery) != 0 {
		r.URL.Path = r.URL.Path + "?" + r.URL.RawQuery
	}
	for _, route := range routes {
		if route.Method == r.Method && route.Path == r.URL.Path {
			for k, v := range route.Header {
				w.Header().Set(k, v)
			}
			w.WriteHeader(route.Status)
			io.WriteString(w, route.Body)
			return
		}
	}
	w.WriteHeader(404)
}

var routes = []struct {
	Method string
	Path   string
	Body   string
	Status int
	Header map[string]string
}{

	// GET /api/v1/repos/gogits/gogs/branches/master
	{
		Method: "GET",
		Path:   "/api/v1/repos/gogits/gogs/branches/master",
		Status: 200,
		Body:   "{\n  \"name\": \"master\",\n  \"commit\": {\n    \"id\": \"f05f642b892d59a0a9ef6a31f6c905a24b5db13a\",\n    \"message\": \"update README\\n\",\n    \"author\": {\n      \"name\": \"Jane Doe\",\n      \"email\": \"jane.doe@mail.com\",\n      \"username\": \"janedoe\"\n    },\n    \"committer\": {\n      \"name\": \"Jane Doe\",\n      \"email\": \"jane.doe@mail.com\",\n      \"username\": \"janedoe\"\n    },\n    \"added\": null,\n    \"removed\": null,\n    \"modified\": null,\n    \"timestamp\": \"2017-11-16T22:06:53Z\"\n  }\n}\n",
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// GET /api/v1/repos/gogits/gogs/branches
	{
		Method: "GET",
		Path:   "/api/v1/repos/gogits/gogs/branches",
		Status: 200,
		Body:   "[\n  {\n    \"name\": \"master\",\n    \"commit\": {\n      \"id\": \"f05f642b892d59a0a9ef6a31f6c905a24b5db13a\",\n      \"message\": \"update README\\n\",\n      \"author\": {\n        \"name\": \"Jane Doe\",\n        \"email\": \"jane.doe@mail.com\",\n        \"username\": \"janedoe\"\n      },\n      \"committer\": {\n        \"name\": \"Jane Doe\",\n        \"email\": \"jane.doe@mail.com\",\n        \"username\": \"janedoe\"\n      },\n      \"added\": null,\n      \"removed\": null,\n      \"modified\": null,\n      \"timestamp\": \"2017-11-16T22:06:53Z\"\n    }\n  }\n]",
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// GET /api/v1/repos/gogits/gogs/raw/f05f642b892d59a0a9ef6a31f6c905a24b5db13a/README.md
	{
		Method: "GET",
		Path:   "/api/v1/repos/gogits/gogs/raw/f05f642b892d59a0a9ef6a31f6c905a24b5db13a/README.md",
		Status: 200,
		Body:   "Hello World!\n",
		Header: map[string]string{
			"Content-Type": "plain/text; charset=UTF-8",
		},
	},

	// DELETE /api/v1/repos/gogits/gogs/hooks/20
	{
		Method: "DELETE",
		Path:   "/api/v1/repos/gogits/gogs/hooks/20",
		Status: 204,
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// GET /api/v1/repos/gogits/gogs/hooks/20
	{
		Method: "GET",
		Path:   "/api/v1/repos/gogits/gogs/hooks/20",
		Status: 200,
		Body:   "{\n  \"id\": 20,\n  \"type\": \"gogs\",\n  \"config\": {\n    \"content_type\": \"json\",\n    \"url\": \"http://gogs.io\"\n  },\n  \"events\": [\n    \"create\",\n    \"push\"\n  ],\n  \"active\": true,\n  \"updated_at\": \"2015-08-29T11:31:22.453572732+08:00\",\n  \"created_at\": \"2015-08-29T11:31:22.453569275+08:00\"\n}\n",
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// GET /api/v1/repos/gogits/gogs/hooks
	{
		Method: "GET",
		Path:   "/api/v1/repos/gogits/gogs/hooks",
		Status: 200,
		Body:   "[\n  {\n    \"id\": 20,\n    \"type\": \"gogs\",\n    \"config\": {\n      \"content_type\": \"json\",\n      \"url\": \"http:\\/\\/gogs.io\"\n    },\n    \"events\": [\n      \"create\",\n      \"push\"\n    ],\n    \"active\": true,\n    \"updated_at\": \"2015-08-29T11:31:22.453572732+08:00\",\n    \"created_at\": \"2015-08-29T11:31:22.453569275+08:00\"\n  }\n]\n",
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// POST /api/v1/repos/gogits/gogs/hooks
	{
		Method: "POST",
		Path:   "/api/v1/repos/gogits/gogs/hooks",
		Status: 201,
		Body:   "{\n  \"id\": 20,\n  \"type\": \"gogs\",\n  \"config\": {\n    \"content_type\": \"json\",\n    \"url\": \"http://gogs.io\"\n  },\n  \"events\": [\n    \"create\",\n    \"push\"\n  ],\n  \"active\": true,\n  \"updated_at\": \"2015-08-29T11:31:22.453572732+08:00\",\n  \"created_at\": \"2015-08-29T11:31:22.453569275+08:00\"\n}\n",
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// DELETE /api/v1/repos/gogits/gogs/issues/1/comments/1
	{
		Method: "DELETE",
		Path:   "/api/v1/repos/gogits/gogs/issues/1/comments/1",
		Status: 204,
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// GET /api/v1/repos/gogits/gogs/issues/1/comments
	{
		Method: "GET",
		Path:   "/api/v1/repos/gogits/gogs/issues/1/comments",
		Status: 200,
		Body:   "[\n  {\n    \"id\": 74,\n    \"user\": {\n      \"id\": 1,\n      \"username\": \"unknwon\",\n      \"full_name\": \"无闻\",\n      \"email\": \"u@gogs.io\",\n      \"avatar_url\": \"http://localhost:3000/avatars/1\"\n    },\n    \"body\": \"what?\",\n    \"created_at\": \"2016-08-26T11:58:18-07:00\",\n    \"updated_at\": \"2016-08-26T11:58:18-07:00\"\n  }\n]",
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// POST /api/v1/repos/gogits/gogs/issues/1/comments
	{
		Method: "POST",
		Path:   "/api/v1/repos/gogits/gogs/issues/1/comments",
		Status: 200,
		Body:   "{\n  \"id\": 74,\n  \"user\": {\n    \"id\": 1,\n    \"username\": \"unknwon\",\n    \"full_name\": \"\\u65e0\\u95fb\",\n    \"email\": \"u@gogs.io\",\n    \"avatar_url\": \"http:\\/\\/localhost:3000\\/avatars\\/1\"\n  },\n  \"body\": \"what?\",\n  \"created_at\": \"2016-08-26T11:58:18-07:00\",\n  \"updated_at\": \"2016-08-26T11:58:18-07:00\"\n}\n",
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// GET /api/v1/repos/gogits/gogs/issues/1
	{
		Method: "GET",
		Path:   "/api/v1/repos/gogits/gogs/issues/1",
		Status: 200,
		Body:   "{\n  \"id\": 1,\n  \"number\": 1,\n  \"user\": {\n    \"id\": 1,\n    \"login\": \"janedoe\",\n    \"full_name\": \"\",\n    \"email\": \"janedoe@mail.com\",\n    \"avatar_url\": \"https:\\/\\/secure.gravatar.com\\/avatar\\/8c58a0be77ee441bb8f8595b7f1b4e87\",\n    \"username\": \"janedoe\"\n  },\n  \"title\": \"Bug found\",\n  \"body\": \"I'm having a problem with this.\",\n  \"labels\": [\n    \n  ],\n  \"milestone\": null,\n  \"assignee\": null,\n  \"state\": \"open\",\n  \"comments\": 0,\n  \"created_at\": \"2017-09-23T19:24:01Z\",\n  \"updated_at\": \"2017-09-23T19:24:01Z\",\n  \"pull_request\": null\n}",
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// GET /api/v1/repos/gogits/gogs/issues
	{
		Method: "GET",
		Path:   "/api/v1/repos/gogits/gogs/issues",
		Status: 200,
		Body:   "[\n  {\n    \"id\": 1,\n    \"number\": 1,\n    \"user\": {\n      \"id\": 1,\n      \"login\": \"janedoe\",\n      \"full_name\": \"\",\n      \"email\": \"janedoe@mail.com\",\n      \"avatar_url\": \"https:\\/\\/secure.gravatar.com\\/avatar\\/8c58a0be77ee441bb8f8595b7f1b4e87\",\n      \"username\": \"janedoe\"\n    },\n    \"title\": \"Bug found\",\n    \"body\": \"I'm having a problem with this.\",\n    \"labels\": [\n      \n    ],\n    \"milestone\": null,\n    \"assignee\": null,\n    \"state\": \"open\",\n    \"comments\": 0,\n    \"created_at\": \"2017-09-23T19:24:01Z\",\n    \"updated_at\": \"2017-09-23T19:24:01Z\",\n    \"pull_request\": null\n  }\n]\n",
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// POST /api/v1/repos/gogits/gogs/issues
	{
		Method: "POST",
		Path:   "/api/v1/repos/gogits/gogs/issues",
		Status: 200,
		Body:   "{\n  \"id\": 1,\n  \"number\": 1,\n  \"user\": {\n    \"id\": 1,\n    \"login\": \"janedoe\",\n    \"full_name\": \"\",\n    \"email\": \"janedoe@mail.com\",\n    \"avatar_url\": \"https:\\/\\/secure.gravatar.com\\/avatar\\/8c58a0be77ee441bb8f8595b7f1b4e87\",\n    \"username\": \"janedoe\"\n  },\n  \"title\": \"Bug found\",\n  \"body\": \"I'm having a problem with this.\",\n  \"labels\": [\n    \n  ],\n  \"milestone\": null,\n  \"assignee\": null,\n  \"state\": \"open\",\n  \"comments\": 0,\n  \"created_at\": \"2017-09-23T19:24:01Z\",\n  \"updated_at\": \"2017-09-23T19:24:01Z\",\n  \"pull_request\": null\n}",
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// GET /api/v1/orgs/gogits
	{
		Method: "GET",
		Path:   "/api/v1/orgs/gogits",
		Status: 200,
		Body:   "{\n  \"id\": 1,\n  \"username\": \"gogits\",\n  \"full_name\": \"gogits\",\n  \"avatar_url\": \"http:\\/\\/gogits.io\\/avatars\\/1\",\n  \"description\": \"\",\n  \"website\": \"\",\n  \"location\": \"\"\n}\n",
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// GET /api/v1/user/orgs
	{
		Method: "GET",
		Path:   "/api/v1/user/orgs",
		Status: 200,
		Body:   "[\n  {\n    \"id\": 1,\n    \"username\": \"gogits\",\n    \"full_name\": \"gogits\",\n    \"avatar_url\": \"http:\\/\\/gogits.io\\/avatars\\/1\",\n    \"description\": \"\",\n    \"website\": \"\",\n    \"location\": \"\"\n  }\n]\n",
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// GET /api/v1/repos/gogits/gogs
	{
		Method: "GET",
		Path:   "/api/v1/repos/gogits/gogs",
		Status: 200,
		Body:   "{\n  \"id\": 1,\n  \"owner\": {\n    \"id\": 1,\n    \"login\": \"gogits\",\n    \"full_name\": \"gogits\",\n    \"email\": \"\",\n    \"avatar_url\": \"http:\\/\\/gogs.io\\/avatars\\/1\",\n    \"username\": \"gogits\"\n  },\n  \"name\": \"gogs\",\n  \"full_name\": \"gogits\\/gogs\",\n  \"description\": \"\",\n  \"private\": true,\n  \"fork\": false,\n  \"parent\": null,\n  \"empty\": false,\n  \"mirror\": false,\n  \"size\": 4485120,\n  \"html_url\": \"http:\\/\\/gogs.io\\/drone\\/cover\",\n  \"ssh_url\": \"git@localhost:drone\\/cover.git\",\n  \"clone_url\": \"http:\\/\\/gogs.io\\/drone\\/cover.git\",\n  \"website\": \"\",\n  \"stars_count\": 0,\n  \"forks_count\": 0,\n  \"watchers_count\": 2,\n  \"open_issues_count\": 0,\n  \"default_branch\": \"master\",\n  \"created_at\": \"2017-10-22T18:25:33Z\",\n  \"updated_at\": \"2017-11-16T22:07:01Z\",\n  \"permissions\": {\n    \"admin\": true,\n    \"push\": true,\n    \"pull\": true\n  }\n}",
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// GET /api/v1/user/repos
	{
		Method: "GET",
		Path:   "/api/v1/user/repos",
		Status: 200,
		Body:   "[\n  {\n    \"id\": 1,\n    \"owner\": {\n      \"id\": 1,\n      \"login\": \"gogits\",\n      \"full_name\": \"gogits\",\n      \"email\": \"\",\n      \"avatar_url\": \"http:\\/\\/gogs.io\\/avatars\\/1\",\n      \"username\": \"gogits\"\n    },\n    \"name\": \"gogs\",\n    \"full_name\": \"gogits\\/gogs\",\n    \"description\": \"\",\n    \"private\": true,\n    \"fork\": false,\n    \"parent\": null,\n    \"empty\": false,\n    \"mirror\": false,\n    \"size\": 4485120,\n    \"html_url\": \"http:\\/\\/gogs.io\\/drone\\/cover\",\n    \"ssh_url\": \"git@localhost:drone\\/cover.git\",\n    \"clone_url\": \"http:\\/\\/gogs.io\\/drone\\/cover.git\",\n    \"website\": \"\",\n    \"stars_count\": 0,\n    \"forks_count\": 0,\n    \"watchers_count\": 2,\n    \"open_issues_count\": 0,\n    \"default_branch\": \"master\",\n    \"created_at\": \"2017-10-22T18:25:33Z\",\n    \"updated_at\": \"2017-11-16T22:07:01Z\",\n    \"permissions\": {\n      \"admin\": true,\n      \"push\": true,\n      \"pull\": true\n    }\n  }\n]\n",
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// GET /api/v1/repos/gogits/go-gogs-client
	{
		Method: "GET",
		Path:   "/api/v1/repos/gogits/go-gogs-client",
		Status: 404,
		Header: map[string]string{},
	},

	// GET /api/v1/users/janedoe
	{
		Method: "GET",
		Path:   "/api/v1/users/janedoe",
		Status: 200,
		Body:   "{\n  \"id\": 1,\n  \"login\": \"janedoe\",\n  \"full_name\": \"\",\n  \"email\": \"janedoe@gmail.com\",\n  \"avatar_url\": \"https:\\/\\/secure.gravatar.com\\/avatar\\/8c58a0be77ee441bb8f8595b7f1b4e87\",\n  \"username\": \"janedoe\"\n}\n",
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},

	// GET /api/v1/user
	{
		Method: "GET",
		Path:   "/api/v1/user",
		Status: 200,
		Body:   "{\n  \"id\": 1,\n  \"login\": \"janedoe\",\n  \"full_name\": \"\",\n  \"email\": \"janedoe@gmail.com\",\n  \"avatar_url\": \"https:\\/\\/secure.gravatar.com\\/avatar\\/8c58a0be77ee441bb8f8595b7f1b4e87\",\n  \"username\": \"janedoe\"\n}",
		Header: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	},
}
