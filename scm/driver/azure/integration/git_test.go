package integration

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/azure"
	"github.com/drone/go-scm/scm/transport"
)

func TestListBranches(t *testing.T) {
	if token == "" {
		t.Skip("Skipping, Acceptance test")
	}
	client = azure.NewDefault(organization, project)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
			},
		},
	}
	references, response, listerr := client.Git.ListBranches(context.Background(), repoID, scm.ListOptions{})
	if listerr != nil {
		t.Errorf("ListBranches got an error %v", listerr)
	}
	if response.Status != http.StatusOK {
		t.Errorf("ListBranches did not get a 200 back %v", response.Status)
	}
	if len(references) < 1 {
		t.Errorf("ListBranches  should have at least 1 branch %d", len(references))
	}
	if references[0].Sha == "" {
		t.Errorf("ListBranches first entry did not get a sha back %v", references[0].Sha)
	}
}

func TestCreateBranch(t *testing.T) {
	if token == "" {
		t.Skip("Skipping, Acceptance test")
	}
	client = azure.NewDefault(organization, project)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
			},
		},
	}
	currentCommit, commitErr := GetCurrentCommitOfBranch(client, "main")
	if commitErr != nil {
		t.Errorf("we got an error %v", commitErr)
	}
	input := &scm.CreateBranch{
		Name: "test_branch",
		Sha:  currentCommit,
	}
	response, listerr := client.Git.CreateBranch(context.Background(), repoID, input)
	if listerr != nil {
		t.Errorf("CreateBranch got an error %v", listerr)
	}
	if response.Status != http.StatusOK {
		t.Errorf("CreateBranch did not get a 200 back %v", response.Status)
	}

}

func TestListCommits(t *testing.T) {
	if token == "" {
		t.Skip("Skipping, Acceptance test")
	}
	client = azure.NewDefault(organization, project)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
			},
		},
	}
	commits, response, listerr := client.Git.ListCommits(context.Background(), repoID, scm.CommitListOptions{})
	if listerr != nil {
		t.Errorf("ListCommits  got an error %v", listerr)
	}
	if response.Status != http.StatusOK {
		t.Errorf("ListCommitsdid not get a 200 back %v", response.Status)
	}
	if len(commits) < 1 {
		t.Errorf("Contents.List there should be at least 1 commit %d", len(commits))
	}
	if commits[0].Sha == "" {
		t.Errorf("Contents.List first entry did not get a sha back %v", commits[0].Sha)
	}
}
