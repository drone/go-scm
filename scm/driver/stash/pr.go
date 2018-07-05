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

type pullService struct {
	client *wrapper
}

func (s *pullService) Find(ctx context.Context, repo string, number int) (*scm.PullRequest, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/pull-requests/%d", namespace, name, number)
	out := new(pullRequest)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertPullRequest(out), res, err
}

func (s *pullService) FindComment(ctx context.Context, repo string, number int, id int) (*scm.Comment, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/pull-requests/%d/comments/%d", namespace, name, number, id)
	out := new(pullRequestComment)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertPullRequestComment(out), res, err
}

func (s *pullService) List(ctx context.Context, repo string, opts scm.PullRequestListOptions) ([]*scm.PullRequest, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/pull-requests", namespace, name)
	out := new(pullRequests)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	copyPagination(out.pagination, res)
	return convertPullRequests(out), res, err
}

func (s *pullService) ListChanges(ctx context.Context, repo string, number int, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/pull-requests/%d/changes", namespace, name, number)
	out := new(diffstats)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	copyPagination(out.pagination, res)
	return convertDiffstats(out), res, err
}

func (s *pullService) ListComments(context.Context, string, int, scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	// rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/activities
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) Merge(ctx context.Context, repo string, number int) (*scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/pull-requests/%d/merge", namespace, name, number)
	res, err := s.client.do(ctx, "POST", path, nil, nil)
	return res, err
}

func (s *pullService) Close(ctx context.Context, repo string, number int) (*scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/pull-requests/%d/decline", namespace, name, number)
	res, err := s.client.do(ctx, "POST", path, nil, nil)
	return res, err
}

func (s *pullService) CreateComment(context.Context, string, int, *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *pullService) DeleteComment(context.Context, string, int, int) (*scm.Response, error) {
	// {"errors":[{"context":null,"message":"You are attempting to modify a comment based on out-of-date information.","exceptionName":"com.atlassian.bitbucket.comment.CommentOutOfDateException","currentVersion":1,"expectedVersion":0}]}
	// rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/comments/3?version=0
	return nil, scm.ErrNotSupported
}

type pullRequest struct {
	ID          int    `json:"id"`
	Version     int    `json:"version"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"`
	Open        bool   `json:"open"`
	Closed      bool   `json:"closed"`
	CreatedDate int64  `json:"createdDate"`
	UpdatedDate int64  `json:"updatedDate"`
	FromRef     struct {
		ID           string `json:"id"`
		DisplayID    string `json:"displayId"`
		LatestCommit string `json:"latestCommit"`
		Repository   struct {
			Slug          string `json:"slug"`
			ID            int    `json:"id"`
			Name          string `json:"name"`
			ScmID         string `json:"scmId"`
			State         string `json:"state"`
			StatusMessage string `json:"statusMessage"`
			Forkable      bool   `json:"forkable"`
			Project       struct {
				Key    string `json:"key"`
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Public bool   `json:"public"`
				Type   string `json:"type"`
				Links  struct {
					Self []struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"links"`
			} `json:"project"`
			Public bool `json:"public"`
			Links  struct {
				Clone []struct {
					Href string `json:"href"`
					Name string `json:"name"`
				} `json:"clone"`
				Self []struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"repository"`
	} `json:"fromRef"`
	ToRef struct {
		ID           string `json:"id"`
		DisplayID    string `json:"displayId"`
		LatestCommit string `json:"latestCommit"`
		Repository   struct {
			Slug          string `json:"slug"`
			ID            int    `json:"id"`
			Name          string `json:"name"`
			ScmID         string `json:"scmId"`
			State         string `json:"state"`
			StatusMessage string `json:"statusMessage"`
			Forkable      bool   `json:"forkable"`
			Project       struct {
				Key    string `json:"key"`
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Public bool   `json:"public"`
				Type   string `json:"type"`
				Links  struct {
					Self []struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"links"`
			} `json:"project"`
			Public bool `json:"public"`
			Links  struct {
				Clone []struct {
					Href string `json:"href"`
					Name string `json:"name"`
				} `json:"clone"`
				Self []struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"repository"`
	} `json:"toRef"`
	Locked bool `json:"locked"`
	Author struct {
		User struct {
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
		} `json:"user"`
		Role     string `json:"role"`
		Approved bool   `json:"approved"`
		Status   string `json:"status"`
	} `json:"author"`
	Reviewers    []interface{} `json:"reviewers"`
	Participants []interface{} `json:"participants"`
	Links        struct {
		Self []link `json:"self"`
	} `json:"links"`
}

type pullRequests struct {
	pagination
	Values []*pullRequest `json:"values"`
}

func convertPullRequests(from *pullRequests) []*scm.PullRequest {
	to := []*scm.PullRequest{}
	for _, v := range from.Values {
		to = append(to, convertPullRequest(v))
	}
	return to
}

func convertPullRequest(from *pullRequest) *scm.PullRequest {
	return &scm.PullRequest{
		Number:  from.ID,
		Title:   from.Title,
		Body:    from.Description,
		Sha:     from.FromRef.LatestCommit,
		Ref:     fmt.Sprintf("refs/pull-requests/%d/from", from.ID),
		Source:  from.FromRef.DisplayID,
		Target:  from.ToRef.DisplayID,
		Link:    extractSelfLink(from.Links.Self),
		Closed:  from.Closed,
		Merged:  from.State == "MERGED",
		Created: time.Unix(from.CreatedDate/1000, 0),
		Updated: time.Unix(from.UpdatedDate/1000, 0),
		Author: scm.User{
			Login:  from.Author.User.Slug,
			Name:   from.Author.User.DisplayName,
			Email:  from.Author.User.EmailAddress,
			Avatar: avatarLink(from.Author.User.EmailAddress),
		},
	}
}

type pullRequestComment struct {
	Properties struct {
		RepositoryID int `json:"repositoryId"`
	} `json:"properties"`
	ID      int    `json:"id"`
	Version int    `json:"version"`
	Text    string `json:"text"`
	Author  struct {
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
	CreatedDate         int64         `json:"createdDate"`
	UpdatedDate         int64         `json:"updatedDate"`
	Comments            []interface{} `json:"comments"`
	Tasks               []interface{} `json:"tasks"`
	PermittedOperations struct {
		Editable  bool `json:"editable"`
		Deletable bool `json:"deletable"`
	} `json:"permittedOperations"`
}

func convertPullRequestComment(from *pullRequestComment) *scm.Comment {
	return &scm.Comment{
		ID:      from.ID,
		Body:    from.Text,
		Created: time.Unix(from.CreatedDate/1000, 0),
		Updated: time.Unix(from.UpdatedDate/1000, 0),
		Author: scm.User{
			Login:  from.Author.Slug,
			Name:   from.Author.DisplayName,
			Email:  from.Author.EmailAddress,
			Avatar: avatarLink(from.Author.EmailAddress),
		},
	}
}
