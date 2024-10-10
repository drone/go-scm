package gitea

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x/go-scm/scm"
	"gopkg.in/h2non/gock.v1"
)

func TestConvertAPIURLToHTMLURL(t *testing.T) {

	got := ConvertAPIURLToHTMLURL("https://try.gitea.com/api/v1/repos/octocat/Hello-World/123", "v1.0.0")
	want := "https://try.gitea.com/octocat/Hello-World/releases/tag/v1.0.0"

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		t.Log("got:")
		t.Log(got)
	}
}

func TestConvertAPIURLToHTMLURLEmptyLinkWhenURLParseFails(t *testing.T) {

	broken := []string{"http s://try.gitea.com/api/v1/repos/octocat/Hello-World/123", "https://try.gitea.com/api/v1/repos/octocat/Hello-World"}
	for _, url := range broken {

		got := ConvertAPIURLToHTMLURL(url, "v1.0.0")
		want := ""

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Unexpected Results")
			t.Log(diff)

			t.Log("got:")
			t.Log(got)
		}
	}
}
func TestReleaseFind(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/repos/octocat/hello-world/releases/1").
		Reply(200).
		Type("application/json").
		File("testdata/release.json")

	client, err := New("https://demo.gitea.com")
	if err != nil {
		t.Error(err)
		return
	}
	got, _, err := client.Releases.Find(context.Background(), "octocat/hello-world", 1)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Release)
	raw, _ := os.ReadFile("testdata/release.json.golden")
	err = json.Unmarshal(raw, want)
	assert.NoError(t, err)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		data, _ := json.Marshal(got)
		t.Log("got JSON:")
		t.Log(string(data))
	}
}

func TestReleaseList(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/repos/octocat/hello-world/releases").
		MatchParam("page", "1").
		MatchParam("limit", "30").
		Reply(200).
		Type("application/json").
		File("testdata/releases.json")

	client, err := New("https://demo.gitea.com")
	if err != nil {
		t.Error(err)
		return
	}

	got, _, err := client.Releases.List(context.Background(), "octocat/hello-world", scm.ReleaseListOptions{Page: 1, Size: 30, Open: true, Closed: true})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Release{}
	raw, _ := os.ReadFile("testdata/releases.json.golden")
	err = json.Unmarshal(raw, &want)
	assert.NoError(t, err)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		data, _ := json.Marshal(got)
		t.Log("got JSON:")
		t.Log(string(data))
	}

}

func TestReleaseCreate(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Post("/repos/octocat/hello-world/releases").
		File("testdata/release_create.json").
		Reply(200).
		Type("application/json").
		File("testdata/release.json")

	client, err := New("https://demo.gitea.com")
	if err != nil {
		t.Error(err)
		return
	}
	input := &scm.ReleaseInput{
		Title:       "v1.0.0",
		Description: "Description of the release",
		Tag:         "v1.0.0",
		Commitish:   "master",
		Draft:       false,
		Prerelease:  false,
	}

	got, _, err := client.Releases.Create(context.Background(), "octocat/hello-world", input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Release)
	raw, _ := os.ReadFile("testdata/release.json.golden")
	err = json.Unmarshal(raw, want)
	assert.NoError(t, err)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		data, _ := json.Marshal(got)
		t.Log("got JSON:")
		t.Log(string(data))
	}

}

func TestReleaseUpdate(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Patch("/repos/octocat/hello-world/releases/1").
		File("testdata/release_update.json").
		Reply(200).
		Type("application/json").
		File("testdata/release.json")

	client, err := New("https://demo.gitea.com")
	if err != nil {
		t.Error(err)
		return
	}
	input := &scm.ReleaseInput{
		Title:       "v1.0.0",
		Description: "Description of the release",
		Tag:         "v1.0.0",
		Commitish:   "master",
		Draft:       false,
		Prerelease:  false,
	}
	got, _, err := client.Releases.Update(context.Background(), "octocat/hello-world", 1, input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Release)
	raw, _ := os.ReadFile("testdata/release.json.golden")
	err = json.Unmarshal(raw, want)
	assert.NoError(t, err)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

}

func TestReleaseDelete(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Delete("/repos/octocat/hello-world/releases/1").
		Reply(200).
		Type("application/json")

	client, err := New("https://demo.gitea.com")
	if err != nil {
		t.Error(err)
	}
	_, err = client.Releases.Delete(context.Background(), "octocat/hello-world", 1)
	if err != nil {
		t.Error(err)
		return
	}

}
