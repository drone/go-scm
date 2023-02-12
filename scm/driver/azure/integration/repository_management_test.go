package integration_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jenkins-x/go-scm/scm"
)

func TestRepositoryManagement(t *testing.T) {
	canRun(t)
	const testRepoName = "repo-management-test"

	clean := func(t *testing.T) {
		_, err := client.Repositories.Delete(context.Background(), repoFQ(testRepoName))
		if err == nil {
			t.Logf("cleaned up leftover repository: %s", repoFQ(testRepoName))
		}
	}

	tests := []TestCase{
		{
			Name: "create new repo",
			Test: func(t *testing.T) {
				ri := scm.RepositoryInput{
					Namespace: projectFQ,
					Name:      testRepoName,
				}
				createdRepo, res, err := client.Repositories.Create(context.Background(), &ri)
				if err != nil {
					t.Fatalf("could not create repo: %+v", err)
				}
				if res.Status != http.StatusCreated {
					t.Fatalf("expected 201, got %d", res.Status)
				}

				if createdRepo.FullName != repoFQ(testRepoName) {
					t.Fatalf("expected %s, but got %s", repoFQ(testRepoName), createdRepo.FullName)
				}
			},
		},
		{
			Name: "list for org",
			Test: func(t *testing.T) {
				repoList, res, err := client.Repositories.ListOrganisation(context.Background(), organization, &scm.ListOptions{})
				if err != nil {
					t.Fatalf("could not list for org: %+v", err)
				}
				if res.Status != http.StatusOK {
					t.Fatalf("expected 200, got %d", res.Status)
				}

				found := false
				for _, repo := range repoList {
					if repo.FullName == repoFQ(testRepoName) {
						found = true
						break
					}
				}

				if !found {
					t.Fatalf("could not find repo: %s", repoFQ(testRepoName))
				}
			},
		},
		{
			Name: "create existing repo",
			Test: func(t *testing.T) {
				ri := scm.RepositoryInput{
					Namespace: projectFQ,
					Name:      testRepoName,
				}
				_, res, err := client.Repositories.Create(context.Background(), &ri)
				if err == nil {
					t.Fatalf("create repo should have failed")
				}
				if res.Status != http.StatusConflict {
					t.Fatalf("expected 409, got %d", res.Status)
				}
			},
		},
		{
			Name: "create with wrong org and project",
			Test: func(t *testing.T) {
				ri := scm.RepositoryInput{
					Namespace: "not/exist",
					Name:      testRepoName,
				}
				_, res, err := client.Repositories.Create(context.Background(), &ri)
				if err == nil {
					t.Fatalf("create repo should have failed")
				}
				if res.Status != http.StatusUnauthorized {
					t.Fatalf("expected 401, got %d", res.Status)
				}

			},
		},
		{
			Name: "create with wrong project",
			Test: func(t *testing.T) {
				ri := scm.RepositoryInput{
					Namespace: organization + "/not_exist",
					Name:      testRepoName,
				}
				_, res, err := client.Repositories.Create(context.Background(), &ri)
				if err == nil {
					t.Fatalf("create repo should have failed")
				}
				if res.Status != http.StatusNotFound {
					t.Fatalf("expected 404, got %d", res.Status)
				}
			},
		},
		{
			Name: "delete existing repo",
			Test: func(t *testing.T) {
				res, err := client.Repositories.Delete(context.Background(), repoFQ(testRepoName))
				if err != nil {
					t.Fatalf("could not delete repo")
				}
				if res.Status < 200 || res.Status >= 300 {
					t.Fatalf("expected 200 level status, got %d", res.Status)
				}
			},
		},
		{
			Name: "delete missing repo",
			Test: func(t *testing.T) {
				res, err := client.Repositories.Delete(context.Background(), repoFQ("not-exist"))
				if err == nil {
					t.Fatalf("delete nonexistant repo should fail")
				}
				if res.Status != http.StatusNotFound {
					t.Fatalf("expected 404, got %d", res.Status)
				}
			},
		},
	}

	clean(t)

	for _, tc := range tests {
		if !t.Run(tc.Name, tc.Test) {
			t.Fatal("test aborted")
		}
	}
}
