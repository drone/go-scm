// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/stretchr/testify/require"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestContentFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/repository/files/app/models/key.rb").
		MatchParam("ref", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content.json")

	client := NewDefault()
	got, res, err := client.Contents.Find(
		context.Background(),
		"diaspora/diaspora",
		"app/models/key.rb",
		"7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
	)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Content)
	raw, _ := ioutil.ReadFile("testdata/content.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestContentList(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitlab.com").
		Get("/api/v4/projects/diaspora/diaspora/repository/tree").
		MatchParam("ref", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content_list.json")

	client := NewDefault()
	got, res, err := client.Contents.List(
		context.Background(),
		"diaspora/diaspora",
		"app/models/key.rb",
		"7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
	)
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.FileEntry{}
	raw, _ := ioutil.ReadFile("testdata/content_list.json.golden")
	err = json.Unmarshal(raw, &want)
	require.NoError(t, err, "failed to unmarshal json")

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		data, err := json.MarshalIndent(got, "", "  ")
		if err == nil {
			t.Logf("got json:\n%s\n", string(data))
		} else {
			t.Logf("failed to marshal: %s\n", err.Error())
		}
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestContentCreate(t *testing.T) {
	defer gock.Off()
	message := "just a test message"
	content := []byte("testing")
	branch := "my-test-branch"

	encoded := base64.StdEncoding.EncodeToString([]byte(content))

	gock.New("https://gitlab.com").
		Post("api/v4/projects/octocat/hello-world/repository/commits").
		MatchType("json").
		JSON(map[string]interface{}{
			"branch":         branch,
			"id":             "octocat%2Fhello-world",
			"commit_message": message,
			"actions": []interface{}{
				map[string]interface{}{
					"action":    "create",
					"file_path": "README",
					"content":   encoded,
					"encoding":  "base64",
				},
			},
		}).
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content.json")

	params := &scm.ContentParams{
		Branch:  branch,
		Message: message,
		Data:    content,
	}
	client := NewDefault()

	_, err := client.Contents.Create(context.Background(), "octocat/hello-world", "README", params)
	if err != nil {
		t.Fatal(err)
	}
}

func TestContentUpdate(t *testing.T) {
	defer gock.Off()
	message := "just a test message"
	content := []byte("testing")
	branch := "my-test-branch"

	gock.New("https://gitlab.com").
		Put("api/v4/projects/octocat/hello-world/repository/files/README").
		MatchType("json").
		JSON(map[string]string{
			"branch":         branch,
			"content":        string(content),
			"commit_message": message,
		}).
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/content.json")

	params := &scm.ContentParams{
		Branch:  branch,
		Message: message,
		Data:    content,
	}
	client := NewDefault()

	res, err := client.Contents.Update(context.Background(), "octocat/hello-world", "README", params)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestContentDelete(t *testing.T) {
	content := new(contentService)
	_, err := content.Delete(context.Background(), "octocat/hello-world", "README", "master")
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

var fileContent = []byte(`require 'digest/md5'

class Key < ActiveRecord::Base
  include Gitlab::CurrentSettings
  include Sortable

  belongs_to :user

  before_validation :generate_fingerprint

  validates :title,
    presence: true,
    length: { maximum: 255 }

  validates :key,
    presence: true,
    length: { maximum: 5000 },
    format: { with: /\A(ssh|ecdsa)-.*\Z/ }

  validates :fingerprint,
    uniqueness: true,
    presence: { message: 'cannot be generated' }

  validate :key_meets_restrictions

  delegate :name, :email, to: :user, prefix: true

  after_commit :add_to_shell, on: :create
  after_create :post_create_hook
  after_create :refresh_user_cache
  after_commit :remove_from_shell, on: :destroy
  after_destroy :post_destroy_hook
  after_destroy :refresh_user_cache

  def key=(value)
    value&.delete!("\n\r")
    value.strip! unless value.blank`)
