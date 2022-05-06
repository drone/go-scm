package integration

import (
	"os"

	"github.com/drone/go-scm/scm"
)

var (
	client   *scm.Client
	token    = os.Getenv("BITBUCKET_SERVER_TOKEN")
	endpoint = "https://bitbucket.dev.harness.io/"
	repoID   = "har/scm-integration-test-repo"
	username = "harnessadmin"
	commitId = "f675c4b55841908d7c338c500c8f4cb844fd9be7"
)
