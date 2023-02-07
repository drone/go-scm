// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/drone/go-scm/scm"
)

type repositoryService struct {
	client *wrapper
}

func (s *repositoryService) Find(ctx context.Context, repo string) (*scm.Repository, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) FindHook(ctx context.Context, repo string, id string) (*scm.Hook, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) FindPerms(ctx context.Context, repo string) (*scm.Perm, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, "")
	path := fmt.Sprintf("api/v1/spaces/%s/repos?sort=path&order=asc&%s", harnessURI, encodeListOptions(opts))
	out := []*repository{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertRepositoryList(out), res, err
}

func (s *repositoryService) ListHooks(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Hook, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) ListStatus(ctx context.Context, repo string, ref string, opts scm.ListOptions) ([]*scm.Status, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) CreateHook(ctx context.Context, repo string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) CreateStatus(ctx context.Context, repo string, ref string, input *scm.StatusInput) (*scm.Status, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) UpdateHook(ctx context.Context, repo, id string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) DeleteHook(ctx context.Context, repo string, id string) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
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

	// gitea permissions details.
	perm struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	}

	// gitea hook resource.
	hook struct {
		ID     int        `json:"id"`
		Type   string     `json:"type"`
		Events []string   `json:"events"`
		Active bool       `json:"active"`
		Config hookConfig `json:"config"`
	}

	// gitea hook configuration details.
	hookConfig struct {
		URL         string `json:"url"`
		ContentType string `json:"content_type"`
		Secret      string `json:"secret"`
	}

	// gitea status resource.
	status struct {
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		State       string    `json:"status"`
		TargetURL   string    `json:"target_url"`
		Description string    `json:"description"`
		Context     string    `json:"context"`
	}

	// gitea status creation request.
	statusInput struct {
		State       string `json:"state"`
		TargetURL   string `json:"target_url"`
		Description string `json:"description"`
		Context     string `json:"context"`
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
func convertHookList(src []*hook) []*scm.Hook {
	var dst []*scm.Hook
	for _, v := range src {
		dst = append(dst, convertHook(v))
	}
	return dst
}

func convertHook(from *hook) *scm.Hook {
	return &scm.Hook{
		ID:     strconv.Itoa(from.ID),
		Active: from.Active,
		Target: from.Config.URL,
		Events: from.Events,
	}
}

func convertHookEvent(from scm.HookEvents) []string {
	var events []string
	if from.PullRequest {
		events = append(events, "pull_request")
	}
	if from.Issue {
		events = append(events, "issues")
	}
	if from.IssueComment || from.PullRequestComment {
		events = append(events, "issue_comment")
	}
	if from.Branch || from.Tag {
		events = append(events, "create")
		events = append(events, "delete")
	}
	if from.Push {
		events = append(events, "push")
	}
	return events
}

func convertStatusList(src []*status) []*scm.Status {
	var dst []*scm.Status
	for _, v := range src {
		dst = append(dst, convertStatus(v))
	}
	return dst
}

func convertStatus(from *status) *scm.Status {
	return &scm.Status{
		State:  convertState(from.State),
		Label:  from.Context,
		Desc:   from.Description,
		Target: from.TargetURL,
	}
}

func convertState(from string) scm.State {
	switch from {
	case "error":
		return scm.StateError
	case "failure":
		return scm.StateFailure
	case "pending":
		return scm.StatePending
	case "success":
		return scm.StateSuccess
	default:
		return scm.StateUnknown
	}
}

func convertFromState(from scm.State) string {
	switch from {
	case scm.StatePending, scm.StateRunning:
		return "pending"
	case scm.StateSuccess:
		return "success"
	case scm.StateFailure:
		return "failure"
	default:
		return "error"
	}
}
