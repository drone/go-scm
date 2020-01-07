// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"
	"fmt"
	"time"

	"github.com/drone/go-scm/scm"
)

// TODO(bradrydzewski) commit link is an empty string.
// TODO(bradrydzewski) commit list is not supported; behavior incompatible with other drivers.

type gitService struct {
	client *wrapper
}

func (s *gitService) FindBranch(ctx context.Context, repo, branch string) (*scm.Reference, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/branches?filterText=%s", namespace, name, branch)
	out := new(branches)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	if err != nil {
		return nil, res, err
	}
	for _, v := range out.Values {
		if v.DisplayID == branch {
			return convertBranch(v), res, err
		}
	}
	return nil, res, scm.ErrNotFound
}

func (s *gitService) FindCommit(ctx context.Context, repo, ref string) (*scm.Commit, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/commits/%s", namespace, name, ref)
	out := new(commit)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertCommit(out), res, err
}

func (s *gitService) FindTag(ctx context.Context, repo, tag string) (*scm.Reference, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/tags?filterText=%s", namespace, name, tag)
	out := new(branches)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	if err != nil {
		return nil, res, err
	}
	for _, v := range out.Values {
		if v.DisplayID == tag {
			return convertTag(v), res, err
		}
	}
	return nil, res, scm.ErrNotFound
}

func (s *gitService) ListBranches(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/branches?%s", namespace, name, encodeListOptions(opts))
	out := new(branches)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	if !out.pagination.LastPage.Bool {
		res.Page.First = 1
		res.Page.Next = opts.Page + 1
	}
	return convertBranchList(out), res, err
}

