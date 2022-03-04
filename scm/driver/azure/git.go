// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"context"
	"fmt"
	"time"

	"github.com/drone/go-scm/scm"
)

type gitService struct {
	client *wrapper
}

func (s *gitService) CreateBranch(ctx context.Context, repo string, params *scm.CreateBranch) (*scm.Response, error) {
	// https://docs.microsoft.com/en-us/rest/api/azure/devops/git/refs/update-refs?view=azure-devops-rest-6.0
	endpoint := fmt.Sprintf("%s/%s/_apis/git/repositories/%s/refs?api-version=6.0", s.client.owner, s.client.project, repo)

	in := make(crudBranch, 1)
	in[0].Name = scm.ExpandRef(params.Name, "refs/heads")
	in[0].NewObjectID = params.Sha
	in[0].OldObjectID = "0000000000000000000000000000000000000000"
	return s.client.do(ctx, "POST", endpoint, in, nil)
}

func (s *gitService) FindBranch(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) FindCommit(ctx context.Context, repo, ref string) (*scm.Commit, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) FindTag(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListBranches(ctx context.Context, repo string, _ scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	// https://docs.microsoft.com/en-us/rest/api/azure/devops/git/refs/list?view=azure-devops-rest-6.0
	endpoint := fmt.Sprintf("%s/%s/_apis/git/repositories/%s/refs?api-version=6.0", s.client.owner, s.client.project, repo)
	out := new(branchList)
	res, err := s.client.do(ctx, "GET", endpoint, nil, &out)
	return convertBranchList(out.Value), res, err
}

func (s *gitService) ListCommits(ctx context.Context, repo string, opts scm.CommitListOptions) ([]*scm.Commit, *scm.Response, error) {
	// https://docs.microsoft.com/en-us/rest/api/azure/devops/git/commits/get-commits?view=azure-devops-rest-6.0
	endpoint := fmt.Sprintf("%s/%s/_apis/git/repositories/%s/commits?", s.client.owner, s.client.project, repo)
	if opts.Ref != "" {
		endpoint += fmt.Sprintf("searchCriteria.itemVersion.version=%s&api-version=6.0", opts.Ref)
	} else {
		endpoint += "&api-version=6.0"
	}
	out := new(commitList)
	res, err := s.client.do(ctx, "GET", endpoint, nil, &out)
	return convertCommitList(out.Value), res, err
}

func (s *gitService) ListTags(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListChanges(ctx context.Context, repo, ref string, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) CompareChanges(ctx context.Context, repo, source, target string, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

type crudBranch []struct {
	Name        string `json:"name"`
	OldObjectID string `json:"oldObjectId"`
	NewObjectID string `json:"newObjectId"`
}

type branchList struct {
	Value []*branch `json:"value"`
	Count int       `json:"count"`
}

type branch struct {
	Name     string `json:"name"`
	ObjectID string `json:"objectId"`
	Creator  struct {
		DisplayName string `json:"displayName"`
		URL         string `json:"url"`
		Links       struct {
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"_links"`
		ID         string `json:"id"`
		UniqueName string `json:"uniqueName"`
		ImageURL   string `json:"imageUrl"`
		Descriptor string `json:"descriptor"`
	} `json:"creator"`
	URL string `json:"url"`
}

type commitList struct {
	Value []*gitCommit `json:"value"`
	Count int          `json:"count"`
}
type gitCommit struct {
	CommitID string `json:"commitId"`
	Author   struct {
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Date  time.Time `json:"date"`
	} `json:"author"`
	Committer struct {
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Date  time.Time `json:"date"`
	} `json:"committer"`
	Comment          string `json:"comment"`
	CommentTruncated bool   `json:"commentTruncated"`
	ChangeCounts     struct {
		Add    int `json:"Add"`
		Edit   int `json:"Edit"`
		Delete int `json:"Delete"`
	} `json:"changeCounts"`
	URL       string `json:"url"`
	RemoteURL string `json:"remoteUrl"`
}

func convertBranchList(from []*branch) []*scm.Reference {
	to := []*scm.Reference{}
	for _, v := range from {
		to = append(to, convertBranch(v))
	}
	return to
}

func convertBranch(from *branch) *scm.Reference {
	return &scm.Reference{
		Name: scm.TrimRef(from.Name),
		Path: from.Name,
		Sha:  from.ObjectID,
	}
}

func convertCommitList(from []*gitCommit) []*scm.Commit {
	to := []*scm.Commit{}
	for _, v := range from {
		to = append(to, convertCommit(v))
	}
	return to
}

func convertCommit(from *gitCommit) *scm.Commit {
	return &scm.Commit{
		Message: from.Comment,
		Sha:     from.CommitID,
		Link:    from.URL,
		Author: scm.Signature{
			Name:  from.Author.Name,
			Email: from.Author.Email,
			Date:  from.Author.Date,
		},
		Committer: scm.Signature{
			Name:  from.Committer.Name,
			Email: from.Committer.Email,
			Date:  from.Committer.Date,
		},
	}
}
