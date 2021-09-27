package gitea

import (
	"context"
	"fmt"
	"time"

	"github.com/drone/go-scm/scm"
)

type milestoneService struct {
	client *wrapper
}

func (s *milestoneService) Find(ctx context.Context, repo string, id int) (*scm.Milestone, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("api/v1/repos/%s/%s/milestones/%d", namespace, name, id)
	out := new(milestone)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertMilestone(out), res, err
}

func (s *milestoneService) List(ctx context.Context, repo string, opts scm.MilestoneListOptions) ([]*scm.Milestone, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("api/v1/repos/%s/%s/milestones%s", namespace, name, encodeMilestoneListOptions(opts))
	out := []*milestone{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertMilestoneList(out), res, err
}

func (s *milestoneService) Create(ctx context.Context, repo string, input *scm.MilestoneInput) (*scm.Milestone, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("api/v1/repos/%s/%s/milestones", namespace, name)
	in := &milestoneInput{
		Title:       input.Title,
		Description: input.Description,
		State:       StateOpen,
		Deadline:    input.DueDate,
	}
	if input.State == "closed" {
		in.State = StateClosed
	}
	out := new(milestone)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertMilestone(out), res, err
}

func (s *milestoneService) Delete(ctx context.Context, repo string, id int) (*scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("api/v1/repos/%s/%s/milestones/%d", namespace, name, id)
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

func (s *milestoneService) Update(ctx context.Context, repo string, id int, input *scm.MilestoneInput) (*scm.Milestone, *scm.Response, error) {
	namespace, name := scm.Split(repo)
	path := fmt.Sprintf("api/v1/repos/%s/%s/milestones/%d", namespace, name, id)
	in := milestoneInput{}
	if input.Title != "" {
		in.Title = input.Title
	}
	switch input.State {
	case "open":
		in.State = StateOpen
	case "close", "closed":
		in.State = StateClosed
	}
	if input.Description != "" {
		in.Description = input.Description
	}
	if !input.DueDate.IsZero() {
		in.Deadline = input.DueDate
	}
	out := new(milestone)
	res, err := s.client.do(ctx, "PATCH", path, in, out)
	return convertMilestone(out), res, err
}

// StateType issue state type
type StateType string

const (
	// StateOpen pr/issue is open
	StateOpen StateType = "open"
	// StateClosed pr/issue is closed
	StateClosed StateType = "closed"
	// StateAll is all
	StateAll StateType = "all"
)

type milestone struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	State        StateType `json:"state"`
	OpenIssues   int       `json:"open_issues"`
	ClosedIssues int       `json:"closed_issues"`
	Created      time.Time `json:"created_at"`
	Updated      time.Time `json:"updated_at"`
	Closed       time.Time `json:"closed_at"`
	Deadline     time.Time `json:"due_on"`
}

type milestoneInput struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	State       StateType `json:"state"`
	Deadline    time.Time `json:"due_on"`
}

func convertMilestoneList(src []*milestone) []*scm.Milestone {
	var dst []*scm.Milestone
	for _, v := range src {
		dst = append(dst, convertMilestone(v))
	}
	return dst
}

func convertMilestone(src *milestone) *scm.Milestone {
	if src == nil || src.Deadline.IsZero() {
		return nil
	}
	return &scm.Milestone{
		Number:      int(src.ID),
		ID:          int(src.ID),
		Title:       src.Title,
		Description: src.Description,
		State:       string(src.State),
		DueDate:     src.Deadline,
	}
}
