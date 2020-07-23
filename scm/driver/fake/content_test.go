package fake

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContent(t *testing.T) {
	client, _ := NewDefault()

	ctx := context.Background()
	sha := "master"

	repo := "myorg/myrepo"
	files, _, err := client.Contents.List(ctx, repo, "/", sha)
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
