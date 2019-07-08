package fake

import (
	"context"

	"github.com/jenkins-x/go-scm/scm"
)

type pullService struct {
	client *wrapper
	data   *Data
}

func (s *pullService) Find(context.Context, string, int) (*scm.PullRequest, *scm.Response, error) {
	panic("implement me")
}

func (s *pullService) FindComment(context.Context, string, int, int) (*scm.Comment, *scm.Response, error) {
	panic("implement me")
}

func (s *pullService) List(context.Context, string, scm.PullRequestListOptions) ([]*scm.PullRequest, *scm.Response, error) {
	panic("implement me")
}

func (s *pullService) ListChanges(ctx context.Context, repo string, number int, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	f := s.data
	return f.PullRequestChanges[number], nil, nil
}

func (s *pullService) ListComments(ctx context.Context, repo string, number int, opts scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	f := s.data
	return append([]*scm.Comment{}, f.PullRequestComments[number]...), nil, nil
}

func (s *pullService) Merge(context.Context, string, int) (*scm.Response, error) {
	panic("implement me")
}

func (s *pullService) Close(context.Context, string, int) (*scm.Response, error) {
	panic("implement me")
}

func (s *pullService) CreateComment(context.Context, string, int, *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	panic("implement me")
}

func (s *pullService) DeleteComment(context.Context, string, int, int) (*scm.Response, error) {
	panic("implement me")
}
