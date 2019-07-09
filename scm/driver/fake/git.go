package fake

import (
	"context"

	"github.com/jenkins-x/go-scm/scm"
)

type gitService struct {
	client *wrapper
	data   *Data
}

func (s *gitService) FindBranch(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	panic("implement me")
}

func (s *gitService) FindCommit(ctx context.Context, repo, SHA string) (*scm.Commit, *scm.Response, error) {
	f := s.data
	return f.Commits[SHA], nil, nil
}

func (s *gitService) FindTag(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	panic("implement me")
}

func (s *gitService) ListBranches(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	panic("implement me")
}

func (s *gitService) ListCommits(ctx context.Context, repo string, opts scm.CommitListOptions) ([]*scm.Commit, *scm.Response, error) {
	panic("implement me")
}

func (s *gitService) ListChanges(ctx context.Context, repo, ref string, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	panic("implement me")
}

func (s *gitService) ListTags(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	panic("implement me")
}
