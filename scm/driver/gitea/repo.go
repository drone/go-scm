// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"net/url"
	"strconv"
	"time"

	"github.com/jenkins-x/go-scm/scm"
)

type repositoryService struct {
	client *wrapper
}

func (s *repositoryService) Create(_ context.Context, input *scm.RepositoryInput) (*scm.Repository, *scm.Response, error) {
	var out *gitea.Repository
	var err error
	in := gitea.CreateRepoOption{
		Name:        input.Name,
		Description: input.Description,
		Private:     input.Private,
	}

	if input.Namespace == "" {
		out, err = s.client.GiteaClient.CreateRepo(in)
	} else {
		out, err = s.client.GiteaClient.CreateOrgRepo(input.Namespace, in)
	}
	return convertGiteaRepository(out), nil, err
}

func (s *repositoryService) Fork(context.Context, *scm.RepositoryInput, string) (*scm.Repository, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) FindCombinedStatus(_ context.Context, repo, ref string) (*scm.CombinedStatus, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	out, err := s.client.GiteaClient.GetCombinedStatus(namespace, name, ref)
	if err != nil {
		return nil, nil, err
	}
	return &scm.CombinedStatus{
		State:    convertState(out.State),
		Sha:      out.SHA,
		Statuses: convertStatusList(out.Statuses),
	}, nil, nil
}

func (s *repositoryService) FindUserPermission(_ context.Context, repo string, user string) (string, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	members, err := s.client.GiteaClient.ListCollaborators(namespace, name, gitea.ListCollaboratorsOptions{})
	if err != nil {
		return "", nil, err
	}
	for _, m := range members {
		if m.UserName == user {
			if m.IsAdmin {
				return scm.AdminPermission, nil, nil
			}
			return scm.WritePermission, nil, nil
		}
	}

	return scm.NoPermission, nil, nil
}

func (s *repositoryService) AddCollaborator(_ context.Context, repo, user, permission string) (bool, bool, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	opt := gitea.AddCollaboratorOption{Permission: &permission}
	err := s.client.GiteaClient.AddCollaborator(namespace, name, user, opt)
	if err != nil {
		return false, false, nil, err
	}
	return true, false, nil, nil
}

func (s *repositoryService) IsCollaborator(_ context.Context, repo, user string) (bool, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	isCollab, err := s.client.GiteaClient.IsCollaborator(namespace, name, user)
	return isCollab, nil, err
}

func (s *repositoryService) ListCollaborators(_ context.Context, repo string, ops scm.ListOptions) ([]scm.User, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	out, err := s.client.GiteaClient.ListCollaborators(namespace, name, gitea.ListCollaboratorsOptions{})
	return convertGiteaUsers(out), nil, err
}

func (s *repositoryService) ListLabels(_ context.Context, repo string, _ scm.ListOptions) ([]*scm.Label, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	out, err := s.client.GiteaClient.ListRepoLabels(namespace, name, gitea.ListLabelsOptions{})
	return convertGiteaLabels(out), nil, err
}

func (s *repositoryService) Find(_ context.Context, repo string) (*scm.Repository, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	out, err := s.client.GiteaClient.GetRepo(namespace, name)
	return convertGiteaRepository(out), nil, err
}

func (s *repositoryService) FindHook(_ context.Context, repo string, id string) (*scm.Hook, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, nil, err
	}
	out, err := s.client.GiteaClient.GetRepoHook(namespace, name, idInt)
	return convertHook(out), nil, err
}

func (s *repositoryService) FindPerms(ctx context.Context, repo string) (*scm.Perm, *scm.Response, error) {
	r, _, err := s.Find(ctx, repo)
	if err != nil || r == nil {
		return nil, nil, err
	}
	return r.Perm, nil, err
}

func (s *repositoryService) List(_ context.Context, _ scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	out, err := s.client.GiteaClient.ListMyRepos(gitea.ListReposOptions{})
	return convertRepositoryList(out), nil, err
}

func (s *repositoryService) ListOrganisation(_ context.Context, org string, opts scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	out, err := s.client.GiteaClient.ListOrgRepos(org, gitea.ListOrgReposOptions{})
	return convertRepositoryList(out), nil, err
}

func (s *repositoryService) ListUser(_ context.Context, username string, opts scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	out, err := s.client.GiteaClient.ListUserRepos(username, gitea.ListReposOptions{})
	return convertRepositoryList(out), nil, err
}

func (s *repositoryService) ListHooks(_ context.Context, repo string, _ scm.ListOptions) ([]*scm.Hook, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	out, err := s.client.GiteaClient.ListRepoHooks(namespace, name, gitea.ListHooksOptions{})
	return convertHookList(out), nil, err
}

func (s *repositoryService) ListStatus(_ context.Context, repo string, ref string, _ scm.ListOptions) ([]*scm.Status, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	out, err := s.client.GiteaClient.ListStatuses(namespace, name, ref, gitea.ListStatusesOption{})
	return convertStatusList(out), nil, err
}

