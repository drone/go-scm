// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"time"

	"github.com/jenkins-x/go-scm/scm"
)

type stateType string

const (
	// stateOpen pr/issue is opend
	stateOpen stateType = "open"
	// stateClosed pr/issue is closed
	stateClosed stateType = "closed"
	// stateAll is all
	stateAll stateType = "all"
)

type issueService struct {
	client *wrapper
}

func (s *issueService) Search(context.Context, scm.SearchOptions) ([]*scm.SearchIssue, *scm.Response, error) {
	// TODO implemment
	return nil, nil, scm.ErrNotSupported
}

func (s *issueService) AssignIssue(ctx context.Context, repo string, number int, logins []string) (*scm.Response, error) {
	issue, res, err := s.Find(ctx, repo, number)
	if err != nil {
		return res, errors.Wrapf(err, "couldn't lookup issue %d in repository %s", number, repo)
	}
	if issue == nil {
		return res, fmt.Errorf("couldn't find issue %d in repository %s", number, repo)
	}
	assignees := sets.NewString(logins...)
	for _, existingAssignee := range issue.Assignees {
		assignees.Insert(existingAssignee.Login)
	}

	namespace, name := scm.Split(repo)
	in := gitea.EditIssueOption{
		Title:     issue.Title,
		Assignees: assignees.List(),
	}
	_, err = s.client.GiteaClient.EditIssue(namespace, name, int64(number), in)
	return nil, err
}

func (s *issueService) UnassignIssue(ctx context.Context, repo string, number int, logins []string) (*scm.Response, error) {
	issue, res, err := s.Find(ctx, repo, number)
	if err != nil {
		return res, errors.Wrapf(err, "couldn't lookup issue %d in repository %s", number, repo)
	}
	if issue == nil {
		return res, fmt.Errorf("couldn't find issue %d in repository %s", number, repo)
	}
	assignees := sets.NewString()
	for _, existingAssignee := range issue.Assignees {
		assignees.Insert(existingAssignee.Login)
	}
	assignees.Delete(logins...)

	namespace, name := scm.Split(repo)
	in := gitea.EditIssueOption{
		Title:     issue.Title,
		Assignees: assignees.List(),
	}
	_, err = s.client.GiteaClient.EditIssue(namespace, name, int64(number), in)
	return nil, err
}

func (s *issueService) ListEvents(context.Context, string, int, scm.ListOptions) ([]*scm.ListedIssueEvent, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *issueService) ListLabels(ctx context.Context, repo string, number int, _ scm.ListOptions) ([]*scm.Label, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	out, err := s.client.GiteaClient.GetIssueLabels(namespace, name, int64(number), gitea.ListLabelsOptions{})
	return convertGiteaLabels(out), nil, err
}

func (s *issueService) lookupLabel(ctx context.Context, repo string, lbl string) (int64, *scm.Response, error) {
	var labelID int64
	labelID = -1
	repoLabels, res, err := s.client.Repositories.ListLabels(ctx, repo, scm.ListOptions{})
	if err != nil {
		return labelID, res, errors.Wrapf(err, "listing labels in repository %s", repo)
	}
	for _, l := range repoLabels {
		if l.Name == lbl {
			labelID = l.ID
			break
		}
	}
	return labelID, res, nil
}

func (s *issueService) AddLabel(ctx context.Context, repo string, number int, lbl string) (*scm.Response, error) {
	labelID, res, err := s.lookupLabel(ctx, repo, lbl)
	if err != nil {
		return res, err
	}
	namespace, name := scm.Split(repo)

	if labelID == -1 {
		lblInput := gitea.CreateLabelOption{
			Color:       "#00aabb",
			Description: "",
			Name:        lbl,
		}
		newLabel, err := s.client.GiteaClient.CreateLabel(namespace, name, lblInput)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create label %s in repository %s", lbl, repo)
		}
		labelID = newLabel.ID
	}

	in := gitea.IssueLabelsOption{Labels: []int64{labelID}}
	_, err = s.client.GiteaClient.AddIssueLabels(namespace, name, int64(number), in)
	return nil, err
}

