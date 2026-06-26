// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"testing"

	"github.com/drone/go-scm/scm"
)

func TestParseDiffGitPath(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{
			name: "simple path",
			line: "diff --git a/main.go b/main.go",
			want: "main.go",
		},
		{
			name: "nested path",
			line: "diff --git a/src/pkg/file.go b/src/pkg/file.go",
			want: "src/pkg/file.go",
		},
		{
			name: "rename uses destination (b/) path",
			line: "diff --git a/old.go b/new.go",
			want: "new.go",
		},
		{
			name: "too few fields",
			line: "diff --git a/main.go",
			want: "",
		},
		{
			name: "too many fields",
			line: "diff --git a/a b/b extra",
			want: "",
		},
		{
			name: "empty after prefix",
			line: "diff --git ",
			want: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := parseDiffGitPath(tc.line); got != tc.want {
				t.Errorf("parseDiffGitPath(%q) = %q, want %q", tc.line, got, tc.want)
			}
		})
	}
}

func TestSplitUnifiedDiff(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want map[string]string
	}{
		{
			name: "empty input",
			raw:  "",
			want: map[string]string{},
		},
		{
			name: "no diff --git headers",
			raw:  "some random text\nwithout headers",
			want: map[string]string{},
		},
		{
			name: "single file",
			raw: "diff --git a/main.go b/main.go\n" +
				"index 111..222 100644\n" +
				"--- a/main.go\n" +
				"+++ b/main.go\n" +
				"@@ -1,1 +1,1 @@\n" +
				"-old\n" +
				"+new\n",
			want: map[string]string{
				"main.go": "index 111..222 100644\n" +
					"--- a/main.go\n" +
					"+++ b/main.go\n" +
					"@@ -1,1 +1,1 @@\n" +
					"-old\n" +
					"+new",
			},
		},
		{
			name: "multiple files",
			raw: "diff --git a/a.go b/a.go\n" +
				"@@ -1 +1 @@\n" +
				"-a\n" +
				"+A\n" +
				"diff --git a/b.go b/b.go\n" +
				"@@ -1 +1 @@\n" +
				"-b\n" +
				"+B\n",
			want: map[string]string{
				"a.go": "@@ -1 +1 @@\n-a\n+A",
				"b.go": "@@ -1 +1 @@\n-b\n+B",
			},
		},
		{
			name: "rename keyed by destination path",
			raw: "diff --git a/old.go b/new.go\n" +
				"similarity index 90%\n" +
				"rename from old.go\n" +
				"rename to new.go\n",
			want: map[string]string{
				"new.go": "similarity index 90%\nrename from old.go\nrename to new.go",
			},
		},
		{
			name: "preamble before first header is ignored",
			raw: "junk line\n" +
				"diff --git a/x.go b/x.go\n" +
				"@@ -1 +1 @@\n" +
				"+x\n",
			want: map[string]string{
				"x.go": "@@ -1 +1 @@\n+x",
			},
		},
		{
			name: "malformed header yields empty path and is skipped",
			raw: "diff --git a/only-one-field\n" +
				"@@ -1 +1 @@\n" +
				"+content\n",
			want: map[string]string{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := splitUnifiedDiff(tc.raw)
			if len(got) != len(tc.want) {
				t.Fatalf("splitUnifiedDiff() returned %d entries, want %d: %#v", len(got), len(tc.want), got)
			}
			for path, wantPatch := range tc.want {
				if got[path] != wantPatch {
					t.Errorf("patch for %q:\n got: %q\nwant: %q", path, got[path], wantPatch)
				}
			}
		})
	}
}

func TestApplyUnifiedDiff(t *testing.T) {
	t.Run("no changes is a no-op", func(t *testing.T) {
		applyUnifiedDiff(nil, "diff --git a/a b/a\n@@ @@\n")
	})

	t.Run("empty raw leaves patches empty", func(t *testing.T) {
		changes := []*scm.Change{{Path: "a.go"}}
		applyUnifiedDiff(changes, "")
		if changes[0].Patch != "" {
			t.Errorf("expected empty patch, got %q", changes[0].Patch)
		}
	})

	t.Run("matches patch by path", func(t *testing.T) {
		changes := []*scm.Change{{Path: "a.go"}, {Path: "b.go"}}
		raw := "diff --git a/a.go b/a.go\n@@ -1 +1 @@\n+a\n" +
			"diff --git a/b.go b/b.go\n@@ -1 +1 @@\n+b\n"
		applyUnifiedDiff(changes, raw)
		if changes[0].Patch != "@@ -1 +1 @@\n+a" {
			t.Errorf("a.go patch = %q", changes[0].Patch)
		}
		if changes[1].Patch != "@@ -1 +1 @@\n+b" {
			t.Errorf("b.go patch = %q", changes[1].Patch)
		}
	})

	t.Run("falls back to PrevFilePath for renames", func(t *testing.T) {
		changes := []*scm.Change{{Path: "new.go", PrevFilePath: "old.go"}}
		raw := "diff --git a/old.go b/old.go\n@@ -1 +1 @@\n-x\n+y\n"
		applyUnifiedDiff(changes, raw)
		if changes[0].Patch != "@@ -1 +1 @@\n-x\n+y" {
			t.Errorf("rename patch = %q", changes[0].Patch)
		}
	})

	t.Run("leaves patch empty when no match", func(t *testing.T) {
		changes := []*scm.Change{{Path: "absent.go"}}
		raw := "diff --git a/other.go b/other.go\n@@ -1 +1 @@\n+z\n"
		applyUnifiedDiff(changes, raw)
		if changes[0].Patch != "" {
			t.Errorf("expected empty patch, got %q", changes[0].Patch)
		}
	})
}
