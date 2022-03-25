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

func TestListRepos(t *testing.T) {
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
	references, response, listerr := client.Repositories.List(context.Background(), scm.ListOptions{})
	if listerr != nil {
		t.Errorf("List got an error %v", listerr)
	}
	if response.Status != http.StatusOK {
		t.Errorf("List did not get a 200 back %v", response.Status)
	}
	if len(references) < 1 {
		t.Errorf("List should have at least 1 repo %d", len(references))
	}
}
