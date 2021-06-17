package github

import (
	"context"
	"fmt"

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
