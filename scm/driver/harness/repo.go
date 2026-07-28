// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/drone/go-scm/scm"
)

type repositoryService struct {
	client *wrapper
}

func (s *repositoryService) Find(ctx context.Context, repo string) (*scm.Repository, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(harnessURI)
	if err != nil {
		return nil, nil, err
	}
	path := fmt.Sprintf("api/v1/repos/%s?%s", repoId, queryParams)
	out := new(repository)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	if err != nil {
		return nil, res, err
	}
	convertedRepo := convertRepository(out)
	if convertedRepo == nil {
		return nil, res, errors.New("Harness returned an unexpected null repository")
	}
	return convertedRepo, res, err
}

func (s *repositoryService) FindHook(ctx context.Context, repo string, id string) (*scm.Hook, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(harnessURI)
	if err != nil {
		return nil, nil, err
	}
	path := fmt.Sprintf("api/v1/repos/%s/webhooks/%s?%s", repoId, id, queryParams)
	out := new(hook)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertHook(out), res, err
}

func (s *repositoryService) FindPerms(ctx context.Context, repo string) (*scm.Perm, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	return s.list(ctx, scm.RepoListOptions{
		ListOptions: opts,
	})
}

func (s *repositoryService) list(ctx context.Context, opts scm.RepoListOptions) ([]*scm.Repository, *scm.Response, error) {
	queryParams := fmt.Sprintf("%s=%s&%s=%s&%s=%s&%s=%s",
		projectIdentifier, s.client.project, orgIdentifier, s.client.organization, accountIdentifier, s.client.account,
		routingId, s.client.account)

	if opts.RepoSearchTerm.RepoName != "" {
		queryParams = fmt.Sprintf("%s&query=%s", queryParams, opts.RepoSearchTerm.RepoName)
	}

	if opts.ListOptions.SortKey == "" {
		opts.ListOptions.SortKey = defaultSortKey
	}

	if opts.ListOptions.Order == "" {
		opts.ListOptions.Order = defaultOrder
	}

	path := fmt.Sprintf("api/v1/repos?%s&%s", encodeListOptions(opts.ListOptions), queryParams)
	out := []*repository{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertRepositoryList(out), res, err
}

func (s *repositoryService) ListV2(ctx context.Context, opts scm.RepoListOptions) ([]*scm.Repository, *scm.Response, error) {
	return s.list(ctx, opts)
}

func (s *repositoryService) ListNamespace(ctx context.Context, _ string, opts scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	// Client already has context about namespace
	return s.List(ctx, opts)
}

func (s *repositoryService) ListHooks(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Hook, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(harnessURI)
	if err != nil {
		return nil, nil, err
	}
	if opts.SortKey == "" {
		opts.SortKey = "display_name"
	}
	if opts.Order == "" {
		opts.Order = "asc"
	}
	path := fmt.Sprintf("api/v1/repos/%s/webhooks?%s&%s", repoId, encodeListOptions(opts), queryParams)
	out := []*hook{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertHookList(out), res, err
}

func (s *repositoryService) ListStatus(ctx context.Context, repo string, ref string, opts scm.ListOptions) ([]*scm.Status, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) CreateHook(ctx context.Context, repo string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(harnessURI)
	if err != nil {
		return nil, nil, err
	}
	path := fmt.Sprintf("api/v1/repos/%s/webhooks?%s", repoId, queryParams)
	in := new(hook)
	in.Enabled = true
	in.Identifier = input.Name
	in.Secret = input.Secret
	in.Insecure = input.SkipVerify
	in.URL = input.Target
	in.Triggers = input.NativeEvents
	out := new(hook)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertHook(out), res, err
}

func (s *repositoryService) CreateStatus(ctx context.Context, repo string, ref string, input *scm.StatusInput) (*scm.Status, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(harnessURI)
	if err != nil {
		return nil, nil, err
	}
	path := fmt.Sprintf("api/v1/repos/%s/checks/commits/%s?%s", repoId, ref, queryParams)
	summary := input.Desc
	if summary == "" {
		summary = input.Title
	}
	in := &checkReportInput{
		Identifier: input.Label,
		Status:     convertToHarnessCheckStatus(input.State),
		Summary:    summary,
		Link:       input.Target,
		Payload: checkPayload{
			Kind: "raw",
		},
	}
	out := new(check)
	res, err := s.client.do(ctx, "PUT", path, in, out)
	return convertCheck(out), res, err
}

func (s *repositoryService) UpdateHook(ctx context.Context, repo, id string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) DeleteHook(ctx context.Context, repo string, id string) (*scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	repoId, queryParams, err := getRepoAndQueryParams(harnessURI)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("api/v1/repos/%s/webhooks/%s?%s", repoId, id, queryParams)
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

//
// native data structures
//

type (
	// harness repository resource.
	repository struct {
		ID             int    `json:"id"`
		ParentID       int    `json:"parent_id"`
		UID            string `json:"uid"`
		Path           string `json:"path"`
		Description    string `json:"description"`
		IsPublic       bool   `json:"is_public"`
		CreatedBy      int    `json:"created_by"`
		Created        int64  `json:"created"`
		Updated        int64  `json:"updated"`
		DefaultBranch  string `json:"default_branch"`
		ForkID         int    `json:"fork_id"`
		NumForks       int    `json:"num_forks"`
		NumPulls       int    `json:"num_pulls"`
		NumClosedPulls int    `json:"num_closed_pulls"`
		NumOpenPulls   int    `json:"num_open_pulls"`
		NumMergedPulls int    `json:"num_merged_pulls"`
		GitURL         string `json:"git_url"`
	}
	// harness commit status check report input.
	checkReportInput struct {
		Identifier string       `json:"identifier"`
		Status     string       `json:"status"`
		Summary    string       `json:"summary,omitempty"`
		Link       string       `json:"link,omitempty"`
		Payload    checkPayload `json:"payload"`
	}
	checkPayload struct {
		Kind string `json:"kind"`
	}
	// harness commit status check resource.
	check struct {
		Identifier string `json:"identifier"`
		Status     string `json:"status"`
		Summary    string `json:"summary"`
		Link       string `json:"link"`
	}
	hook struct {
		Created               int      `json:"created"`
		CreatedBy             int      `json:"created_by"`
		Description           string   `json:"description"`
		Enabled               bool     `json:"enabled"`
		HasSecret             bool     `json:"has_secret"`
		Secret                string   `json:"secret"`
		Identifier            string   `json:"identifier"`
		Insecure              bool     `json:"insecure"`
		LatestExecutionResult string   `json:"latest_execution_result"`
		ParentID              int      `json:"parent_id"`
		ParentType            string   `json:"parent_type"`
		Triggers              []string `json:"triggers"`
		Updated               int      `json:"updated"`
		URL                   string   `json:"url"`
		Version               int      `json:"version"`
	}
)

//
// native data structure conversion
//

func convertRepositoryList(src []*repository) []*scm.Repository {
	var dst []*scm.Repository
	for _, v := range src {
		dst = append(dst, convertRepository(v))
	}
	return dst
}

func convertRepository(src *repository) *scm.Repository {
	return &scm.Repository{
		ID:        strconv.Itoa(src.ID),
		Namespace: src.Path,
		Name:      src.UID,
		Branch:    src.DefaultBranch,
		Private:   !src.IsPublic,
		Clone:     src.GitURL,
		CloneSSH:  src.GitURL,
		Link:      src.GitURL,
		// Created:   time.Unix(src.Created, 0),
		//		Updated:   time.Unix(src.Updated, 0),
	}
}

func convertCheck(from *check) *scm.Status {
	return &scm.Status{
		State:  convertFromHarnessCheckStatus(from.Status),
		Label:  from.Identifier,
		Desc:   from.Summary,
		Target: from.Link,
	}
}

// convertToHarnessCheckStatus maps a scm.State to the Harness Code
// status check status string.
func convertToHarnessCheckStatus(from scm.State) string {
	switch from {
	case scm.StatePending:
		return "pending"
	case scm.StateRunning:
		return "running"
	case scm.StateSuccess:
		return "success"
	case scm.StateFailure:
		return "failure"
	default:
		return "error"
	}
}

// convertFromHarnessCheckStatus maps a Harness Code status check status
// string to a scm.State.
func convertFromHarnessCheckStatus(from string) scm.State {
	switch from {
	case "pending":
		return scm.StatePending
	case "running":
		return scm.StateRunning
	case "success", "failure_ignored":
		return scm.StateSuccess
	case "failure":
		return scm.StateFailure
	case "error":
		return scm.StateError
	default:
		return scm.StateUnknown
	}
}

func convertHookList(from []*hook) []*scm.Hook {
	to := []*scm.Hook{}
	for _, v := range from {
		to = append(to, convertHook(v))
	}
	return to
}

func convertHook(from *hook) *scm.Hook {
	return &scm.Hook{
		// keeping id same as name
		ID:         from.Identifier,
		Name:       from.Identifier,
		Active:     from.Enabled,
		Target:     from.URL,
		Events:     from.Triggers,
		SkipVerify: from.Insecure,
	}
}
