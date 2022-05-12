package integration

import (
	"context"
	"net/http"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/stash"
	"github.com/drone/go-scm/scm/transport"
)

func TestCreateBranch(t *testing.T) {
	if token == "" {
		t.Skip("Skipping, Acceptance test")
	}
	client, _ = stash.New(endpoint)
	client.Client = &http.Client{
		Transport: &transport.BasicAuth{
			Username: username,
			Password: token,
		},
	}

	commitId, _ := GetCurrentCommitOfBranch(client, "master")
	input := &scm.CreateBranch{
		Name: "test_branch",
		Sha:  commitId,
	}
	response, listerr := client.Git.CreateBranch(context.Background(), repoID, input)
	if listerr != nil {
		t.Errorf("CreateBranch got an error %v", listerr)
	}
	if response.Status != http.StatusOK {
		t.Errorf("CreateBranch did not get a 200 back %v", response.Status)
	}

}
