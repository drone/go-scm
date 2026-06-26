// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"testing"

	"github.com/drone/go-scm/scm"
)

func TestRenderHunks(t *testing.T) {
	tests := []struct {
		name  string
		hunks []*hunk
		want  string
	}{
		{
			name:  "no hunks",
			hunks: nil,
			want:  "",
		},
		{
			name: "added lines",
			hunks: []*hunk{
				{
					SourceLine: 1, SourceSpan: 0, DestinationLine: 1, DestinationSpan: 3,
					Segments: []*segment{
						{Type: "ADDED", Lines: []*segmentLine{
							{Line: "package main"},
							{Line: ""},
							{Line: "func main() {}"},
						}},
					},
				},
			},
			want: "@@ -1,0 +1,3 @@\n+package main\n+\n+func main() {}",
		},
		{
			name: "removed lines",
			hunks: []*hunk{
				{
					SourceLine: 1, SourceSpan: 1, DestinationLine: 0, DestinationSpan: 0,
					Segments: []*segment{
						{Type: "REMOVED", Lines: []*segmentLine{{Line: "old readme"}}},
					},
				},
			},
			want: "@@ -1,1 +0,0 @@\n-old readme",
		},
		{
			name: "context plus added and removed",
			hunks: []*hunk{
				{
					SourceLine: 1, SourceSpan: 2, DestinationLine: 1, DestinationSpan: 2,
					Segments: []*segment{
						{Type: "CONTEXT", Lines: []*segmentLine{{Line: "keep"}}},
						{Type: "REMOVED", Lines: []*segmentLine{{Line: "drop"}}},
						{Type: "ADDED", Lines: []*segmentLine{{Line: "insert"}}},
					},
				},
			},
			want: "@@ -1,2 +1,2 @@\n keep\n-drop\n+insert",
		},
		{
			name: "unknown segment type defaults to context prefix",
			hunks: []*hunk{
				{
					SourceLine: 5, SourceSpan: 1, DestinationLine: 5, DestinationSpan: 1,
					Segments: []*segment{
						{Type: "MYSTERY", Lines: []*segmentLine{{Line: "unchanged"}}},
					},
				},
			},
			want: "@@ -5,1 +5,1 @@\n unchanged",
		},
		{
			name: "multiple hunks concatenated",
			hunks: []*hunk{
				{
					SourceLine: 1, SourceSpan: 1, DestinationLine: 1, DestinationSpan: 1,
					Segments: []*segment{{Type: "ADDED", Lines: []*segmentLine{{Line: "a"}}}},
				},
				{
					SourceLine: 10, SourceSpan: 1, DestinationLine: 10, DestinationSpan: 1,
					Segments: []*segment{{Type: "REMOVED", Lines: []*segmentLine{{Line: "b"}}}},
				},
			},
			want: "@@ -1,1 +1,1 @@\n+a\n@@ -10,1 +10,1 @@\n-b",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := renderHunks(tc.hunks); got != tc.want {
				t.Errorf("renderHunks():\n got: %q\nwant: %q", got, tc.want)
			}
		})
	}
}

func TestApplyDiffs(t *testing.T) {
	t.Run("nil response is a no-op", func(t *testing.T) {
		changes := []*scm.Change{{Path: "a.go"}}
		applyDiffs(changes, nil)
		if changes[0].Patch != "" {
			t.Errorf("expected empty patch, got %q", changes[0].Patch)
		}
	})

	t.Run("empty diffs is a no-op", func(t *testing.T) {
		changes := []*scm.Change{{Path: "a.go"}}
		applyDiffs(changes, &diffResponse{})
		if changes[0].Patch != "" {
			t.Errorf("expected empty patch, got %q", changes[0].Patch)
		}
	})

	t.Run("added file keyed by destination", func(t *testing.T) {
		changes := []*scm.Change{{Path: "main.go", Added: true}}
		resp := &diffResponse{Diffs: []*diff{
			{
				Destination: &diffPath{ToString: "main.go"},
				Hunks: []*hunk{{
					SourceLine: 1, SourceSpan: 0, DestinationLine: 1, DestinationSpan: 1,
					Segments: []*segment{{Type: "ADDED", Lines: []*segmentLine{{Line: "x"}}}},
				}},
			},
		}}
		applyDiffs(changes, resp)
		if changes[0].Patch != "@@ -1,0 +1,1 @@\n+x" {
			t.Errorf("patch = %q", changes[0].Patch)
		}
	})

	t.Run("deleted file keyed by source when destination nil", func(t *testing.T) {
		changes := []*scm.Change{{Path: "README", Deleted: true}}
		resp := &diffResponse{Diffs: []*diff{
			{
				Source: &diffPath{ToString: "README"},
				Hunks: []*hunk{{
					SourceLine: 1, SourceSpan: 1, DestinationLine: 0, DestinationSpan: 0,
					Segments: []*segment{{Type: "REMOVED", Lines: []*segmentLine{{Line: "bye"}}}},
				}},
			},
		}}
		applyDiffs(changes, resp)
		if changes[0].Patch != "@@ -1,1 +0,0 @@\n-bye" {
			t.Errorf("patch = %q", changes[0].Patch)
		}
	})

	t.Run("falls back to PrevFilePath for renames", func(t *testing.T) {
		changes := []*scm.Change{{Path: "new.go", PrevFilePath: "old.go", Renamed: true}}
		resp := &diffResponse{Diffs: []*diff{
			{
				Source: &diffPath{ToString: "old.go"},
				Hunks: []*hunk{{
					SourceLine: 1, SourceSpan: 1, DestinationLine: 1, DestinationSpan: 1,
					Segments: []*segment{{Type: "ADDED", Lines: []*segmentLine{{Line: "moved"}}}},
				}},
			},
		}}
		applyDiffs(changes, resp)
		if changes[0].Patch != "@@ -1,1 +1,1 @@\n+moved" {
			t.Errorf("rename patch = %q", changes[0].Patch)
		}
	})

	t.Run("diff with no usable path is skipped", func(t *testing.T) {
		changes := []*scm.Change{{Path: "a.go"}}
		resp := &diffResponse{Diffs: []*diff{
			{Source: nil, Destination: nil, Hunks: nil},
		}}
		applyDiffs(changes, resp)
		if changes[0].Patch != "" {
			t.Errorf("expected empty patch, got %q", changes[0].Patch)
		}
	})

	t.Run("non-matching change left empty", func(t *testing.T) {
		changes := []*scm.Change{{Path: "absent.go"}}
		resp := &diffResponse{Diffs: []*diff{
			{
				Destination: &diffPath{ToString: "present.go"},
				Hunks: []*hunk{{
					DestinationLine: 1, DestinationSpan: 1,
					Segments: []*segment{{Type: "ADDED", Lines: []*segmentLine{{Line: "z"}}}},
				}},
			},
		}}
		applyDiffs(changes, resp)
		if changes[0].Patch != "" {
			t.Errorf("expected empty patch, got %q", changes[0].Patch)
		}
	})
}
