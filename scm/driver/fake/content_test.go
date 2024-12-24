package fake_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jenkins-x/go-scm/scm"

	"github.com/jenkins-x/go-scm/scm/driver/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContent(t *testing.T) {
	client, _ := fake.NewDefault()

	ctx := context.Background()
	sha := "master"

	repo := "myorg/myrepo"
	files, _, err := client.Contents.List(ctx, repo, "/", sha, &scm.ListOptions{})
	require.NoError(t, err, "could not list files in repo %s", repo)
	require.Len(t, files, 2, "should have found 2 files")

	for _, f := range files {
		t.Logf("file: %s path: %s type: %s size: %v\n", f.Name, f.Path, f.Type, f.Size)

		switch f.Name {
		case "README.md":
			assert.Equal(t, "file", f.Type, "for path %s", f.Path)
		case "somedir":
			assert.Equal(t, "dir", f.Type, "for path %s", f.Path)
		default:
			assert.Fail(t, "invalid file name %s", f.Name)
		}
	}

	path := "README.md"
	c, _, err := client.Contents.Find(ctx, repo, path, sha)
	require.NoError(t, err, "could not find file in repo %s path %s", repo, path)

	text := string(c.Data)
	t.Logf("loaded repo %s path %s got %s\n", repo, path, text)

	assert.Contains(t, text, "root dir of a repo", "for repo %s path %s", repo, path)
}

func TestContentWithRefs(t *testing.T) {
	client, fakeData := fake.NewDefault()
	fakeData.ContentDir = filepath.Join("testdata", "test_refs")

	ctx := context.Background()
	repo := "myorg/myrepo"
	path := "something.txt"

	testCases := []struct {
		ref      string
		expected string
	}{
		{
			ref:      "master",
			expected: "main text",
		},
		{
			ref:      "mybranch",
			expected: "my changes on a branch",
		},
	}
	for _, tc := range testCases {
		ref := tc.ref
		c, _, err := client.Contents.Find(ctx, repo, path, ref)
		require.NoError(t, err, "could not find file in repo %s path %s", repo, path)

		text := strings.TrimSpace(string(c.Data))

		assert.Equal(t, tc.expected, text, "for repo path %s ref %", repo, ref, path)
		t.Logf("loaded repo %s path %s ref %s got %s\n", repo, ref, path, text)
	}
}
