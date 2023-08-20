// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/drone/go-scm/scm"
)

type gitService struct {
	client *wrapper
}

func (s *gitService) CreateBranch(ctx context.Context, repo string, params *scm.ReferenceInput) (*scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/branches", harnessURI)
	in := &branchInput{
		Name:   params.Name,
		Target: params.Sha,
	}
	return s.client.do(ctx, "POST", path, in, nil)
}

func (s *gitService) FindBranch(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/branches/%s", harnessURI, name)
	out := new(branch)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertBranch(out), res, err
}

func (s *gitService) FindCommit(ctx context.Context, repo, ref string) (*scm.Commit, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/commits/%s", harnessURI, ref)
	out := new(commitInfo)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertCommitInfo(out), res, err
}

func (s *gitService) FindTag(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListBranches(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/branches?%s", harnessURI, encodeListOptions(opts))
	out := []*branch{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertBranchList(out), res, err
}

func (s *gitService) ListBranchesV2(ctx context.Context, repo string, opts scm.BranchListOptions) ([]*scm.Reference, *scm.Response, error) {
	// Harness doesnt provide support listing based on searchTerm
	// Hence calling the ListBranches
	return s.ListBranches(ctx, repo, opts.PageListOptions)
}

func (s *gitService) ListCommits(ctx context.Context, repo string, _ scm.CommitListOptions) ([]*scm.Commit, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/commits", harnessURI)
	out := []*commitInfo{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertCommitList(out), res, err
}

func (s *gitService) ListTags(ctx context.Context, repo string, _ scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListChanges(ctx context.Context, repo, ref string, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) CompareChanges(ctx context.Context, repo, source, target string, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	harnessURI := buildHarnessURI(s.client.account, s.client.organization, s.client.project, repo)
	path := fmt.Sprintf("api/v1/repos/%s/compare/%s...%s", harnessURI, source, target)
	res, err := s.client.do(ctx, "GET", path, nil, nil)
	// convert response to a string
	buf := new(strings.Builder)
	_, _ = io.Copy(buf, res.Body)

	return convertCompareChanges(buf.String()), res, err
}

// native data structures
type (
	commitInfo struct {
		Author struct {
			Identity struct {
				Email string `json:"email"`
				Name  string `json:"name"`
			} `json:"identity"`
			When time.Time `json:"when"`
		} `json:"author"`
		Committer struct {
			Identity struct {
				Email string `json:"email"`
				Name  string `json:"name"`
			} `json:"identity"`
			When time.Time `json:"when"`
		} `json:"committer"`
		Message string `json:"message"`
		Sha     string `json:"sha"`
		Title   string `json:"title"`
	}
	branchInput struct {
		Name   string `json:"name"`
		Target string `json:"target"`
	}
	branch struct {
		Commit struct {
			Author struct {
				Identity struct {
					Email string `json:"email"`
					Name  string `json:"name"`
				} `json:"identity"`
				When time.Time `json:"when"`
			} `json:"author"`
			Committer struct {
				Identity struct {
					Email string `json:"email"`
					Name  string `json:"name"`
				} `json:"identity"`
				When time.Time `json:"when"`
			} `json:"committer"`
			Message string `json:"message"`
			Sha     string `json:"sha"`
			Title   string `json:"title"`
		} `json:"commit"`
		Name string `json:"name"`
		Sha  string `json:"sha"`
	}
)

//
// native data structure conversion
//

func convertBranchList(src []*branch) []*scm.Reference {
	dst := []*scm.Reference{}
	for _, v := range src {
		dst = append(dst, convertBranch(v))
	}
	return dst
}

func convertBranch(src *branch) *scm.Reference {
	return &scm.Reference{
		Name: src.Name,
		Path: scm.ExpandRef(src.Name, "refs/heads/"),
		Sha:  src.Sha,
	}
}

func convertCommitList(src []*commitInfo) []*scm.Commit {
	dst := []*scm.Commit{}
	for _, v := range src {
		dst = append(dst, convertCommitInfo(v))
	}
	return dst
}

func convertCompareChanges(src string) []*scm.Change {
	files, _, err := gitdiff.Parse(strings.NewReader(src))
	if err != nil {
		return nil
	}

	changes := make([]*scm.Change, 0)
	for _, f := range files {
		changes = append(changes, &scm.Change{
			Path:         f.NewName,
			PrevFilePath: f.OldName,
			Added:        f.IsNew,
			Deleted:      f.IsDelete,
			Renamed:      f.IsRename,
		})
	}

	return changes
}

func convertCommitInfo(src *commitInfo) *scm.Commit {
	return &scm.Commit{
		Sha:     src.Sha,
		Message: src.Message,
		Author: scm.Signature{
			Name:  src.Author.Identity.Name,
			Email: src.Author.Identity.Email,
			Date:  src.Author.When,
		},
		Committer: scm.Signature{
			Name:  src.Committer.Identity.Name,
			Email: src.Committer.Identity.Email,
			Date:  src.Committer.When,
		},
	}
}