func (s *issueService) DeleteLabel(ctx context.Context, repo string, number int, lbl string) (*scm.Response, error) {
	labelID, res, err := s.lookupLabel(ctx, repo, lbl)
	if err != nil {
		return res, err
	}
	if labelID == -1 {
		return nil, nil
	}

	namespace, name := scm.Split(repo)
	err = s.client.GiteaClient.DeleteIssueLabel(namespace, name, int64(number), labelID)
	return nil, err
}

func (s *issueService) Find(ctx context.Context, repo string, number int) (*scm.Issue, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	out, err := s.client.GiteaClient.GetIssue(namespace, name, int64(number))
	return convertGiteaIssue(out), nil, err
}

func (s *issueService) FindComment(ctx context.Context, repo string, index, id int) (*scm.Comment, *scm.Response, error) {
	comments, res, err := s.ListComments(ctx, repo, index, scm.ListOptions{})
	if err != nil {
		return nil, res, err
	}
	for _, comment := range comments {
		if comment.ID == id {
			return comment, res, nil
		}
	}
	return nil, res, nil
}

func (s *issueService) List(ctx context.Context, repo string, _ scm.IssueListOptions) ([]*scm.Issue, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	in := gitea.ListIssueOption{
		Type: gitea.IssueTypeIssue,
	}
	out, err := s.client.GiteaClient.ListRepoIssues(namespace, name, in)
	return convertIssueList(out), nil, err
}

func (s *issueService) ListComments(ctx context.Context, repo string, index int, _ scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	out, err := s.client.GiteaClient.ListIssueComments(namespace, name, int64(index), gitea.ListIssueCommentOptions{})
	return convertIssueCommentList(out), nil, err
}

func (s *issueService) Create(ctx context.Context, repo string, input *scm.IssueInput) (*scm.Issue, *scm.Response, error) {
	namespace, name := scm.Split(repo)

	in := gitea.CreateIssueOption{
		Title: input.Title,
		Body:  input.Body,
	}
	out, err := s.client.GiteaClient.CreateIssue(namespace, name, in)
	return convertGiteaIssue(out), nil, err
}

func (s *issueService) CreateComment(ctx context.Context, repo string, index int, input *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	in := gitea.CreateIssueCommentOption{Body: input.Body}
	out, err := s.client.GiteaClient.CreateIssueComment(namespace, name, int64(index), in)
	return convertGiteaIssueComment(out), nil, err
}

func (s *issueService) DeleteComment(ctx context.Context, repo string, index, id int) (*scm.Response, error) {
	namespace, name := scm.Split(repo)
	return nil, s.client.GiteaClient.DeleteIssueComment(namespace, name, int64(id))
}

func (s *issueService) EditComment(ctx context.Context, repo string, number int, id int, input *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	in := gitea.EditIssueCommentOption{Body: input.Body}
	out, err := s.client.GiteaClient.EditIssueComment(namespace, name, int64(id), in)
	return convertGiteaIssueComment(out), nil, err
}

func (s *issueService) Close(ctx context.Context, repo string, number int) (*scm.Response, error) {
	namespace, name := scm.Split(repo)
	closed := gitea.StateClosed
	in := gitea.EditIssueOption{
		State: &closed,
	}
	_, err := s.client.GiteaClient.EditIssue(namespace, name, int64(number), in)
	return nil, err
}

func (s *issueService) Reopen(ctx context.Context, repo string, number int) (*scm.Response, error) {
	namespace, name := scm.Split(repo)
	reopen := gitea.StateOpen
	in := gitea.EditIssueOption{
		State: &reopen,
	}
	_, err := s.client.GiteaClient.EditIssue(namespace, name, int64(number), in)
	return nil, err
}

