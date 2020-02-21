package github

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jenkins-x/go-scm/scm"
)

type deploymentService struct {
	client *wrapper
}

type deployment struct {
	Namespace             string
	Name                  string
	FullName              string
	ID                    int       `json:"id"`
	Link                  string    `json:"url"`
	Sha                   string    `json:"sha"`
	Ref                   string    `json:"ref"`
	Description           string    `json:"description"`
	OriginalEnvironment   string    `json:"original_environment"`
	Environment           string    `json:"environment"`
	RepositoryLink        string    `json:"repository_url"`
	StatusLink            string    `json:"statuses_url"`
	Author                *user     `json:"creator"`
	Created               time.Time `json:"created_at"`
	Updated               time.Time `json:"updated_at"`
	TransientEnvironment  bool      `json:"transient_environment"`
	ProductionEnvironment bool      `json:"production_environment"`
}

type deploymentInput struct {
	Ref                   string   `json:"ref"`
	Task                  string   `json:"task"`
	Payload               string   `json:"payload"`
	Environment           string   `json:"environment"`
	Description           string   `json:"description"`
	RequiredContexts      []string `json:"required_contexts"`
	AutoMerge             bool     `json:"auto_merge"`
	TransientEnvironment  bool     `json:"transient_environment"`
	ProductionEnvironment bool     `json:"production_environment"`
}

func (s *deploymentService) Find(ctx context.Context, repoFullName string, deploymentID string) (*scm.Deployment, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/deployments/%s?", repoFullName, deploymentID)
	out := new(deployment)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertDeployment(out, repoFullName), res, err
}

func (s *deploymentService) List(ctx context.Context, repoFullName string, opts scm.ListOptions) ([]*scm.Deployment, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/deployments?%s", repoFullName, encodeListOptions(opts))
	out := []*deployment{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertDeploymentList(out, repoFullName), res, err
}

func (s *deploymentService) Create(ctx context.Context, repoFullName string, deploymentInput *scm.DeploymentInput) (*scm.Deployment, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/deployments", repoFullName)
	in := convertToDeploymentInput(deploymentInput)
	out := new(deployment)
	res, err := s.client.do(ctx, "POST", path, in, out)
	return convertDeployment(out, repoFullName), res, err
}

func (s *deploymentService) Delete(ctx context.Context, repoFullName string, deploymentID string) (*scm.Response, error) {
	panic("implement me")
}

func (s *deploymentService) FindStatus(ctx context.Context, repoFullName string, deploymentID string, statusID string) (*scm.DeploymentStatus, *scm.Response, error) {
	panic("implement me")
}

func (s *deploymentService) ListStatus(ctx context.Context, repoFullName string, options scm.ListOptions) ([]*scm.DeploymentStatus, *scm.Response, error) {
	panic("implement me")
}

func (s *deploymentService) CreateStatus(ctx context.Context, repoFullName string, deployment *scm.DeploymentStatus) (*scm.DeploymentStatus, *scm.Response, error) {
	panic("implement me")
}

func (s *deploymentService) UpdateStatus(ctx context.Context, repoFullName string, deployment *scm.DeploymentStatus) (*scm.DeploymentStatus, *scm.Response, error) {
	panic("implement me")
}

func (s *deploymentService) DeleteStatus(ctx context.Context, repoFullName string, deploymentID string, statusID string) (*scm.Response, error) {
	panic("implement me")
}

func convertDeploymentList(out []*deployment, fullName string) []*scm.Deployment {
	answer := []*scm.Deployment{}
	for _, o := range out {
		answer = append(answer, convertDeployment(o, fullName))
	}
	return answer
}

func convertToDeploymentInput(from *scm.DeploymentInput) *deploymentInput {
	return &deploymentInput{
		Ref:                   from.Ref,
		Task:                  from.Task,
		Payload:               from.Payload,
		Environment:           from.Environment,
		Description:           from.Description,
		RequiredContexts:      from.RequiredContexts,
		AutoMerge:             from.AutoMerge,
		TransientEnvironment:  from.TransientEnvironment,
		ProductionEnvironment: from.ProductionEnvironment,
	}
}

func convertDeployment(from *deployment, fullName string) *scm.Deployment {
	dst := &scm.Deployment{
		ID:                    strconv.Itoa(from.ID),
		Link:                  from.Link,
		Sha:                   from.Sha,
		Ref:                   from.Ref,
		FullName:              fullName,
		Description:           from.Description,
		OriginalEnvironment:   from.OriginalEnvironment,
		Environment:           from.Environment,
		RepositoryLink:        from.RepositoryLink,
		StatusLink:            from.StatusLink,
		Author:                convertUser(from.Author),
		Created:               from.Created,
		Updated:               from.Updated,
		TransientEnvironment:  from.TransientEnvironment,
		ProductionEnvironment: from.ProductionEnvironment,
	}
	names := strings.Split(fullName, "/")
	if len(names) > 1 {
		dst.Namespace = names[0]
		dst.Name = names[1]
	}
	return dst
}
