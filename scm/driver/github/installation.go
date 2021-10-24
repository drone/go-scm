package github

import (
	"context"
	"fmt"
	"strconv"

	"github.com/livecycle/go-scm/scm"
)

type installationService struct {
	client *wrapper
}

func (s *installationService) ListRepositories(ctx context.Context, opts scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	path := fmt.Sprintf("installation/repositories?%s", encodeListOptions(opts))
	var out struct {
		Repositories []*repository `json:"repositories"`
	}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertRepositoryList(out.Repositories), res, err
}

// Installation
type installation struct {
	Id      int           `json:"id"`
	Account *organization `json:"account"`
	AppId   int           `json:"app_id"`
	AppSlug string        `json:"app_slug"`
}

func (s *installationService) ListInstallationsForAuthenticatedUser(ctx context.Context, opts scm.ListOptions) ([]*scm.Installation, *scm.Response, error) {
	path := fmt.Sprintf("user/installations?%s", encodeListOptions(opts))
	var out struct {
		Installations []*installation `json:"installations"`
	}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertInstallationList(out.Installations), res, err
}

// helper function to convert from the Github repository list to
// the common repository structure.
func convertInstallationList(from []*installation) []*scm.Installation {
	to := []*scm.Installation{}
	for _, v := range from {
		to = append(to, convertInstallation(v))
	}
	return to
}

// helper function to convert from the Github repository structure
// to the common repository structure.
func convertInstallation(from *installation) *scm.Installation {
	return &scm.Installation{
		ID:           strconv.Itoa(from.Id),
		Organization: convertOrganization(from.Account),
		AppName:      string(from.AppSlug),
	}
}
