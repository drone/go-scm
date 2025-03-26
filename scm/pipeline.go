package scm

import "time"

type (
	Pipeline struct {
		ID       string     `json:"id"`
		Status   string     `json:"status"`
		Created  time.Time  `json:"created_at"`
		Updated  *time.Time `json:"updated_at,omitempty"`
		URL      string     `json:"pipeline_url"`
		RepoName string     `json:"repo_name"`
	}
)
