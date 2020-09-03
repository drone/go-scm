// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"time"

	"github.com/jenkins-x/go-scm/scm"
)

type gitService struct {
	client *wrapper
}

func (s *gitService) FindRef(ctx context.Context, repo, ref string) (string, *scm.Response, error) {
	namespace, name := scm.Split(repo)

	out, err := s.client.GiteaClient.GetRepoRef(namespace, name, ref)
	if err != nil || out.Object == nil {
		return "", nil, err
	}
	return out.Object.SHA, nil, nil
}

func (s *gitService) CreateRef(ctx context.Context, repo, ref, sha string) (*scm.Reference, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) DeleteRef(ctx context.Context, repo, ref string) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *gitService) FindBranch(ctx context.Context, repo, branchName string) (*scm.Reference, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	out, err := s.client.GiteaClient.GetRepoBranch(namespace, name, branchName)
	return convertBranch(out), nil, err
}

func (s *gitService) FindCommit(ctx context.Context, repo, ref string) (*scm.Commit, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	out, err := s.client.GiteaClient.GetSingleCommit(namespace, name, ref)
	return convertCommit(out), nil, err
}

func (s *gitService) FindTag(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListBranches(ctx context.Context, repo string, _ scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	out, err := s.client.GiteaClient.ListRepoBranches(namespace, name, gitea.ListRepoBranchesOptions{})
	return convertBranchList(out), nil, err
}

func (s *gitService) ListCommits(ctx context.Context, repo string, _ scm.CommitListOptions) ([]*scm.Commit, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListTags(ctx context.Context, repo string, _ scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListChanges(ctx context.Context, repo, ref string, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

//
// native data structures
//

type (
	// gitea commit object.
	commit struct {
		ID        string    `json:"id"`
		Sha       string    `json:"sha"`
		Message   string    `json:"message"`
		URL       string    `json:"url"`
		Author    signature `json:"author"`
		Committer signature `json:"committer"`
		Timestamp time.Time `json:"timestamp"`
	}

	// gitea signature object.
	signature struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}
)

//
// native data structure conversion
//

func convertBranchList(src []*gitea.Branch) []*scm.Reference {
	dst := []*scm.Reference{}
	for _, v := range src {
		dst = append(dst, convertBranch(v))
	}
	return dst
}

func convertBranch(src *gitea.Branch) *scm.Reference {
	if src == nil || src.Commit == nil {
		return nil
	}
	return &scm.Reference{
		Name: scm.TrimRef(src.Name),
		Path: scm.ExpandRef(src.Name, "refs/heads/"),
		Sha:  src.Commit.ID,
	}
}

func convertCommit(src *gitea.Commit) *scm.Commit {
	if src == nil || src.RepoCommit == nil {
		return nil
	}
	return &scm.Commit{
		Sha:       src.SHA,
		Link:      src.URL,
		Message:   src.RepoCommit.Message,
		Author:    convertUserSignature(src.Author),
		Committer: convertUserSignature(src.Committer),
	}
}

func convertUserSignature(src *gitea.User) scm.Signature {
	if src == nil {
		return scm.Signature{}
	}
	return scm.Signature{
		Login:  src.UserName,
		Email:  src.Email,
		Name:   src.FullName,
		Avatar: src.AvatarURL,
	}
}
