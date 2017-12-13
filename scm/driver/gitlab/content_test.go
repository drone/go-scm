// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"bytes"
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitlab/fixtures"
)

func TestContentFind(t *testing.T) {
	server := fixtures.NewServer()
	defer server.Close()

	client, _ := New(server.URL)
	result, _, err := client.Contents.Find(context.Background(), "diaspora/diaspora", "app/models/key.rb", "d5a3ff139356ce33e37e73add446f16869741b50")
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := result.Path, "app/models/key.rb"; got != want {
		t.Errorf("Want content Path %q, got %q", want, got)
	}
	if !bytes.Equal(result.Data, fileContent) {
		t.Errorf("Downloaded content does not match")
	}
}

func TestContentCreate(t *testing.T) {
	content := new(contentService)
	_, err := content.Create(context.Background(), "octocat/hello-world", "README", nil)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestContentUpdate(t *testing.T) {
	content := new(contentService)
	_, err := content.Update(context.Background(), "octocat/hello-world", "README", nil)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
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
