// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestRepositoryList(t *testing.T) {
	defer gock.Off()

	gock.New("https:/dev.azure.com/").
		Get("/ORG/PROJ/_apis/git/repositories").
		Reply(200).
		Type("application/json").
		File("testdata/repos.json")

	client := NewDefault("ORG", "PROJ")
	got, _, err := client.Repositories.List(context.Background(), scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Repository{}
	raw, _ := ioutil.ReadFile("testdata/repos.json.golden")
	jsonErr := json.Unmarshal(raw, &want)
	if jsonErr != nil {
		t.Error(jsonErr)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
