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

	// GET /api/v4/projects/diaspora%2Fdiaspora/repository/branches/master
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/repository/branches/master",
		Status: 200,
		Body:   "{\n  \"name\": \"master\",\n  \"merged\": false,\n  \"protected\": true,\n  \"developers_can_push\": false,\n  \"developers_can_merge\": false,\n  \"commit\": {\n    \"author_email\": \"john@example.com\",\n    \"author_name\": \"John Smith\",\n    \"authored_date\": \"2012-06-27T05:51:39-07:00\",\n    \"committed_date\": \"2012-06-28T03:44:20-07:00\",\n    \"committer_email\": \"john@example.com\",\n    \"committer_name\": \"John Smith\",\n    \"id\": \"7b5c3cc8be40ee161ae89a06bba6229da1032a0c\",\n    \"short_id\": \"7b5c3cc\",\n    \"title\": \"add projects API\",\n    \"message\": \"add projects API\",\n    \"parent_ids\": [\n      \"4ad91d3c1144c406e50c7b33bae684bd6837faf8\"\n    ]\n  }\n}\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora/repository/branches?page=1&per_page=30
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/repository/branches?page=1&per_page=30",
		Status: 200,
		Body:   "[\n  {\n    \"name\": \"master\",\n    \"merged\": false,\n    \"protected\": true,\n    \"developers_can_push\": false,\n    \"developers_can_merge\": false,\n    \"commit\": {\n      \"author_email\": \"john@example.com\",\n      \"author_name\": \"John Smith\",\n      \"authored_date\": \"2012-06-27T05:51:39-07:00\",\n      \"committed_date\": \"2012-06-28T03:44:20-07:00\",\n      \"committer_email\": \"john@example.com\",\n      \"committer_name\": \"John Smith\",\n      \"id\": \"7b5c3cc8be40ee161ae89a06bba6229da1032a0c\",\n      \"short_id\": \"7b5c3cc\",\n      \"title\": \"add projects API\",\n      \"message\": \"add projects API\",\n      \"parent_ids\": [\n        \"4ad91d3c1144c406e50c7b33bae684bd6837faf8\"\n      ]\n    }\n  }\n]\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Link":                "<https://api.github.com/resource?page=2>; rel=\"next\", <https://api.github.com/resource?page=1>; rel=\"prev\", <https://api.github.com/resource?page=1>; rel=\"first\", <https://api.github.com/resource?page=5>; rel=\"last\"",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora/repository/commits/6104942438c14ec7bd21c6cd5bd995272b3faff6
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/repository/commits/6104942438c14ec7bd21c6cd5bd995272b3faff6",
		Status: 200,
		Body:   "{\n  \"id\": \"6104942438c14ec7bd21c6cd5bd995272b3faff6\",\n  \"short_id\": \"6104942438c\",\n  \"title\": \"Sanitize for network graph\",\n  \"author_name\": \"randx\",\n  \"author_email\": \"dmitriy.zaporozhets@gmail.com\",\n  \"committer_name\": \"Dmitriy\",\n  \"committer_email\": \"dmitriy.zaporozhets@gmail.com\",\n  \"created_at\": \"2012-06-28T03:44:20-07:00\",\n  \"message\": \"Sanitize for network graph\",\n  \"committed_date\": \"2012-06-28T03:44:20-07:00\",\n  \"authored_date\": \"2012-06-28T03:44:20-07:00\",\n  \"parent_ids\": [\n    \"ae1d9fb46aa2b07ee9836d49862ec4e2c46fbbba\"\n  ],\n  \"last_pipeline\" : {\n    \"id\": 8,\n    \"ref\": \"master\",\n    \"sha\": \"2dc6aa325a317eda67812f05600bdf0fcdc70ab0\",\n    \"status\": \"created\"\n  },\n  \"stats\": {\n    \"additions\": 15,\n    \"deletions\": 10,\n    \"total\": 25\n  },\n  \"status\": \"running\"\n}\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora/repository/commits?page=1&per_page=30&ref_name=master
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/repository/commits?page=1&per_page=30&ref_name=master",
		Status: 200,
		Body:   "[\n  {\n    \"id\": \"6104942438c14ec7bd21c6cd5bd995272b3faff6\",\n    \"short_id\": \"6104942438c\",\n    \"title\": \"Sanitize for network graph\",\n    \"author_name\": \"randx\",\n    \"author_email\": \"dmitriy.zaporozhets@gmail.com\",\n    \"authored_date\": \"2012-06-28T03:44:20-07:00\",\n    \"committer_name\": \"Dmitriy\",\n    \"committer_email\": \"dmitriy.zaporozhets@gmail.com\",\n    \"committed_date\": \"2012-06-28T03:44:20-07:00\",\n    \"created_at\": \"2012-09-20T09:06:12+03:00\",\n    \"message\": \"Sanitize for network graph\",\n    \"parent_ids\": [\n      \"ae1d9fb46aa2b07ee9836d49862ec4e2c46fbbba\"\n    ]\n  }\n]\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Link":                "<https://api.github.com/resource?page=2>; rel=\"next\", <https://api.github.com/resource?page=1>; rel=\"prev\", <https://api.github.com/resource?page=1>; rel=\"first\", <https://api.github.com/resource?page=5>; rel=\"last\"",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora/repository/commits/6104942438c14ec7bd21c6cd5bd995272b3faff6/diff
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/repository/commits/6104942438c14ec7bd21c6cd5bd995272b3faff6/diff",
		Status: 200,
		Body:   "[\n  {\n    \"diff\": \"--- a/doc/update/5.4-to-6.0.md\\n+++ b/doc/update/5.4-to-6.0.md\\n@@ -71,6 +71,8 @@\\n sudo -u git -H bundle exec rake migrate_keys RAILS_ENV=production\\n sudo -u git -H bundle exec rake migrate_inline_notes RAILS_ENV=production\\n \\n+sudo -u git -H bundle exec rake gitlab:assets:compile RAILS_ENV=production\\n+\\n ```\\n \\n ### 6. Update config files\",\n    \"new_path\": \"doc/update/5.4-to-6.0.md\",\n    \"old_path\": \"doc/update/5.4-to-6.0.md\",\n    \"a_mode\": null,\n    \"b_mode\": \"100644\",\n    \"new_file\": true,\n    \"renamed_file\": false,\n    \"deleted_file\": false\n  }\n]\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora/repository/files/app%2Fmodels%2Fkey%2Erb?ref=d5a3ff139356ce33e37e73add446f16869741b50
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/repository/files/app%2Fmodels%2Fkey%2Erb?ref=d5a3ff139356ce33e37e73add446f16869741b50",
		Status: 200,
		Body:   "{\n  \"file_name\": \"key.rb\",\n  \"file_path\": \"app/models/key.rb\",\n  \"size\": 1476,\n  \"encoding\": \"base64\",\n  \"content\": \"cmVxdWlyZSAnZGlnZXN0L21kNScKCmNsYXNzIEtleSA8IEFjdGl2ZVJlY29yZDo6QmFzZQogIGluY2x1ZGUgR2l0bGFiOjpDdXJyZW50U2V0dGluZ3MKICBpbmNsdWRlIFNvcnRhYmxlCgogIGJlbG9uZ3NfdG8gOnVzZXIKCiAgYmVmb3JlX3ZhbGlkYXRpb24gOmdlbmVyYXRlX2ZpbmdlcnByaW50CgogIHZhbGlkYXRlcyA6dGl0bGUsCiAgICBwcmVzZW5jZTogdHJ1ZSwKICAgIGxlbmd0aDogeyBtYXhpbXVtOiAyNTUgfQoKICB2YWxpZGF0ZXMgOmtleSwKICAgIHByZXNlbmNlOiB0cnVlLAogICAgbGVuZ3RoOiB7IG1heGltdW06IDUwMDAgfSwKICAgIGZvcm1hdDogeyB3aXRoOiAvXEEoc3NofGVjZHNhKS0uKlxaLyB9CgogIHZhbGlkYXRlcyA6ZmluZ2VycHJpbnQsCiAgICB1bmlxdWVuZXNzOiB0cnVlLAogICAgcHJlc2VuY2U6IHsgbWVzc2FnZTogJ2Nhbm5vdCBiZSBnZW5lcmF0ZWQnIH0KCiAgdmFsaWRhdGUgOmtleV9tZWV0c19yZXN0cmljdGlvbnMKCiAgZGVsZWdhdGUgOm5hbWUsIDplbWFpbCwgdG86IDp1c2VyLCBwcmVmaXg6IHRydWUKCiAgYWZ0ZXJfY29tbWl0IDphZGRfdG9fc2hlbGwsIG9uOiA6Y3JlYXRlCiAgYWZ0ZXJfY3JlYXRlIDpwb3N0X2NyZWF0ZV9ob29rCiAgYWZ0ZXJfY3JlYXRlIDpyZWZyZXNoX3VzZXJfY2FjaGUKICBhZnRlcl9jb21taXQgOnJlbW92ZV9mcm9tX3NoZWxsLCBvbjogOmRlc3Ryb3kKICBhZnRlcl9kZXN0cm95IDpwb3N0X2Rlc3Ryb3lfaG9vawogIGFmdGVyX2Rlc3Ryb3kgOnJlZnJlc2hfdXNlcl9jYWNoZQoKICBkZWYga2V5PSh2YWx1ZSkKICAgIHZhbHVlJi5kZWxldGUhKCJcblxyIikKICAgIHZhbHVlLnN0cmlwISB1bmxlc3MgdmFsdWUuYmxhbms\",\n  \"ref\": \"master\",\n  \"blob_id\": \"79f7bbd25901e8334750839545a9bd021f0e4c83\",\n  \"commit_id\": \"d5a3ff139356ce33e37e73add446f16869741b50\",\n  \"last_commit_id\": \"570e7b2abdd848b95f2f578043fc23bd6f6fd24d\"\n}\n",
		Header: map[string]string{
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/groups/Twitter
	{
		Method: "GET",
		Path:   "/api/v4/groups/Twitter",
		Status: 200,
		Body:   "{\n  \"id\": 4,\n  \"name\": \"Twitter\",\n  \"path\": \"twitter\",\n  \"description\": \"Aliquid qui quis dignissimos distinctio ut commodi voluptas est.\",\n  \"visibility\": \"public\",\n  \"avatar_url\": \"http://localhost:3000/uploads/group/avatar/1/twitter.jpg\",\n  \"web_url\": \"https://gitlab.example.com/groups/twitter\",\n  \"request_access_enabled\": false,\n  \"full_name\": \"Twitter\",\n  \"full_path\": \"twitter\",\n  \"parent_id\": null,\n  \"projects\": [\n    {\n      \"id\": 7,\n      \"description\": \"Voluptas veniam qui et beatae voluptas doloremque explicabo facilis.\",\n      \"default_branch\": \"master\",\n      \"tag_list\": [],\n      \"archived\": false,\n      \"visibility\": \"public\",\n      \"ssh_url_to_repo\": \"git@gitlab.example.com:twitter/typeahead-js.git\",\n      \"http_url_to_repo\": \"https://gitlab.example.com/twitter/typeahead-js.git\",\n      \"web_url\": \"https://gitlab.example.com/twitter/typeahead-js\",\n      \"name\": \"Typeahead.Js\",\n      \"name_with_namespace\": \"Twitter / Typeahead.Js\",\n      \"path\": \"typeahead-js\",\n      \"path_with_namespace\": \"twitter/typeahead-js\",\n      \"issues_enabled\": true,\n      \"merge_requests_enabled\": true,\n      \"wiki_enabled\": true,\n      \"jobs_enabled\": true,\n      \"snippets_enabled\": false,\n      \"container_registry_enabled\": true,\n      \"created_at\": \"2016-06-17T07:47:25.578Z\",\n      \"last_activity_at\": \"2016-06-17T07:47:25.881Z\",\n      \"shared_runners_enabled\": true,\n      \"creator_id\": 1,\n      \"namespace\": {\n        \"id\": 4,\n        \"name\": \"Twitter\",\n        \"path\": \"twitter\",\n        \"kind\": \"group\"\n      },\n      \"avatar_url\": null,\n      \"star_count\": 0,\n      \"forks_count\": 0,\n      \"open_issues_count\": 3,\n      \"public_jobs\": true,\n      \"shared_with_groups\": [],\n      \"request_access_enabled\": false\n    },\n    {\n      \"id\": 6,\n      \"description\": \"Aspernatur omnis repudiandae qui voluptatibus eaque.\",\n      \"default_branch\": \"master\",\n      \"tag_list\": [],\n      \"archived\": false,\n      \"visibility\": \"internal\",\n      \"ssh_url_to_repo\": \"git@gitlab.example.com:twitter/flight.git\",\n      \"http_url_to_repo\": \"https://gitlab.example.com/twitter/flight.git\",\n      \"web_url\": \"https://gitlab.example.com/twitter/flight\",\n      \"name\": \"Flight\",\n      \"name_with_namespace\": \"Twitter / Flight\",\n      \"path\": \"flight\",\n      \"path_with_namespace\": \"twitter/flight\",\n      \"issues_enabled\": true,\n      \"merge_requests_enabled\": true,\n      \"wiki_enabled\": true,\n      \"jobs_enabled\": true,\n      \"snippets_enabled\": false,\n      \"container_registry_enabled\": true,\n      \"created_at\": \"2016-06-17T07:47:24.661Z\",\n      \"last_activity_at\": \"2016-06-17T07:47:24.838Z\",\n      \"shared_runners_enabled\": true,\n      \"creator_id\": 1,\n      \"namespace\": {\n        \"id\": 4,\n        \"name\": \"Twitter\",\n        \"path\": \"twitter\",\n        \"kind\": \"group\"\n      },\n      \"avatar_url\": null,\n      \"star_count\": 0,\n      \"forks_count\": 0,\n      \"open_issues_count\": 8,\n      \"public_jobs\": true,\n      \"shared_with_groups\": [],\n      \"request_access_enabled\": false\n    }\n  ],\n  \"shared_projects\": [\n    {\n      \"id\": 8,\n      \"description\": \"Velit eveniet provident fugiat saepe eligendi autem.\",\n      \"default_branch\": \"master\",\n      \"tag_list\": [],\n      \"archived\": false,\n      \"visibility\": \"private\",\n      \"ssh_url_to_repo\": \"git@gitlab.example.com:h5bp/html5-boilerplate.git\",\n      \"http_url_to_repo\": \"https://gitlab.example.com/h5bp/html5-boilerplate.git\",\n      \"web_url\": \"https://gitlab.example.com/h5bp/html5-boilerplate\",\n      \"name\": \"Html5 Boilerplate\",\n      \"name_with_namespace\": \"H5bp / Html5 Boilerplate\",\n      \"path\": \"html5-boilerplate\",\n      \"path_with_namespace\": \"h5bp/html5-boilerplate\",\n      \"issues_enabled\": true,\n      \"merge_requests_enabled\": true,\n      \"wiki_enabled\": true,\n      \"jobs_enabled\": true,\n      \"snippets_enabled\": false,\n      \"container_registry_enabled\": true,\n      \"created_at\": \"2016-06-17T07:47:27.089Z\",\n      \"last_activity_at\": \"2016-06-17T07:47:27.310Z\",\n      \"shared_runners_enabled\": true,\n      \"creator_id\": 1,\n      \"namespace\": {\n        \"id\": 5,\n        \"name\": \"H5bp\",\n        \"path\": \"h5bp\",\n        \"kind\": \"group\"\n      },\n      \"avatar_url\": null,\n      \"star_count\": 0,\n      \"forks_count\": 0,\n      \"open_issues_count\": 4,\n      \"public_jobs\": true,\n      \"shared_with_groups\": [\n        {\n          \"group_id\": 4,\n          \"group_name\": \"Twitter\",\n          \"group_access_level\": 30\n        },\n        {\n          \"group_id\": 3,\n          \"group_name\": \"Gitlab Org\",\n          \"group_access_level\": 10\n        }\n      ]\n    }\n  ]\n}\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/groups?page=1&per_page=30
	{
		Method: "GET",
		Path:   "/api/v4/groups?page=1&per_page=30",
		Status: 200,
		Body:   "[\n  {\n    \"id\": 1,\n    \"name\": \"Twitter\",\n    \"path\": \"twitter\",\n    \"description\": \"An interesting group\",\n    \"visibility\": \"public\",\n    \"lfs_enabled\": true,\n    \"avatar_url\": \"http://localhost:3000/uploads/group/avatar/1/twitter.jpg\",\n    \"web_url\": \"http://localhost:3000/groups/twitter\",\n    \"request_access_enabled\": false,\n    \"full_name\": \"Twitter\",\n    \"full_path\": \"twitter\",\n    \"parent_id\": null\n  }\n]",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Link":                "<https://api.github.com/resource?page=2>; rel=\"next\", <https://api.github.com/resource?page=1>; rel=\"prev\", <https://api.github.com/resource?page=1>; rel=\"first\", <https://api.github.com/resource?page=5>; rel=\"last\"",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// POST /api/v4/projects/diaspora%2Fdiaspora/hooks?push_events=true&url=http%3A%2F%2Fexample.com%2Fhook
	{
		Method: "POST",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/hooks?push_events=true&url=http%3A%2F%2Fexample.com%2Fhook",
		Status: 200,
		Body:   "{\n  \"id\": 1,\n  \"url\": \"http://example.com/hook\",\n  \"project_id\": 3,\n  \"push_events\": true,\n  \"issues_events\": true,\n  \"merge_requests_events\": true,\n  \"tag_push_events\": true,\n  \"note_events\": true,\n  \"job_events\": true,\n  \"pipeline_events\": true,\n  \"wiki_page_events\": true,\n  \"enable_ssl_verification\": true,\n  \"created_at\": \"2012-10-12T17:04:47Z\"\n}\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// DELETE /api/v4/projects/diaspora%2Fdiaspora/hooks/1
	{
		Method: "DELETE",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/hooks/1",
		Status: 200,
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora/hooks/1
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/hooks/1",
		Status: 200,
		Body:   "{\n  \"id\": 1,\n  \"url\": \"http://example.com/hook\",\n  \"project_id\": 3,\n  \"push_events\": true,\n  \"issues_events\": true,\n  \"merge_requests_events\": true,\n  \"tag_push_events\": true,\n  \"note_events\": true,\n  \"job_events\": true,\n  \"pipeline_events\": true,\n  \"wiki_page_events\": true,\n  \"enable_ssl_verification\": true,\n  \"created_at\": \"2012-10-12T17:04:47Z\"\n}\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora/hooks?page=1&per_page=30
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/hooks?page=1&per_page=30",
		Status: 200,
		Body:   "[{\n  \"id\": 1,\n  \"url\": \"http://example.com/hook\",\n  \"project_id\": 3,\n  \"push_events\": true,\n  \"issues_events\": true,\n  \"merge_requests_events\": true,\n  \"tag_push_events\": true,\n  \"note_events\": true,\n  \"job_events\": true,\n  \"pipeline_events\": true,\n  \"wiki_page_events\": true,\n  \"enable_ssl_verification\": true,\n  \"created_at\": \"2012-10-12T17:04:47Z\"\n}]\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Link":                "<https://api.github.com/resource?page=2>; rel=\"next\", <https://api.github.com/resource?page=1>; rel=\"prev\", <https://api.github.com/resource?page=1>; rel=\"first\", <https://api.github.com/resource?page=5>; rel=\"last\"",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// POST /api/v4/projects/diaspora%2Fdiaspora/issues?description=I%27m+having+a+problem+with+this.&title=Found+a+bug
	{
		Method: "POST",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/issues?description=I%27m+having+a+problem+with+this.&title=Found+a+bug",
		Status: 200,
		Body:   "{\n   \"project_id\" : 4,\n   \"milestone\" : {\n      \"due_date\" : null,\n      \"project_id\" : 4,\n      \"state\" : \"closed\",\n      \"description\" : \"Rerum est voluptatem provident consequuntur molestias similique ipsum dolor.\",\n      \"iid\" : 3,\n      \"id\" : 11,\n      \"title\" : \"v3.0\",\n      \"created_at\" : \"2016-01-04T15:31:39.788Z\",\n      \"updated_at\" : \"2016-01-04T15:31:39.788Z\",\n      \"closed_at\" : \"2016-01-05T15:31:46.176Z\"\n   },\n   \"author\" : {\n      \"state\" : \"active\",\n      \"web_url\" : \"https://gitlab.example.com/root\",\n      \"avatar_url\" : null,\n      \"username\" : \"root\",\n      \"id\" : 1,\n      \"name\" : \"Administrator\"\n   },\n   \"description\" : \"Omnis vero earum sunt corporis dolor et placeat.\",\n   \"state\" : \"closed\",\n   \"iid\" : 1,\n   \"assignees\" : [{\n      \"avatar_url\" : null,\n      \"web_url\" : \"https://gitlab.example.com/lennie\",\n      \"state\" : \"active\",\n      \"username\" : \"lennie\",\n      \"id\" : 9,\n      \"name\" : \"Dr. Luella Kovacek\"\n   }],\n   \"assignee\" : {\n      \"avatar_url\" : null,\n      \"web_url\" : \"https://gitlab.example.com/lennie\",\n      \"state\" : \"active\",\n      \"username\" : \"lennie\",\n      \"id\" : 9,\n      \"name\" : \"Dr. Luella Kovacek\"\n   },\n   \"labels\" : [],\n   \"id\" : 41,\n   \"title\" : \"Ut commodi ullam eos dolores perferendis nihil sunt.\",\n   \"updated_at\" : \"2016-01-04T15:31:46.176Z\",\n   \"created_at\" : \"2016-01-04T15:31:46.176Z\",\n   \"subscribed\": false,\n   \"user_notes_count\": 1,\n   \"due_date\": null,\n   \"web_url\": \"http://example.com/example/example/issues/1\",\n   \"time_stats\": {\n      \"time_estimate\": 0,\n      \"total_time_spent\": 0,\n      \"human_time_estimate\": null,\n      \"human_total_time_spent\": null\n   },\n   \"confidential\": false,\n   \"discussion_locked\": false,\n   \"_links\": {\n      \"self\": \"http://example.com/api/v4/projects/1/issues/2\",\n      \"notes\": \"http://example.com/api/v4/projects/1/issues/2/notes\",\n      \"award_emoji\": \"http://example.com/api/v4/projects/1/issues/2/award_emoji\",\n      \"project\": \"http://example.com/api/v4/projects/1\"\n   }\n}\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// PUT /api/v4/projects/diaspora%2Fdiaspora/issues/1?state_event=close
	{
		Method: "PUT",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/issues/1?state_event=close",
		Status: 200,
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora/issues/1
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/issues/1",
		Status: 200,
		Body:   "{\n   \"project_id\" : 4,\n   \"milestone\" : {\n      \"due_date\" : null,\n      \"project_id\" : 4,\n      \"state\" : \"closed\",\n      \"description\" : \"Rerum est voluptatem provident consequuntur molestias similique ipsum dolor.\",\n      \"iid\" : 3,\n      \"id\" : 11,\n      \"title\" : \"v3.0\",\n      \"created_at\" : \"2016-01-04T15:31:39.788Z\",\n      \"updated_at\" : \"2016-01-04T15:31:39.788Z\",\n      \"closed_at\" : \"2016-01-05T15:31:46.176Z\"\n   },\n   \"author\" : {\n      \"state\" : \"active\",\n      \"web_url\" : \"https://gitlab.example.com/root\",\n      \"avatar_url\" : null,\n      \"username\" : \"root\",\n      \"id\" : 1,\n      \"name\" : \"Administrator\"\n   },\n   \"description\" : \"Omnis vero earum sunt corporis dolor et placeat.\",\n   \"state\" : \"closed\",\n   \"iid\" : 1,\n   \"assignees\" : [{\n      \"avatar_url\" : null,\n      \"web_url\" : \"https://gitlab.example.com/lennie\",\n      \"state\" : \"active\",\n      \"username\" : \"lennie\",\n      \"id\" : 9,\n      \"name\" : \"Dr. Luella Kovacek\"\n   }],\n   \"assignee\" : {\n      \"avatar_url\" : null,\n      \"web_url\" : \"https://gitlab.example.com/lennie\",\n      \"state\" : \"active\",\n      \"username\" : \"lennie\",\n      \"id\" : 9,\n      \"name\" : \"Dr. Luella Kovacek\"\n   },\n   \"labels\" : [],\n   \"id\" : 41,\n   \"title\" : \"Ut commodi ullam eos dolores perferendis nihil sunt.\",\n   \"updated_at\" : \"2016-01-04T15:31:46.176Z\",\n   \"created_at\" : \"2016-01-04T15:31:46.176Z\",\n   \"subscribed\": false,\n   \"user_notes_count\": 1,\n   \"due_date\": null,\n   \"web_url\": \"http://example.com/example/example/issues/1\",\n   \"time_stats\": {\n      \"time_estimate\": 0,\n      \"total_time_spent\": 0,\n      \"human_time_estimate\": null,\n      \"human_total_time_spent\": null\n   },\n   \"confidential\": false,\n   \"discussion_locked\": false,\n   \"_links\": {\n      \"self\": \"http://example.com/api/v4/projects/1/issues/2\",\n      \"notes\": \"http://example.com/api/v4/projects/1/issues/2/notes\",\n      \"award_emoji\": \"http://example.com/api/v4/projects/1/issues/2/award_emoji\",\n      \"project\": \"http://example.com/api/v4/projects/1\"\n   }\n}\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora/issues?page=1&per_page=30&state=opened
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/issues?page=1&per_page=30&state=opened",
		Status: 200,
		Body:   "[\n   {\n      \"project_id\" : 4,\n      \"milestone\" : {\n         \"due_date\" : null,\n         \"project_id\" : 4,\n         \"state\" : \"closed\",\n         \"description\" : \"Rerum est voluptatem provident consequuntur molestias similique ipsum dolor.\",\n         \"iid\" : 3,\n         \"id\" : 11,\n         \"title\" : \"v3.0\",\n         \"created_at\" : \"2016-01-04T15:31:39.788Z\",\n         \"updated_at\" : \"2016-01-04T15:31:39.788Z\"\n      },\n      \"author\" : {\n         \"state\" : \"active\",\n         \"web_url\" : \"https://gitlab.example.com/root\",\n         \"avatar_url\" : null,\n         \"username\" : \"root\",\n         \"id\" : 1,\n         \"name\" : \"Administrator\"\n      },\n      \"description\" : \"Omnis vero earum sunt corporis dolor et placeat.\",\n      \"state\" : \"closed\",\n      \"iid\" : 1,\n      \"assignees\" : [{\n         \"avatar_url\" : null,\n         \"web_url\" : \"https://gitlab.example.com/lennie\",\n         \"state\" : \"active\",\n         \"username\" : \"lennie\",\n         \"id\" : 9,\n         \"name\" : \"Dr. Luella Kovacek\"\n      }],\n      \"assignee\" : {\n         \"avatar_url\" : null,\n         \"web_url\" : \"https://gitlab.example.com/lennie\",\n         \"state\" : \"active\",\n         \"username\" : \"lennie\",\n         \"id\" : 9,\n         \"name\" : \"Dr. Luella Kovacek\"\n      },\n      \"labels\" : [],\n      \"id\" : 41,\n      \"title\" : \"Ut commodi ullam eos dolores perferendis nihil sunt.\",\n      \"updated_at\" : \"2016-01-04T15:31:46.176Z\",\n      \"created_at\" : \"2016-01-04T15:31:46.176Z\",\n      \"closed_at\" : \"2016-01-05T15:31:46.176Z\",\n      \"user_notes_count\": 1,\n      \"due_date\": \"2016-07-22\",\n      \"web_url\": \"http://example.com/example/example/issues/1\",\n      \"time_stats\": {\n         \"time_estimate\": 0,\n         \"total_time_spent\": 0,\n         \"human_time_estimate\": null,\n         \"human_total_time_spent\": null\n      },\n      \"confidential\": false,\n      \"discussion_locked\": false\n   }\n]\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Link":                "<https://api.github.com/resource?page=2>; rel=\"next\", <https://api.github.com/resource?page=1>; rel=\"prev\", <https://api.github.com/resource?page=1>; rel=\"first\", <https://api.github.com/resource?page=5>; rel=\"last\"",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// PUT /api/v4/projects/diaspora%2Fdiaspora/issues/1?discussion_locked=true
	{
		Method: "PUT",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/issues/1?discussion_locked=true",
		Status: 200,
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// POST /api/v4/projects/diaspora%2Fdiaspora/issues/1/notes?body=what%3F
	{
		Method: "POST",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/issues/1/notes?body=what%3F",
		Status: 200,
		Body:   "{\n  \"id\": 302,\n  \"body\": \"closed\",\n  \"attachment\": null,\n  \"author\": {\n    \"id\": 1,\n    \"username\": \"pipin\",\n    \"email\": \"admin@example.com\",\n    \"name\": \"Pip\",\n    \"state\": \"active\",\n    \"created_at\": \"2013-09-30T13:46:01Z\"\n  },\n  \"created_at\": \"2013-10-02T09:22:45Z\",\n  \"updated_at\": \"2013-10-02T10:22:45Z\",\n  \"system\": true,\n  \"noteable_id\": 377,\n  \"noteable_type\": \"Issue\",\n  \"noteable_iid\": 377\n}\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// DELETE /api/v4/projects/diaspora%2Fdiaspora/issues/1/notes/1
	{
		Method: "DELETE",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/issues/1/notes/1",
		Status: 200,
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora/issues/1/notes/302
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/issues/1/notes/302",
		Status: 200,
		Body:   "{\n  \"id\": 302,\n  \"body\": \"closed\",\n  \"attachment\": null,\n  \"author\": {\n    \"id\": 1,\n    \"username\": \"pipin\",\n    \"email\": \"admin@example.com\",\n    \"name\": \"Pip\",\n    \"state\": \"active\",\n    \"created_at\": \"2013-09-30T13:46:01Z\"\n  },\n  \"created_at\": \"2013-10-02T09:22:45Z\",\n  \"updated_at\": \"2013-10-02T10:22:45Z\",\n  \"system\": true,\n  \"noteable_id\": 377,\n  \"noteable_type\": \"Issue\",\n  \"noteable_iid\": 377\n}\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora/issues/1/notes?page=1&per_page=30
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/issues/1/notes?page=1&per_page=30",
		Status: 200,
		Body:   "[\n  {\n    \"id\": 302,\n    \"body\": \"closed\",\n    \"attachment\": null,\n    \"author\": {\n      \"id\": 1,\n      \"username\": \"pipin\",\n      \"email\": \"admin@example.com\",\n      \"name\": \"Pip\",\n      \"state\": \"active\",\n      \"created_at\": \"2013-09-30T13:46:01Z\"\n    },\n    \"created_at\": \"2013-10-02T09:22:45Z\",\n    \"updated_at\": \"2013-10-02T10:22:45Z\",\n    \"system\": true,\n    \"noteable_id\": 377,\n    \"noteable_type\": \"Issue\",\n    \"noteable_iid\": 377\n  }\n]\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Link":                "<https://api.github.com/resource?page=2>; rel=\"next\", <https://api.github.com/resource?page=1>; rel=\"prev\", <https://api.github.com/resource?page=1>; rel=\"first\", <https://api.github.com/resource?page=5>; rel=\"last\"",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// PUT /api/v4/projects/diaspora%2Fdiaspora/issues/1?discussion_locked=false
	{
		Method: "PUT",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/issues/1?discussion_locked=false",
		Status: 200,
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// PUT /api/v4/projects/gitlab-org%2Ftestme/merge_requests/1?state_event=closed
	{
		Method: "PUT",
		Path:   "/api/v4/projects/gitlab-org%2Ftestme/merge_requests/1?state_event=closed",
		Status: 200,
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/gitlab-org%2Ftestme/merge_requests/1
	{
		Method: "GET",
		Path:   "/api/v4/projects/gitlab-org%2Ftestme/merge_requests/1",
		Status: 200,
		Body:   "{\n  \"id\": 239450,\n  \"iid\": 1,\n  \"project_id\": 32732,\n  \"title\": \"JS fix\",\n  \"description\": \"Signed-off-by: Dmitriy Zaporozhets <dmitriy.zaporozhets@gmail.com>\",\n  \"state\": \"closed\",\n  \"created_at\": \"2015-12-18T18:29:53.563Z\",\n  \"updated_at\": \"2015-12-18T18:30:22.522Z\",\n  \"target_branch\": \"master\",\n  \"source_branch\": \"fix\",\n  \"upvotes\": 0,\n  \"downvotes\": 0,\n  \"author\": {\n    \"id\": 13356,\n    \"name\": \"Drew Blessing\",\n    \"username\": \"dblessing\",\n    \"state\": \"active\",\n    \"avatar_url\": \"https:\\/\\/secure.gravatar.com\\/avatar\\/b5bf44866b4eeafa2d8114bfe15da02f?s=80&d=identicon\",\n    \"web_url\": \"https:\\/\\/gitlab.com\\/dblessing\"\n  },\n  \"assignee\": null,\n  \"source_project_id\": 32732,\n  \"target_project_id\": 32732,\n  \"labels\": [\n    \n  ],\n  \"work_in_progress\": false,\n  \"milestone\": null,\n  \"merge_when_pipeline_succeeds\": false,\n  \"merge_status\": \"can_be_merged\",\n  \"sha\": \"12d65c8dd2b2676fa3ac47d955accc085a37a9c1\",\n  \"merge_commit_sha\": null,\n  \"user_notes_count\": 1,\n  \"approvals_before_merge\": null,\n  \"discussion_locked\": null,\n  \"should_remove_source_branch\": null,\n  \"force_remove_source_branch\": null,\n  \"squash\": false,\n  \"web_url\": \"https:\\/\\/gitlab.com\\/gitlab-org\\/testme\\/merge_requests\\/1\",\n  \"time_stats\": {\n    \"time_estimate\": 0,\n    \"total_time_spent\": 0,\n    \"human_time_estimate\": null,\n    \"human_total_time_spent\": null\n  },\n  \"subscribed\": false,\n  \"changes_count\": null\n}",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/gitlab-org%2Ftestme/merge_requests?page=1&per_page=30&state=opened
	{
		Method: "GET",
		Path:   "/api/v4/projects/gitlab-org%2Ftestme/merge_requests?page=1&per_page=30&state=opened",
		Status: 200,
		Body:   "[\n  {\n    \"id\": 239450,\n    \"iid\": 1,\n    \"project_id\": 32732,\n    \"title\": \"JS fix\",\n    \"description\": \"Signed-off-by: Dmitriy Zaporozhets <dmitriy.zaporozhets@gmail.com>\",\n    \"state\": \"closed\",\n    \"created_at\": \"2015-12-18T18:29:53.563Z\",\n    \"updated_at\": \"2015-12-18T18:30:22.522Z\",\n    \"target_branch\": \"master\",\n    \"source_branch\": \"fix\",\n    \"upvotes\": 0,\n    \"downvotes\": 0,\n    \"author\": {\n      \"id\": 13356,\n      \"name\": \"Drew Blessing\",\n      \"username\": \"dblessing\",\n      \"state\": \"active\",\n      \"avatar_url\": \"https:\\/\\/secure.gravatar.com\\/avatar\\/b5bf44866b4eeafa2d8114bfe15da02f?s=80&d=identicon\",\n      \"web_url\": \"https:\\/\\/gitlab.com\\/dblessing\"\n    },\n    \"assignee\": null,\n    \"source_project_id\": 32732,\n    \"target_project_id\": 32732,\n    \"labels\": [\n      \n    ],\n    \"work_in_progress\": false,\n    \"milestone\": null,\n    \"merge_when_pipeline_succeeds\": false,\n    \"merge_status\": \"can_be_merged\",\n    \"sha\": \"12d65c8dd2b2676fa3ac47d955accc085a37a9c1\",\n    \"merge_commit_sha\": null,\n    \"user_notes_count\": 1,\n    \"approvals_before_merge\": null,\n    \"discussion_locked\": null,\n    \"should_remove_source_branch\": null,\n    \"force_remove_source_branch\": null,\n    \"squash\": false,\n    \"web_url\": \"https:\\/\\/gitlab.com\\/gitlab-org\\/testme\\/merge_requests\\/1\",\n    \"time_stats\": {\n      \"time_estimate\": 0,\n      \"total_time_spent\": 0,\n      \"human_time_estimate\": null,\n      \"human_total_time_spent\": null\n    }\n  }\n]",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Link":                "<https://api.github.com/resource?page=2>; rel=\"next\", <https://api.github.com/resource?page=1>; rel=\"prev\", <https://api.github.com/resource?page=1>; rel=\"first\", <https://api.github.com/resource?page=5>; rel=\"last\"",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// PUT /api/v4/projects/gitlab-org%2Ftestme/merge_requests/1/merge
	{
		Method: "PUT",
		Path:   "/api/v4/projects/gitlab-org%2Ftestme/merge_requests/1/merge",
		Status: 200,
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// POST /api/v4/projects/diaspora%2Fdiaspora/merge_requests/1/notes?body=Comment+for+MR
	{
		Method: "POST",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/merge_requests/1/notes?body=Comment+for+MR",
		Status: 200,
		Body:   "{\n  \"id\": 301,\n  \"body\": \"Comment for MR\",\n  \"attachment\": null,\n  \"author\": {\n    \"id\": 1,\n    \"username\": \"pipin\",\n    \"email\": \"admin@example.com\",\n    \"name\": \"Pip\",\n    \"state\": \"active\",\n    \"created_at\": \"2013-09-30T13:46:01Z\"\n  },\n  \"created_at\": \"2013-10-02T08:57:14Z\",\n  \"updated_at\": \"2013-10-02T08:57:14Z\",\n  \"system\": false,\n  \"noteable_id\": 2,\n  \"noteable_type\": \"MergeRequest\",\n  \"noteable_iid\": 2\n}\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// DELETE /api/v4/projects/diaspora%2Fdiaspora/merge_requests/1/notes/1
	{
		Method: "DELETE",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/merge_requests/1/notes/1",
		Status: 200,
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora/merge_requests/1/notes/301
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/merge_requests/1/notes/301",
		Status: 200,
		Body:   "{\n  \"id\": 301,\n  \"body\": \"Comment for MR\",\n  \"attachment\": null,\n  \"author\": {\n    \"id\": 1,\n    \"username\": \"pipin\",\n    \"email\": \"admin@example.com\",\n    \"name\": \"Pip\",\n    \"state\": \"active\",\n    \"created_at\": \"2013-09-30T13:46:01Z\"\n  },\n  \"created_at\": \"2013-10-02T08:57:14Z\",\n  \"updated_at\": \"2013-10-02T08:57:14Z\",\n  \"system\": false,\n  \"noteable_id\": 2,\n  \"noteable_type\": \"MergeRequest\",\n  \"noteable_iid\": 2\n}\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora/merge_requests/1/notes?page=1&per_page=30
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/merge_requests/1/notes?page=1&per_page=30",
		Status: 200,
		Body:   "[{\n  \"id\": 301,\n  \"body\": \"Comment for MR\",\n  \"attachment\": null,\n  \"author\": {\n    \"id\": 1,\n    \"username\": \"pipin\",\n    \"email\": \"admin@example.com\",\n    \"name\": \"Pip\",\n    \"state\": \"active\",\n    \"created_at\": \"2013-09-30T13:46:01Z\"\n  },\n  \"created_at\": \"2013-10-02T08:57:14Z\",\n  \"updated_at\": \"2013-10-02T08:57:14Z\",\n  \"system\": false,\n  \"noteable_id\": 2,\n  \"noteable_type\": \"MergeRequest\",\n  \"noteable_iid\": 2\n}]",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Link":                "<https://api.github.com/resource?page=2>; rel=\"next\", <https://api.github.com/resource?page=1>; rel=\"prev\", <https://api.github.com/resource?page=1>; rel=\"first\", <https://api.github.com/resource?page=5>; rel=\"last\"",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora",
		Status: 200,
		Body:   "{\n  \"id\": 178504,\n  \"description\": \"\",\n  \"default_branch\": \"master\",\n  \"tag_list\": [\n    \n  ],\n  \"ssh_url_to_repo\": \"git@gitlab.com:diaspora\\/diaspora.git\",\n  \"http_url_to_repo\": \"https:\\/\\/gitlab.com\\/diaspora\\/diaspora.git\",\n  \"web_url\": \"https:\\/\\/gitlab.com\\/diaspora\\/diaspora\",\n  \"name\": \"Diaspora\",\n  \"name_with_namespace\": \"diaspora \\/ Diaspora\",\n  \"path\": \"diaspora\",\n  \"path_with_namespace\": \"diaspora\\/diaspora\",\n  \"avatar_url\": null,\n  \"star_count\": 0,\n  \"forks_count\": 0,\n  \"created_at\": \"2015-03-03T18:37:05.387Z\",\n  \"last_activity_at\": \"2015-03-03T18:37:20.795Z\",\n  \"_links\": {\n    \"self\": \"http:\\/\\/gitlab.com\\/api\\/v4\\/projects\\/178504\",\n    \"issues\": \"http:\\/\\/gitlab.com\\/api\\/v4\\/projects\\/178504\\/issues\",\n    \"merge_requests\": \"http:\\/\\/gitlab.com\\/api\\/v4\\/projects\\/178504\\/merge_requests\",\n    \"repo_branches\": \"http:\\/\\/gitlab.com\\/api\\/v4\\/projects\\/178504\\/repository\\/branches\",\n    \"labels\": \"http:\\/\\/gitlab.com\\/api\\/v4\\/projects\\/178504\\/labels\",\n    \"events\": \"http:\\/\\/gitlab.com\\/api\\/v4\\/projects\\/178504\\/events\",\n    \"members\": \"http:\\/\\/gitlab.com\\/api\\/v4\\/projects\\/178504\\/members\"\n  },\n  \"archived\": false,\n  \"visibility\": \"public\",\n  \"resolve_outdated_diff_discussions\": null,\n  \"container_registry_enabled\": null,\n  \"issues_enabled\": true,\n  \"merge_requests_enabled\": true,\n  \"wiki_enabled\": true,\n  \"jobs_enabled\": true,\n  \"snippets_enabled\": false,\n  \"shared_runners_enabled\": true,\n  \"lfs_enabled\": true,\n  \"creator_id\": 57658,\n  \"namespace\": {\n    \"id\": 120836,\n    \"name\": \"diaspora\",\n    \"path\": \"diaspora\",\n    \"kind\": \"group\",\n    \"full_path\": \"diaspora\",\n    \"parent_id\": null\n  },\n  \"import_status\": \"finished\",\n  \"open_issues_count\": 0,\n  \"public_jobs\": true,\n  \"ci_config_path\": null,\n  \"shared_with_groups\": [\n    \n  ],\n  \"only_allow_merge_if_pipeline_succeeds\": false,\n  \"request_access_enabled\": true,\n  \"only_allow_merge_if_all_discussions_are_resolved\": null,\n  \"printing_merge_request_link_enabled\": true,\n  \"approvals_before_merge\": 0,\n  \"permissions\": {\n    \"project_access\": null,\n    \"group_access\": null\n  }\n}",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/not%2Ffound
	{
		Method: "GET",
		Path:   "/api/v4/projects/not%2Ffound",
		Status: 404,
		Body:   "{\"message\":\"404 Project Not Found\"}\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects?page=1&per_page=30
	{
		Method: "GET",
		Path:   "/api/v4/projects?page=1&per_page=30",
		Status: 200,
		Body:   "[\n  {\n    \"id\": 178504,\n    \"description\": \"\",\n    \"default_branch\": \"master\",\n    \"tag_list\": [\n      \n    ],\n    \"ssh_url_to_repo\": \"git@gitlab.com:diaspora\\/diaspora.git\",\n    \"http_url_to_repo\": \"https:\\/\\/gitlab.com\\/diaspora\\/diaspora.git\",\n    \"web_url\": \"https:\\/\\/gitlab.com\\/diaspora\\/diaspora\",\n    \"name\": \"Diaspora\",\n    \"name_with_namespace\": \"diaspora \\/ Diaspora\",\n    \"path\": \"diaspora\",\n    \"path_with_namespace\": \"diaspora\\/diaspora\",\n    \"avatar_url\": null,\n    \"star_count\": 0,\n    \"forks_count\": 0,\n    \"created_at\": \"2015-03-03T18:37:05.387Z\",\n    \"last_activity_at\": \"2015-03-03T18:37:20.795Z\",\n    \"_links\": {\n      \"self\": \"http:\\/\\/gitlab.com\\/api\\/v4\\/projects\\/178504\",\n      \"issues\": \"http:\\/\\/gitlab.com\\/api\\/v4\\/projects\\/178504\\/issues\",\n      \"merge_requests\": \"http:\\/\\/gitlab.com\\/api\\/v4\\/projects\\/178504\\/merge_requests\",\n      \"repo_branches\": \"http:\\/\\/gitlab.com\\/api\\/v4\\/projects\\/178504\\/repository\\/branches\",\n      \"labels\": \"http:\\/\\/gitlab.com\\/api\\/v4\\/projects\\/178504\\/labels\",\n      \"events\": \"http:\\/\\/gitlab.com\\/api\\/v4\\/projects\\/178504\\/events\",\n      \"members\": \"http:\\/\\/gitlab.com\\/api\\/v4\\/projects\\/178504\\/members\"\n    },\n    \"archived\": false,\n    \"visibility\": \"public\",\n    \"resolve_outdated_diff_discussions\": null,\n    \"container_registry_enabled\": null,\n    \"issues_enabled\": true,\n    \"merge_requests_enabled\": true,\n    \"wiki_enabled\": true,\n    \"jobs_enabled\": true,\n    \"snippets_enabled\": false,\n    \"shared_runners_enabled\": true,\n    \"lfs_enabled\": true,\n    \"creator_id\": 57658,\n    \"namespace\": {\n      \"id\": 120836,\n      \"name\": \"diaspora\",\n      \"path\": \"diaspora\",\n      \"kind\": \"group\",\n      \"full_path\": \"diaspora\",\n      \"parent_id\": null\n    },\n    \"import_status\": \"finished\",\n    \"open_issues_count\": 0,\n    \"public_jobs\": true,\n    \"ci_config_path\": null,\n    \"shared_with_groups\": [\n      \n    ],\n    \"only_allow_merge_if_pipeline_succeeds\": false,\n    \"request_access_enabled\": true,\n    \"only_allow_merge_if_all_discussions_are_resolved\": null,\n    \"printing_merge_request_link_enabled\": true,\n    \"approvals_before_merge\": 0\n  }\n]",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Link":                "<https://api.github.com/resource?page=2>; rel=\"next\", <https://api.github.com/resource?page=1>; rel=\"prev\", <https://api.github.com/resource?page=1>; rel=\"first\", <https://api.github.com/resource?page=5>; rel=\"last\"",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora/repository/commits/18f3e63d05582537db6d183d9d557be09e1f90c8/statuses?page=1&per_page=30
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/repository/commits/18f3e63d05582537db6d183d9d557be09e1f90c8/statuses?page=1&per_page=30",
		Status: 200,
		Body:   "[\n   {\n      \"status\" : \"pending\",\n      \"created_at\" : \"2016-01-19T08:40:25.934Z\",\n      \"started_at\" : null,\n      \"name\" : \"default\",\n      \"allow_failure\" : true,\n      \"author\" : {\n         \"username\" : \"thedude\",\n         \"state\" : \"active\",\n         \"web_url\" : \"https://gitlab.example.com/thedude\",\n         \"avatar_url\" : \"https://gitlab.example.com/uploads/user/avatar/28/The-Big-Lebowski-400-400.png\",\n         \"id\" : 28,\n         \"name\" : \"Jeff Lebowski\"\n      },\n      \"description\" : \"the dude abides\",\n      \"sha\" : \"18f3e63d05582537db6d183d9d557be09e1f90c8\",\n      \"target_url\" : \"https://gitlab.example.com/thedude/gitlab-ce/builds/91\",\n      \"finished_at\" : null,\n      \"id\" : 91,\n      \"ref\" : \"master\"\n   }\n]",
		Header: map[string]string{
			"Link":                "<https://api.github.com/resource?page=2>; rel=\"next\", <https://api.github.com/resource?page=1>; rel=\"prev\", <https://api.github.com/resource?page=1>; rel=\"first\", <https://api.github.com/resource?page=5>; rel=\"last\"",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// POST /api/v4/projects/diaspora%2Fdiaspora/repository/commits/18f3e63d05582537db6d183d9d557be09e1f90c8/statuses?name=continuous-integration%2Fjenkins&state=success&target_url=https%3A%2F%2Fci.example.com%2F1000%2Foutput
	{
		Method: "POST",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/repository/commits/18f3e63d05582537db6d183d9d557be09e1f90c8/statuses?name=continuous-integration%2Fjenkins&state=success&target_url=https%3A%2F%2Fci.example.com%2F1000%2Foutput",
		Status: 200,
		Body:   "{\n   \"author\" : {\n      \"web_url\" : \"https://gitlab.example.com/thedude\",\n      \"name\" : \"Jeff Lebowski\",\n      \"avatar_url\" : \"https://gitlab.example.com/uploads/user/avatar/28/The-Big-Lebowski-400-400.png\",\n      \"username\" : \"thedude\",\n      \"state\" : \"active\",\n      \"id\" : 28\n   },\n   \"name\" : \"default\",\n   \"sha\" : \"18f3e63d05582537db6d183d9d557be09e1f90c8\",\n   \"status\" : \"pending\",\n   \"coverage\": 100.0,\n   \"description\" : \"the dude abides\",\n   \"id\" : 93,\n   \"target_url\" : \"https://gitlab.example.com/thedude/gitlab-ce/builds/91\",\n   \"ref\" : null,\n   \"started_at\" : null,\n   \"created_at\" : \"2016-01-19T09:05:50.355Z\",\n   \"allow_failure\" : false,\n   \"finished_at\" : \"2016-01-19T09:05:50.365Z\"\n}",
		Header: map[string]string{
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora/repository/tags/v1.0.0
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/repository/tags/v1.0.0",
		Status: 200,
		Body:   "{\n  \"name\": \"v1.0.0\",\n  \"message\": null,\n  \"commit\": {\n    \"id\": \"2695effb5807a22ff3d138d593fd856244e155e7\",\n    \"short_id\": \"2695effb\",\n    \"title\": \"Initial commit\",\n    \"created_at\": \"2017-07-26T11:08:53.000+02:00\",\n    \"parent_ids\": [\n      \"2a4b78934375d7f53875269ffd4f45fd83a84ebe\"\n    ],\n    \"message\": \"v1.0.0\\n\",\n    \"author_name\": \"Arthur Verschaeve\",\n    \"author_email\": \"contact@arthurverschaeve.be\",\n    \"authored_date\": \"2015-02-01T21:56:31.000+01:00\",\n    \"committer_name\": \"Arthur Verschaeve\",\n    \"committer_email\": \"contact@arthurverschaeve.be\",\n    \"committed_date\": \"2015-02-01T21:56:31.000+01:00\"\n  },\n  \"release\": null\n}\n",
		Header: map[string]string{
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/projects/diaspora%2Fdiaspora/repository/tags?page=1&per_page=30
	{
		Method: "GET",
		Path:   "/api/v4/projects/diaspora%2Fdiaspora/repository/tags?page=1&per_page=30",
		Status: 200,
		Body:   "[\n  {\n    \"commit\": {\n      \"id\": \"2695effb5807a22ff3d138d593fd856244e155e7\",\n      \"short_id\": \"2695effb\",\n      \"title\": \"Initial commit\",\n      \"created_at\": \"2017-07-26T11:08:53.000+02:00\",\n      \"parent_ids\": [\n        \"2a4b78934375d7f53875269ffd4f45fd83a84ebe\"\n      ],\n      \"message\": \"Initial commit\",\n      \"author_name\": \"John Smith\",\n      \"author_email\": \"john@example.com\",\n      \"authored_date\": \"2015-02-01T21:56:31.000+01:00\",\n      \"committer_name\": \"Jack Smith\",\n      \"committer_email\": \"jack@example.com\",\n      \"committed_date\": \"2015-02-01T21:56:31.000+01:00\"\n    },\n    \"release\": {\n      \"tag_name\": \"1.0.0\",\n      \"description\": \"Amazing release. Wow\"\n    },\n    \"name\": \"v1.0.0\",\n    \"message\": null\n  }\n]\n",
		Header: map[string]string{
			"Link":                "<https://api.github.com/resource?page=2>; rel=\"next\", <https://api.github.com/resource?page=1>; rel=\"prev\", <https://api.github.com/resource?page=1>; rel=\"first\", <https://api.github.com/resource?page=5>; rel=\"last\"",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/user
	{
		Method: "GET",
		Path:   "/api/v4/user",
		Status: 200,
		Body:   "{\n  \"id\": 1,\n  \"username\": \"john_smith\",\n  \"email\": \"john@example.com\",\n  \"name\": \"John Smith\",\n  \"state\": \"active\",\n  \"avatar_url\": \"http://localhost:3000/uploads/user/avatar/1/index.jpg\",\n  \"web_url\": \"http://localhost:3000/john_smith\",\n  \"created_at\": \"2012-05-23T08:00:58Z\",\n  \"bio\": null,\n  \"location\": null,\n  \"skype\": \"\",\n  \"linkedin\": \"\",\n  \"twitter\": \"\",\n  \"website_url\": \"\",\n  \"organization\": \"\",\n  \"last_sign_in_at\": \"2012-06-01T11:41:01Z\",\n  \"confirmed_at\": \"2012-05-23T09:05:22Z\",\n  \"theme_id\": 1,\n  \"last_activity_on\": \"2012-05-23\",\n  \"color_scheme_id\": 2,\n  \"projects_limit\": 100,\n  \"current_sign_in_at\": \"2012-06-02T06:36:55Z\",\n  \"identities\": [\n    {\"provider\": \"github\", \"extern_uid\": \"2435223452345\"},\n    {\"provider\": \"bitbucket\", \"extern_uid\": \"john_smith\"},\n    {\"provider\": \"google_oauth2\", \"extern_uid\": \"8776128412476123468721346\"}\n  ],\n  \"can_create_group\": true,\n  \"can_create_project\": true,\n  \"two_factor_enabled\": true,\n  \"external\": false\n}\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/users?search=john_smith
	{
		Method: "GET",
		Path:   "/api/v4/users?search=john_smith",
		Status: 200,
		Body:   "[{\n  \"id\": 1,\n  \"username\": \"john_smith\",\n  \"email\": \"john@example.com\",\n  \"name\": \"John Smith\",\n  \"state\": \"active\",\n  \"avatar_url\": \"http://localhost:3000/uploads/user/avatar/1/index.jpg\",\n  \"web_url\": \"http://localhost:3000/john_smith\",\n  \"created_at\": \"2012-05-23T08:00:58Z\",\n  \"bio\": null,\n  \"location\": null,\n  \"skype\": \"\",\n  \"linkedin\": \"\",\n  \"twitter\": \"\",\n  \"website_url\": \"\",\n  \"organization\": \"\",\n  \"last_sign_in_at\": \"2012-06-01T11:41:01Z\",\n  \"confirmed_at\": \"2012-05-23T09:05:22Z\",\n  \"theme_id\": 1,\n  \"last_activity_on\": \"2012-05-23\",\n  \"color_scheme_id\": 2,\n  \"projects_limit\": 100,\n  \"current_sign_in_at\": \"2012-06-02T06:36:55Z\",\n  \"identities\": [\n    {\"provider\": \"github\", \"extern_uid\": \"2435223452345\"},\n    {\"provider\": \"bitbucket\", \"extern_uid\": \"john_smith\"},\n    {\"provider\": \"google_oauth2\", \"extern_uid\": \"8776128412476123468721346\"}\n  ],\n  \"can_create_group\": true,\n  \"can_create_project\": true,\n  \"two_factor_enabled\": true,\n  \"external\": false\n}]\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},

	// GET /api/v4/users?search=nobody
	{
		Method: "GET",
		Path:   "/api/v4/users?search=nobody",
		Status: 401,
		Body:   "{\"message\":\"401 Unauthorized\"}\n",
		Header: map[string]string{
			"Content-Type":        "application/json",
			"Ratelimit-Limit":     "600",
			"Ratelimit-Observed":  "1",
			"Ratelimit-Remaining": "599",
			"Ratelimit-Reset":     "1512454441",
			"Ratelimit-Resettime": "Wed, 05 Dec 2017 06:14:01 GMT",
			"X-Request-Id":        "0d511a76-2ade-4c34-af0d-d17e84adb255",
		},
	},
}
