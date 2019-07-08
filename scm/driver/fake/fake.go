package fake

import (
	"net/url"

	"github.com/jenkins-x/go-scm/scm"
)

type fakeClient struct {
}

// NewDefault returns a new GitHub API client using the
// default api.github.com address.
func NewDefault() (*scm.Client, *Data) {
	data := NewData()
	client := &wrapper{new(scm.Client)}
	client.BaseURL = &url.URL{
		Host: "fake.com",
		Path: "/",
	}
	// initialize services
	client.Driver = scm.DriverFake

	client.Repositories = &repositoryService{client: client, data: data}
	client.Issues = &issueService{client: client, data: data}

	// TODO
	/*
		client.Contents = &contentService{client}
		client.Git = &gitService{client}
		client.Organizations = &organizationService{client}
		client.PullRequests = &pullService{&issueService{client}}
		client.Reviews = &reviewService{client}
		client.Users = &userService{client}
		client.Webhooks = &webhookService{client}
	*/
	return client.Client, data
}

// wraper wraps the Client to provide high level helper functions
// for making http requests and unmarshaling the response.
type wrapper struct {
	*scm.Client
}
