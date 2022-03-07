package fake_test

import (
	"context"
	"testing"

	"github.com/jenkins-x/go-scm/scm/driver/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x/go-scm/scm"
)

func TestHookCreateDelete(t *testing.T) {
	client, _ := fake.NewDefault()

	in := &scm.HookInput{
		Target: "https://example.com",
		Name:   "test",
	}

	// create a hook
	createdHook, _, err := client.Repositories.CreateHook(context.Background(), "foo/repo", in)
	if err != nil {
		t.Fatal(err)
	}

	id := createdHook.ID
	if id == "" {
		t.Fatal("created hook must have an ID")
	}

	// list to verify created hook
	hooks, _, err := client.Repositories.ListHooks(context.Background(), "foo/repo", scm.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}

	if len(hooks) != 1 {
		t.Fatal("expect one hook")
	}

	if diff := cmp.Diff(id, hooks[0].ID); diff != "" {
		t.Fatalf("hook id mismatch got\n%s", diff)
	}

	// delete by hook ID
	_, err = client.Repositories.DeleteHook(context.Background(), "foo/repo", id)
	if err != nil {
		t.Fatal(err)
	}

	// list to verify deletion
	hooks, _, err = client.Repositories.ListHooks(context.Background(), "foo/repo", scm.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}

	if len(hooks) != 0 {
		t.Fatal("expect no hooks")
	}
}

func TestForkRepository(t *testing.T) {
	client, _ := fake.NewDefault()

	org := "jenkins-x"
	repoName := "go-scm"
	username := client.Username
	expectedGitURL := "https://fake.com/" + username + "/" + repoName + ".git"
	fullName := scm.Join(org, repoName)
	forkFullName := scm.Join(username, repoName)

	ctx := context.TODO()

	fake.AssertNoRepoExists(ctx, t, client, fullName)
	fake.AssertNoRepoExists(ctx, t, client, forkFullName)

	repo, _, err := client.Repositories.Create(ctx, &scm.RepositoryInput{
		Namespace: org,
		Name:      repoName,
	})
	require.NoError(t, err, "failed to create repo %s", fullName)
	require.NotNil(t, repo, "no repo returned for create repo %s", fullName)

	fake.AssertRepoExists(ctx, t, client, fullName)

	_, _, err = client.Repositories.Fork(ctx, &scm.RepositoryInput{
		Name: repoName,
	}, repoName)
	if err != nil {
		t.Error(err)
	}

	repository := fake.AssertRepoExists(ctx, t, client, forkFullName)
	assert.Equal(t, expectedGitURL, repository.Clone, "forked repository clone URL")
}
