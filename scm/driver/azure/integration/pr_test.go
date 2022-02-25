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

func TestCreatePR(t *testing.T) {
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
	input := &scm.PullRequestInput{
		Title:  "test_pr",
		Body:   "test_pr_body",
		Source: "pr_branch",
		Target: "main",
	}
	outputPR, response, listerr := client.PullRequests.Create(context.Background(), repoID, input)
	if listerr != nil {
		t.Errorf("PullRequests.Create got an error %v", listerr)
	}
	if response.Status != http.StatusCreated {
		t.Errorf("PullRequests.Create did not get a 201 back %v", response.Status)
	}
	if outputPR.Title != "test_pr" {
		t.Errorf("PullRequests.Create does not have the correct title %v", outputPR.Title)
	}
}