func (s *repositoryService) CreateHook(_ context.Context, repo string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	target, err := url.Parse(input.Target)
	if err != nil {
		return nil, nil, err
	}
	params := target.Query()
	params.Set("secret", input.Secret)
	target.RawQuery = params.Encode()

	namespace, name := scm.Split(repo)
	in := gitea.CreateHookOption{
		Type: "gitea",
		Config: map[string]string{
			"secret":       input.Secret,
			"content_type": "json",
			"url":          target.String(),
		},
		Events: append(
			input.NativeEvents,
			convertHookEvent(input.Events)...,
		),
		Active: true,
	}
	out, err := s.client.GiteaClient.CreateRepoHook(namespace, name, in)
	return convertHook(out), nil, err
}

func (s *repositoryService) CreateStatus(_ context.Context, repo string, ref string, input *scm.StatusInput) (*scm.Status, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	in := gitea.CreateStatusOption{
		State:       convertFromState(input.State),
		TargetURL:   input.Target,
		Description: input.Desc,
		Context:     input.Label,
	}
	out, err := s.client.GiteaClient.CreateStatus(namespace, name, ref, in)
	return convertStatus(out), nil, err
}

func (s *repositoryService) DeleteHook(_ context.Context, repo string, id string) (*scm.Response, error) {
	namespace, name := scm.Split(repo)
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	err = s.client.GiteaClient.DeleteRepoHook(namespace, name, idInt)
	return nil, err
}

//
// native data structures
//

type (
	// gitea repository resource.
	repository struct {
		ID            int       `json:"id"`
		Owner         user      `json:"owner"`
		Name          string    `json:"name"`
		FullName      string    `json:"full_name"`
		Private       bool      `json:"private"`
		Fork          bool      `json:"fork"`
		HTMLURL       string    `json:"html_url"`
		SSHURL        string    `json:"ssh_url"`
		CloneURL      string    `json:"clone_url"`
		DefaultBranch string    `json:"default_branch"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		Permissions   perm      `json:"permissions"`
	}

	// gitea permissions details.
	perm struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	}
)

//
// native data structure conversion
//

func convertRepositoryList(src []*gitea.Repository) []*scm.Repository {
	var dst []*scm.Repository
	for _, v := range src {
		dst = append(dst, convertGiteaRepository(v))
	}
	return dst
}

func convertGiteaRepository(src *gitea.Repository) *scm.Repository {
	if src == nil || src.Owner == nil {
		return nil
	}
	return &scm.Repository{
		ID:        strconv.FormatInt(src.ID, 10),
		Namespace: src.Owner.UserName,
		Name:      src.Name,
		FullName:  src.FullName,
		Perm:      convertGiteaPerm(src.Permissions),
		Branch:    src.DefaultBranch,
		Private:   src.Private,
		Clone:     src.CloneURL,
		CloneSSH:  src.SSHURL,
		Link:      src.HTMLURL,
		Created:   src.Created,
		Updated:   src.Updated,
	}
}

func convertGiteaPerm(src *gitea.Permission) *scm.Perm {
	if src == nil {
		return nil
	}
	return &scm.Perm{
		Push:  src.Push,
		Pull:  src.Pull,
		Admin: src.Admin,
	}
}

func convertRepository(src *repository) *scm.Repository {
	return &scm.Repository{
		ID:        strconv.Itoa(src.ID),
		Namespace: userLogin(&src.Owner),
		Name:      src.Name,
		FullName:  src.FullName,
		Perm:      convertPerm(src.Permissions),
		Branch:    src.DefaultBranch,
		Private:   src.Private,
		Clone:     src.CloneURL,
		CloneSSH:  src.SSHURL,
	}
}

func convertPerm(src perm) *scm.Perm {
	return &scm.Perm{
		Push:  src.Push,
		Pull:  src.Pull,
		Admin: src.Admin,
	}
}

func convertHookList(src []*gitea.Hook) []*scm.Hook {
	var dst []*scm.Hook
	for _, v := range src {
		dst = append(dst, convertHook(v))
	}
	return dst
}

func convertHook(from *gitea.Hook) *scm.Hook {
	return &scm.Hook{
		ID:     strconv.FormatInt(from.ID, 10),
		Active: from.Active,
		Target: from.Config["url"],
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

func convertStatusList(src []*gitea.Status) []*scm.Status {
	var dst []*scm.Status
	for _, v := range src {
		dst = append(dst, convertStatus(v))
	}
	return dst
}

func convertStatus(from *gitea.Status) *scm.Status {
	return &scm.Status{
		State:  convertState(from.State),
		Label:  from.Context,
		Desc:   from.Description,
		Target: from.TargetURL,
	}
}

func convertState(from gitea.StatusState) scm.State {
	switch from {
	case gitea.StatusError:
		return scm.StateError
	case gitea.StatusFailure:
		return scm.StateFailure
	case gitea.StatusPending:
		return scm.StatePending
	case gitea.StatusSuccess:
		return scm.StateSuccess
	default:
		return scm.StateUnknown
	}
}

func convertFromState(from scm.State) gitea.StatusState {
	switch from {
	case scm.StatePending, scm.StateRunning:
		return gitea.StatusPending
	case scm.StateSuccess:
		return gitea.StatusSuccess
	case scm.StateFailure:
		return gitea.StatusFailure
	default:
		return gitea.StatusError
	}
}
