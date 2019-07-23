package factory

import (
	"context"
	"fmt"
	"strings"

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
	case "bitbucket", "bitbucketcloud":
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
			client, err = github.New(ensureGHEEndpoint(serverURL))
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
	case "stash", "bitbucketserver":
		if serverURL == "" {
			return nil, MissingGitServerURL
		}
		client, err = stash.New(ensureBitBucketServerEndpoint(serverURL))
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

func ensureBitBucketServerEndpoint(s string) string {
	if strings.HasSuffix(s, "/rest") || strings.HasSuffix(s, "/rest/") {
		return s
	}
	return scm.UrlJoin(s, "/rest")
}

// ensureGHEEndpoint lets ensure we have the /api/v3 suffix on the URL
func ensureGHEEndpoint(u string) string {
	if strings.HasPrefix(u, "https://github.com") || strings.HasPrefix(u, "http://github.com") {
		return u
	}
	// lets ensure we use the API endpoint to login
	if strings.Index(u, "/api/") < 0 {
		u = scm.UrlJoin(u, "/api/v3")
	}
	return u
}
