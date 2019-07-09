package fake

import (
	"context"
	"strings"

	"github.com/jenkins-x/go-scm/scm"
)

type repositoryService struct {
	client *wrapper
	data   *Data
}

// NormLogin normalizes login strings
var NormLogin = strings.ToLower

func (s *repositoryService) IsCollaborator(ctx context.Context, repo, login string) (bool, *scm.Response, error) {
	f := s.data
	normed := NormLogin(login)
	for _, collab := range f.Collaborators {
		if NormLogin(collab) == normed {
			return true, nil, nil
		}
	}
	return false, nil, nil
}

func (s *repositoryService) ListCollaborators(ctx context.Context, repo string) ([]scm.User, *scm.Response, error) {
	f := s.data
	result := make([]scm.User, 0, len(f.Collaborators))
	for _, login := range f.Collaborators {
		result = append(result, scm.User{Login: login})
	}
	return result, nil, nil
}

func (s *repositoryService) Find(context.Context, string) (*scm.Repository, *scm.Response, error) {
	panic("implement me")
}

func (s *repositoryService) FindHook(context.Context, string, string) (*scm.Hook, *scm.Response, error) {
	panic("implement me")
}

func (s *repositoryService) FindPerms(context.Context, string) (*scm.Perm, *scm.Response, error) {
	panic("implement me")
}

func (s *repositoryService) List(context.Context, scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	panic("implement me")
}

func (s *repositoryService) ListLabels(context.Context, string, scm.ListOptions) ([]*scm.Label, *scm.Response, error) {
	f := s.data
	la := []*scm.Label{}
	for _, l := range f.RepoLabelsExisting {
		la = append(la, &scm.Label{Name: l})
	}
	return la, nil, nil
}

func (s *repositoryService) ListHooks(context.Context, string, scm.ListOptions) ([]*scm.Hook, *scm.Response, error) {
	panic("implement me")
}

func (s *repositoryService) ListStatus(context.Context, string, string, scm.ListOptions) ([]*scm.Status, *scm.Response, error) {
	panic("implement me")
}

func (s *repositoryService) CreateHook(context.Context, string, *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	panic("implement me")
}

func (s *repositoryService) CreateStatus(context.Context, string, string, *scm.StatusInput) (*scm.Status, *scm.Response, error) {
	panic("implement me")
}

func (s *repositoryService) DeleteHook(context.Context, string, string) (*scm.Response, error) {
	panic("implement me")
}