func (s *gitService) ListCommits(ctx context.Context, repo string, opts scm.CommitListOptions) ([]*scm.Commit, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListTags(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/tags?%s", namespace, name, encodeListOptions(opts))
	out := new(branches)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	if !out.pagination.LastPage.Bool {
		res.Page.First = 1
		res.Page.Next = opts.Page + 1
	}
	return convertTagList(out), res, err
}

func (s *gitService) ListChanges(ctx context.Context, repo, ref string, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/commits/%s/changes?%s", namespace, name, ref, encodeListOptions(opts))
	out := new(diffstats)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	if !out.pagination.LastPage.Bool {
		res.Page.First = 1
		res.Page.Next = opts.Page + 1
	}
	return convertDiffstats(out), res, err
}

func (s *gitService) CompareChanges(ctx context.Context, repo, source, target string, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/compare/changes?from=%s&to=%s&%s", namespace, name, source, target, encodeListOptions(opts))
	out := new(diffstats)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	if !out.pagination.LastPage.Bool {
		res.Page.First = 1
		res.Page.Next = opts.Page + 1
	}
	return convertDiffstats(out), res, err
}

type branch struct {
	ID              string `json:"id"`
	DisplayID       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestChangeset"`
	IsDefault       bool   `json:"isDefault"`
}

type commits struct {
	pagination
	Values []*commit `json:"values"`
}

type branches struct {
	pagination
	Values []*branch `json:"values"`
}

type diffstats struct {
	pagination
	Values []*diffstat
}

type diffstat struct {
	ContentID     string `json:"contentId"`
	FromContentID string `json:"fromContentId"`
	Path          struct {
		Components []string `json:"components"`
		Parent     string   `json:"parent"`
		Name       string   `json:"name"`
		Extension  string   `json:"extension"`
		ToString   string   `json:"toString"`
	} `json:"path"`
	PercentUnchanged int    `json:"percentUnchanged"`
	Type             string `json:"type"`
	NodeType         string `json:"nodeType"`
	SrcExecutable    bool   `json:"srcExecutable"`
	Links            struct {
		Self []struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
	Properties struct {
		GitChangeType string `json:"gitChangeType"`
	} `json:"properties"`
}

type commit struct {
	ID        string `json:"id"`
	DisplayID string `json:"displayId"`
	Author    struct {
		Name         string `json:"name"`
		EmailAddress string `json:"emailAddress"`
		ID           int    `json:"id"`
		DisplayName  string `json:"displayName"`
		Active       bool   `json:"active"`
		Slug         string `json:"slug"`
		Type         string `json:"type"`
		Links        struct {
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"author"`
	AuthorTimestamp int64 `json:"authorTimestamp"`
	Committer       struct {
		Name         string `json:"name"`
		EmailAddress string `json:"emailAddress"`
		ID           int    `json:"id"`
		DisplayName  string `json:"displayName"`
		Active       bool   `json:"active"`
		Slug         string `json:"slug"`
		Type         string `json:"type"`
		Links        struct {
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"committer"`
	CommitterTimestamp int64  `json:"committerTimestamp"`
	Message            string `json:"message"`
	Parents            []struct {
		ID        string `json:"id"`
		DisplayID string `json:"displayId"`
		Author    struct {
			Name         string `json:"name"`
			EmailAddress string `json:"emailAddress"`
		} `json:"author"`
		AuthorTimestamp int64 `json:"authorTimestamp"`
		Committer       struct {
			Name         string `json:"name"`
			EmailAddress string `json:"emailAddress"`
		} `json:"committer"`
		CommitterTimestamp int64  `json:"committerTimestamp"`
		Message            string `json:"message"`
		Parents            []struct {
			ID        string `json:"id"`
			DisplayID string `json:"displayId"`
		} `json:"parents"`
	} `json:"parents"`
}

func convertDiffstats(from *diffstats) []*scm.Change {
	to := []*scm.Change{}
	for _, v := range from.Values {
		to = append(to, convertDiffstat(v))
	}
	return to
}

func convertDiffstat(from *diffstat) *scm.Change {
	return &scm.Change{
		Path:    from.Path.ToString,
		Added:   from.Type == "ADD",
		Renamed: from.Type == "MOVE",
		Deleted: from.Type == "DELETE",
	}
}

func convertCommitList(from *commits) []*scm.Commit {
	to := []*scm.Commit{}
	for _, v := range from.Values {
		to = append(to, convertCommit(v))
	}
	return to
}

func convertCommit(from *commit) *scm.Commit {
	return &scm.Commit{
		Message: from.Message,
		Sha:     from.ID,
		// Link:    "%s/projects/%s/repos/%s/commits/%s",
		Author: scm.Signature{
			Name:   from.Author.DisplayName,
			Email:  from.Author.EmailAddress,
			Date:   time.Unix(from.AuthorTimestamp/1000, 0),
			Login:  from.Author.Slug,
			Avatar: avatarLink(from.Author.EmailAddress),
		},
		Committer: scm.Signature{
			Name:   from.Committer.DisplayName,
			Email:  from.Committer.EmailAddress,
			Date:   time.Unix(from.CommitterTimestamp/1000, 0),
			Login:  from.Committer.Slug,
			Avatar: avatarLink(from.Committer.EmailAddress),
		},
	}
}

func convertBranchList(from *branches) []*scm.Reference {
	to := []*scm.Reference{}
	for _, v := range from.Values {
		to = append(to, convertBranch(v))
	}
	return to
}

func convertBranch(from *branch) *scm.Reference {
	return &scm.Reference{
		Name: scm.TrimRef(from.DisplayID),
		Path: scm.ExpandRef(from.DisplayID, "refs/heads/"),
		Sha:  from.LatestCommit,
	}
}

func convertTagList(from *branches) []*scm.Reference {
	to := []*scm.Reference{}
	for _, v := range from.Values {
		to = append(to, convertTag(v))
	}
	return to
}

func convertTag(from *branch) *scm.Reference {
	return &scm.Reference{
		Name: scm.TrimRef(from.DisplayID),
		Path: scm.ExpandRef(from.DisplayID, "refs/tags/"),
		Sha:  from.LatestCommit,
	}
}
