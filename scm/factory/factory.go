package factory

import (
	"context"
	"fmt"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/driver/bitbucket"
	"github.com/jenkins-x/go-scm/scm/driver/fake"
	"github.com/jenkins-x/go-scm/scm/driver/gitea"
	"github.com/jenkins-x/go-scm/scm/driver/github"
	"github.com/jenkins-x/go-scm/scm/driver/gitlab"
	"github.com/jenkins-x/go-scm/scm/driver/gogs"
	"github.com/jenkins-x/go-scm/scm/driver/stash"
	"golang.org/x/oauth2"
)

// MissingGitServerURL the error returned if you use a git driver that needs a git server URL
var MissingGitServerURL = fmt.Errorf("No git serverURL was specified")

func NewClient(driver, serverURL, oauthToken string) (*scm.Client, error) {
	if driver == "" {
		driver = "github"
	}
	var client *scm.Client
	var err error

	switch driver {
	case "bitbucket":
		if serverURL != "" {
			client, err = bitbucket.New(serverURL)
		} else {
			client = bitbucket.NewDefault()
		}
	case "fake":
		client, _ = fake.NewDefault()
	case "gitea":
		if serverURL == "" {
			return nil, MissingGitServerURL
		}
		client, err = gitea.New(serverURL)
	case "github":
		if serverURL != "" {
			client, err = github.New(serverURL)
		} else {
			client = github.NewDefault()
		}
	case "gitlab":
		if serverURL != "" {
			client, err = gitlab.New(serverURL)
		} else {
			client = gitlab.NewDefault()
		}
	case "gogs":
		if serverURL == "" {
			return nil, MissingGitServerURL
		}
		client, err = gogs.New(serverURL)
	case "stash":
		if serverURL == "" {
			return nil, MissingGitServerURL
		}
		client, err = stash.New(serverURL)
	default:
		return nil, fmt.Errorf("Unsupported $GIT_KIND value: %s", driver)
	}
	if err != nil {
		return client, err
	}
	if oauthToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: oauthToken},
		)
		client.Client = oauth2.NewClient(context.Background(), ts)
	}
	return client, err
}
