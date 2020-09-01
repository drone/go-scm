// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"bytes"
	"code.gitea.io/sdk/gitea"
	"context"
	"fmt"
	"time"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/jenkins-x/go-scm/scm"
)

type pullService struct {
	*issueService
}

func (s *pullService) Find(ctx context.Context, repo string, index int) (*scm.PullRequest, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	out, err := s.client.GiteaClient.GetPullRequest(namespace, name, int64(index))
	return convertPullRequest(out), nil, err
}

func (s *pullService) List(ctx context.Context, repo string, opts scm.PullRequestListOptions) ([]*scm.PullRequest, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	out, err := s.client.GiteaClient.ListRepoPullRequests(namespace, name, gitea.ListPullRequestsOptions{})
	return convertPullRequests(out), nil, err
}

// TODO: Maybe contribute to gitea/go-sdk with .patch function?
func (s *pullService) ListChanges(ctx context.Context, repo string, number int, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	// Get the patch and then parse it.
	path := fmt.Sprintf("api/v1/repos/%s/pulls/%d.patch", repo, number)
	buf := new(bytes.Buffer)
	res, err := s.client.do(ctx, "GET", path, nil, buf)
	if err != nil {
		return nil, res, err
	}
	changedFiles, _, err := gitdiff.Parse(buf)
	if err != nil {
		return nil, res, err
	}
	var changes []*scm.Change
	for _, c := range changedFiles {
		var linesAdded int64
		var linesDeleted int64

		for _, tf := range c.TextFragments {
			linesAdded += tf.LinesAdded
			linesDeleted += tf.LinesDeleted
		}
		changes = append(changes, &scm.Change{
			Path:         c.NewName,
			PreviousPath: c.OldName,
			Added:        c.IsNew,
			Renamed:      c.IsRename,
			Deleted:      c.IsDelete,
			Additions:    int(linesAdded),
			Deletions:    int(linesDeleted),
		})
	}
	return changes, res, nil
}

func (s *pullService) Merge(ctx context.Context, repo string, index int, options *scm.PullRequestMergeOptions) (*scm.Response, error) {
	namespace, name := scm.Split(repo)
	in := gitea.MergePullRequestOption{}

	if options != nil {
		in.Style = convertMergeMethodToMergeStyle(options.MergeMethod)
		in.Title = options.CommitTitle
	}

	_, err := s.client.GiteaClient.MergePullRequest(namespace, name, int64(index), in)
	return nil, err
}

func (s *pullService) Update(ctx context.Context, repo string, number int, input *scm.PullRequestInput) (*scm.PullRequest, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	in := gitea.EditPullRequestOption{
		Title: input.Title,
		Body:  input.Body,
		Base:  input.Base,
	}
	out, err := s.client.GiteaClient.EditPullRequest(namespace, name, int64(number), in)
	return convertPullRequest(out), nil, err
}

func (s *pullService) Close(ctx context.Context, repo string, number int) (*scm.Response, error) {
	namespace, name := scm.Split(repo)
	closed := gitea.StateClosed
	in := gitea.EditPullRequestOption{
		State: &closed,
	}
	_, err := s.client.GiteaClient.EditPullRequest(namespace, name, int64(number), in)
	return nil, err
}

func (s *pullService) Reopen(ctx context.Context, repo string, number int) (*scm.Response, error) {
	namespace, name := scm.Split(repo)
	reopen := gitea.StateOpen
	in := gitea.EditPullRequestOption{
		State: &reopen,
	}
	_, err := s.client.GiteaClient.EditPullRequest(namespace, name, int64(number), in)
	return nil, err
}

func (s *pullService) Create(ctx context.Context, repo string, input *scm.PullRequestInput) (*scm.PullRequest, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	in := gitea.CreatePullRequestOption{
		Head:  input.Head,
		Base:  input.Base,
		Title: input.Title,
		Body:  input.Body,
	}
	out, err := s.client.GiteaClient.CreatePullRequest(namespace, name, in)
	return convertPullRequest(out), nil, err
}

func (s *pullService) RequestReview(ctx context.Context, repo string, number int, logins []string) (*scm.Response, error) {
	return s.AssignIssue(ctx, repo, number, logins)
}

func (s *pullService) UnrequestReview(ctx context.Context, repo string, number int, logins []string) (*scm.Response, error) {
	return s.UnassignIssue(ctx, repo, number, logins)
}

//
// native data structures
//

type pullRequest struct {
	ID         int        `json:"id"`
	Number     int        `json:"number"`
	User       user       `json:"user"`
	Title      string     `json:"title"`
	Body       string     `json:"body"`
	State      string     `json:"state"`
	HeadBranch string     `json:"head_branch"`
	HeadRepo   repository `json:"head_repo"`
	Head       reference  `json:"head"`
	BaseBranch string     `json:"base_branch"`
	BaseRepo   repository `json:"base_repo"`
	Base       reference  `json:"base"`
	HTMLURL    string     `json:"html_url"`
	Mergeable  bool       `json:"mergeable"`
	Merged     bool       `json:"merged"`
	Created    time.Time  `json:"created_at"`
	Updated    time.Time  `json:"updated_at"`
}

type reference struct {
	Repo repository `json:"repo"`
	Name string     `json:"ref"`
	Sha  string     `json:"sha"`
}

//
// native data structure conversion
//

func convertPullRequests(src []*gitea.PullRequest) []*scm.PullRequest {
	dst := []*scm.PullRequest{}
	for _, v := range src {
		dst = append(dst, convertPullRequest(v))
	}
	return dst
}

func convertPullRequest(src *gitea.PullRequest) *scm.PullRequest {
	if src == nil || src.Title == "" {
		return nil
	}
	return &scm.PullRequest{
		Number:    int(src.Index),
		Title:     src.Title,
		Body:      src.Body,
		Sha:       src.Head.Sha,
		Source:    src.Head.Name,
		Target:    src.Base.Name,
		Head:      *convertPullRequestBranch(src.Head),
		Base:      *convertPullRequestBranch(src.Base),
		Link:      src.HTMLURL,
		Fork:      src.Base.Repository.FullName,
		Ref:       fmt.Sprintf("refs/pull/%d/head", src.Index),
		Closed:    src.State == gitea.StateClosed,
		Author:    *convertGiteaUser(src.Poster),
		Labels:    convertGiteaLabels(src.Labels),
		Merged:    src.HasMerged,
		Mergeable: src.Mergeable,
		Created:   *src.Created,
		Updated:   *src.Updated,
	}
}

func convertPullRequestFromIssue(src *gitea.Issue) *scm.PullRequest {
	return &scm.PullRequest{
		Number:  int(src.Index),
		Title:   src.Title,
		Body:    src.Body,
		Closed:  src.State == gitea.StateClosed,
		Author:  *convertGiteaUser(src.Poster),
		Merged:  src.PullRequest.HasMerged,
		Created: src.Created,
		Updated: src.Updated,
	}
}

func convertPullRequestBranch(src *gitea.PRBranchInfo) *scm.PullRequestBranch {
	return &scm.PullRequestBranch{
		Ref:  src.Ref,
		Sha:  src.Sha,
		Repo: *convertGiteaRepository(src.Repository),
	}
}

func convertMergeMethodToMergeStyle(mm string) gitea.MergeStyle {
	switch mm {
	case "merge":
		return gitea.MergeStyleMerge
	case "rebase":
		return gitea.MergeStyleRebase
	case "rebase-merge":
		return gitea.MergeStyleRebaseMerge
	case "squash":
		return gitea.MergeStyleSquash
	default:
		return gitea.MergeStyleMerge
	}
}
