// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestContentFind(t *testing.T) {
	defer gock.Off()

	mockServerVersion()
	gock.New("https://try.gitea.io").
		Get("/api/v1/repos/go-gitea/gitea/contents/.gitignore").
		Reply(200).
		Type("application/json;charset=utf-8").
		File("testdata/content_find.json")

	client, _ := New("https://try.gitea.io")
	result, _, err := client.Contents.Find(
		context.Background(),
		"go-gitea/gitea",
		".gitignore",
		"master",
	)
	if err != nil {
		t.Error(err)
	}

	if got, want := result.Path, ".gitignore"; got != want {
		t.Errorf("Want file Path %q, got %q", want, got)
	}

	want := new(scm.Content)
	raw, _ := os.ReadFile("testdata/content_find.json.golden")
	err = json.Unmarshal(raw, want)
	assert.NoError(t, err)

	if diff := cmp.Diff(result, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}

func TestContentList(t *testing.T) {
	defer gock.Off()

	mockServerVersion()
	gock.New("https://try.gitea.io").
		Get("/api/v1/repos/go-gitea/gitea/contents/").
		Reply(200).
		Type("application/json;charset=utf-8").
		File("testdata/content_list.json")

	client, _ := New("https://try.gitea.io")
	result, _, err := client.Contents.List(
		context.Background(),
		"go-gitea/gitea",
		"/",
		"",
	)
	if err != nil {
		t.Error(err)
	}

	want := []*scm.FileEntry{}
	raw, _ := os.ReadFile("testdata/content_list.json.golden")
	err = json.Unmarshal(raw, &want)
	assert.NoError(t, err)

	if diff := cmp.Diff(result, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestContentCreate(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	message := "add README.md"
	content := []byte("Hello World")
	branch := "master"
	c := encode(content)

	r :=
		gock.New("https://try.gitea.io").
			Post("/api/v1/repos/go-gitea/gitea/contents/README.md").
			JSON(map[string]interface{}{
				"message":    message,
				"branch":     branch,
				"new_branch": "",
				"author":     map[string]string{"name": "", "email": ""},
				"committer":  map[string]string{"name": "", "email": ""},
				"dates":      map[string]string{"author": "0001-01-01T00:00:00Z", "committer": "0001-01-01T00:00:00Z"},
				"signoff":    false,
				"content":    c,
			})

	r.Header.Del("Content-Type")

	r.Reply(200).
		Type("application/json;charset=utf-8").
		File("testdata/content_create.json") // content is ignored anyway

	o := &scm.ContentParams{
		Data:    content,
		Branch:  branch,
		Message: message,
	}
	client, _ := New("https://try.gitea.io")
	_, err := client.Contents.Create(context.Background(), "go-gitea/gitea", "README.md", o)
	if err != nil {
		t.Error(err)
	}

}

func TestContentUpdate(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	previousSHA := "99094f9600108d9913c6e7d91c61ee5914cceb75"
	content := []byte("Hello World")
	message := "add README.md"
	branch := "master"

	r := gock.New("https://try.gitea.io").
		Put("/api/v1/repos/go-gitea/gitea/contents/README.md").
		MatchType("json").
		JSON(map[string]interface{}{
			"message":    message,
			"branch":     branch,
			"new_branch": "",
			"author":     map[string]string{"name": "", "email": ""},
			"committer":  map[string]string{"name": "", "email": ""},
			"dates":      map[string]string{"author": "0001-01-01T00:00:00Z", "committer": "0001-01-01T00:00:00Z"},
			"sha":        previousSHA,
			"content":    encode(content),
			"signoff":    false,
			"from_path":  "",
		})
	r.Header.Del("Content-Type")

	r.Reply(200).
		Type("application/json;charset=utf-8").
		File("testdata/content_create.json") // content is ignored anyway

	o := &scm.ContentParams{
		Data:    content,
		Branch:  branch,
		Message: message,
		Sha:     previousSHA,
	}
	client, _ := New("https://try.gitea.io")
	_, err := client.Contents.Update(context.Background(), "go-gitea/gitea", "README.md", o)
	if err != nil {
		t.Error(err)
	}
}

func TestContentDelete(t *testing.T) {
	// TODO disable for now as its down
	t.SkipNow()
	client, _ := New("https://try.gitea.io")
	_, err := client.Contents.Delete(context.Background(), "go-gitea/gitea", "README.md", "master")
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