func (s *issueService) Lock(ctx context.Context, repo string, number int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *issueService) Unlock(ctx context.Context, repo string, number int) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

//
// native data structures
//

type label struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	// example: 00aabb
	Color       string `json:"color"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type closeReopenInput struct {
	State stateType `json:"state"`
}

type (
	// gitea issue response object.
	issue struct {
		ID          int       `json:"id"`
		Number      int       `json:"number"`
		User        user      `json:"user"`
		Title       string    `json:"title"`
		Body        string    `json:"body"`
		State       stateType `json:"state"`
		Labels      []label   `json:"labels"`
		Comments    int       `json:"comments"`
		Assignees   []user    `json:"assignees"`
		Created     time.Time `json:"created_at"`
		Updated     time.Time `json:"updated_at"`
		PullRequest *struct {
			Merged   bool        `json:"merged"`
			MergedAt interface{} `json:"merged_at"`
		} `json:"pull_request"`
	}

	// gitea issue comment response object.
	issueComment struct {
		ID        int       `json:"id"`
		HTMLURL   string    `json:"html_url"`
		User      user      `json:"user"`
		Body      string    `json:"body"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

//
// native data structure conversion
//

func convertIssue(from *issue) *scm.Issue {
	return &scm.Issue{
		Number:    from.Number,
		Title:     from.Title,
		Body:      from.Body,
		Link:      "", // TODO construct the link to the issue.
		Closed:    from.State == "closed",
		Labels:    convertLabels(from),
		Author:    *convertUser(&from.User),
		Assignees: convertUsers(from.Assignees),
		Created:   from.Created,
		Updated:   from.Updated,
	}
}

func convertIssueList(from []*gitea.Issue) []*scm.Issue {
	to := []*scm.Issue{}
	for _, v := range from {
		to = append(to, convertGiteaIssue(v))
	}
	return to
}

func convertGiteaIssue(from *gitea.Issue) *scm.Issue {
	return &scm.Issue{
		Number:    int(from.Index),
		Title:     from.Title,
		Body:      from.Body,
		Link:      from.URL,
		Closed:    from.State == gitea.StateClosed,
		Labels:    convertIssueLabels(from),
		Author:    *convertGiteaUser(from.Poster),
		Assignees: convertGiteaUsers(from.Assignees),
		Created:   from.Created,
		Updated:   from.Updated,
	}
}

func convertIssueComment(from *issueComment) *scm.Comment {
	return &scm.Comment{
		ID:      from.ID,
		Body:    from.Body,
		Author:  *convertUser(&from.User),
		Created: from.CreatedAt,
		Updated: from.UpdatedAt,
	}
}

func convertIssueCommentList(from []*gitea.Comment) []*scm.Comment {
	to := []*scm.Comment{}
	for _, v := range from {
		to = append(to, convertGiteaIssueComment(v))
	}
	return to
}

func convertGiteaIssueComment(from *gitea.Comment) *scm.Comment {
	if from == nil || from.Poster == nil {
		return nil
	}
	return &scm.Comment{
		ID:      int(from.ID),
		Body:    from.Body,
		Author:  *convertGiteaUser(from.Poster),
		Created: from.Created,
		Updated: from.Updated,
	}
}

func convertLabels(from *issue) []string {
	var labels []string
	for _, label := range from.Labels {
		labels = append(labels, label.Name)
	}
	return labels
}

func convertIssueLabels(from *gitea.Issue) []string {
	var labels []string
	for _, label := range from.Labels {
		labels = append(labels, label.Name)
	}
	return labels
}

func convertGiteaLabels(from []*gitea.Label) []*scm.Label {
	var labels []*scm.Label
	for _, label := range from {
		labels = append(labels, &scm.Label{
			Name:        label.Name,
			Description: label.Description,
			URL:         label.URL,
			Color:       label.Color,
		})
	}
	return labels
}
