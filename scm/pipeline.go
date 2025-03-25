package scm

import "time"

type (
	Pipeline struct {
		ID          string     `json:"id"`
		Status      string     `json:"status"`
		Branch      string     `json:"branch"`
		CommitSHA   string     `json:"commit_sha"`
		Author      string     `json:"author"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   *time.Time `json:"updated_at,omitempty"`
		PipelineURL string     `json:"pipeline_url"`
		RepoName    string     `json:"repo_name"`
	}
)
