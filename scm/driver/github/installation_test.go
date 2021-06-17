package github

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
	"github.com/livecycle/go-scm/scm"
)

func TestInstallationListRepositories(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/installation/repositories").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/installation_repos.json")

	client := NewDefault()
	got, res, err := client.Installation.ListRepositories(context.Background(), scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Repository{}
	raw, _ := ioutil.ReadFile("testdata/installation_repos.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}
