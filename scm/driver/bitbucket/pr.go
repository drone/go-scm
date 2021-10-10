// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"fmt"
	"time"

	"github.com/livecycle/go-scm/scm"
)

type pullService struct {
	client *wrapper
}

func (s *pullService) Find(ctx context.Context, repo string, number int) (*scm.PullRequest, *scm.Response, error) {
	path := fmt.Sprintf("2.0/repositories/%s/pullrequests/%d", repo, number)
	out := new(pr)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertPullRequest(out), res, err
}

func (s *pullService) List(ctx context.Context, repo string, opts scm.PullRequestListOptions) ([]*scm.PullRequest, *scm.Response, error) {
	path := fmt.Sprintf("2.0/repositories/%s/pullrequests?%s", repo, encodePullRequestListOptions(opts))
	out := new(prs)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	copyPagination(out.pagination, res)
	return convertPullRequests(out), res, err
}

func (s *pullService) ListChanges(ctx context.Context, repo string, number int, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	path := fmt.Sprintf("2.0/repositories/%s/pullrequests/%d/diffstat?%s", repo, number, encodeListOptions(opts))
	out := new(diffstats)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	copyPagination(out.pagination, res)
	return convertDiffstats(out), res, err
}

func (s *pullService) Merge(ctx context.Context, repo string, number int) (*scm.Response, error) {
	path := fmt.Sprintf("2.0/repositories/%s/pullrequests/%d/merge", repo, number)
	res, err := s.client.do(ctx, "POST", path, nil, nil)
	return res, err
}

func (s *pullService) Close(ctx context.Context, repo string, number int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) Create(ctx context.Context, repo string, input *scm.PullRequestInput) (*scm.PullRequest, *scm.Response, error) {
	path := fmt.Sprintf("2.0/repositories/%s/pullrequests", repo)
	in := new(prInput)
	in.Title = input.Title
	in.Description = input.Body
	in.Source.Branch.Name = input.Source
	in.Destination.Branch.Name = input.Target
	out := new(pr)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertPullRequest(out), res, err
}

func (s *pullService) CreateComment(ctx context.Context, repo string, prID int, input *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	path := fmt.Sprintf("2.0/repositories/%s/pullrequests/%d/comments", repo, prID)
	in := new(commentIput)
	in.Content.Raw = input.Body
	out := new(commentOutput)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return converComment(out), res, err
}

func (s *pullService) EditComment(ctx context.Context, repo string, prID int, commentID int, input *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) DeleteComment(ctx context.Context, repo string, prID int, commentID int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *pullService) FindComment(ctx context.Context, repo string, prID int, commentID int) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) ListComments(ctx context.Context, repo string, prID int, options scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

type reference struct {
	Commit struct {
		Hash  string `json:"hash"`
		Links struct {
			Self link `json:"self"`
		} `json:"links"`
	} `json:"commit"`
	Branch struct {
		Name string `json:"name"`
	} `json:"branch"`
	Repository struct {
		FullName string `json:"full_name"`
		Type     string `json:"type"`
		Name     string `json:"name"`
		Links    struct {
			Self   link `json:"self"`
			HTML   link `json:"html"`
			Avatar link `json:"avatar"`
		} `json:"links"`
		UUID string `json:"uuid"`
	} `json:"repository"`
}

type pr struct {
	Description string `json:"description"`
	Links       struct {
		HTML link `json:"html"`
		Diff link `json:"diff"`
	} `json:"links"`
	Title        string    `json:"title"`
	ID           int       `json:"id"`
	Destination  reference `json:"destination"`
	CommentCount int       `json:"comment_count"`
	Summary      struct {
		Raw    string `json:"raw"`
		Markup string `json:"markup"`
		HTML   string `json:"html"`
		Type   string `json:"type"`
	} `json:"summary"`
	Source    reference `json:"source"`
	State     string    `json:"state"`
	Author    user      `json:"author"`
	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`
}

type prs struct {
	pagination
	Values []*pr `json:"values"`
}

type prInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Source      struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
	} `json:"source"`
	Destination struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
	} `json:"destination"`
}

type commentIput struct {
	Content struct {
		Raw string `json:"raw"`
	} `json:"content"`
}

type commentOutput struct {
	Content struct {
		Raw string `json:"raw"`
	} `json:"content"`
	CreatedOn time.Time `json:"created_on"`
	User      struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Type        string `json:"type"`
		UUID        string `json:"uuid"`
		Links       struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
	} `json:"user"`
	UpdatedOn time.Time `json:"updated_on"`
	ID        int       `json:"id"`
}

func convertPullRequests(from *prs) []*scm.PullRequest {
	to := []*scm.PullRequest{}
	for _, v := range from.Values {
		to = append(to, convertPullRequest(v))
	}
	return to
}

func convertPullRequest(from *pr) *scm.PullRequest {
	return &scm.PullRequest{
		Number: from.ID,
		Title:  from.Title,
		Body:   from.Description,
		Sha:    from.Source.Commit.Hash,
		Source: from.Source.Branch.Name,
		Target: from.Destination.Branch.Name,
		Fork:   from.Source.Repository.FullName,
		Link:   from.Links.HTML.Href,
		Diff:   from.Links.Diff.Href,
		Closed: from.State != "OPEN",
		Merged: from.State == "MERGED",
		Head: scm.Reference{
			Name: from.Source.Branch.Name,
			Path: scm.ExpandRef(from.Source.Branch.Name, "refs/heads"),
			Sha:  from.Source.Commit.Hash,
		},
		Base: scm.Reference{
			Name: from.Destination.Branch.Name,
			Path: scm.ExpandRef(from.Destination.Branch.Name, "refs/heads"),
			Sha:  from.Destination.Commit.Hash,
		},
		Author: scm.User{
			Login:  from.Author.Nickname,
			Name:   from.Author.DisplayName,
			Avatar: from.Author.Links.Avatar.Href,
		},
		Created: from.CreatedOn,
		Updated: from.UpdatedOn,
	}
}

func converComment(from *commentOutput) *scm.Comment {
	return &scm.Comment{
		ID:   from.ID,
		Body: from.Content.Raw,
		Author: scm.User{
			Login:  from.User.Username,
			Name:   from.User.DisplayName,
			Avatar: from.User.Links.Avatar.Href,
		},
		Created: from.CreatedOn,
		Updated: from.UpdatedOn,
	}
}
