package stash

import (
	"strings"
	"testing"
)

func TestRenderHunks_SingleHunk(t *testing.T) {
	t.Parallel()

	hunks := []*prHunk{
		{
			SourceLine: 1, SourceSpan: 2,
			DestinationLine: 1, DestinationSpan: 3,
			Segments: []*prSegment{
				{Type: "CONTEXT", Lines: []*prDiffLine{{Line: "package main"}}},
				{Type: "ADDED", Lines: []*prDiffLine{
					{Line: ""},
					{Line: "func main() {}"},
				}},
			},
		},
	}

	got := renderHunks(hunks)
	want := strings.Join([]string{
		"@@ -1,2 +1,3 @@",
		" package main",
		"+",
		"+func main() {}",
	}, "\n")

	if got != want {
		t.Fatalf("unexpected patch\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestRenderHunks_MultipleHunks(t *testing.T) {
	t.Parallel()

	hunks := []*prHunk{
		{
			SourceLine: 10, SourceSpan: 5,
			DestinationLine: 10, DestinationSpan: 7,
			Segments: []*prSegment{
				{Type: "CONTEXT", Lines: []*prDiffLine{{Line: " context"}}},
				{Type: "REMOVED", Lines: []*prDiffLine{{Line: "old"}}},
				{Type: "ADDED", Lines: []*prDiffLine{
					{Line: "new"},
					{Line: "added"},
				}},
				{Type: "CONTEXT", Lines: []*prDiffLine{{Line: " context"}}},
			},
		},
		{
			SourceLine: 30, SourceSpan: 3,
			DestinationLine: 32, DestinationSpan: 4,
			Segments: []*prSegment{
				{Type: "CONTEXT", Lines: []*prDiffLine{{Line: " context"}}},
				{Type: "ADDED", Lines: []*prDiffLine{{Line: "inserted"}}},
				{Type: "CONTEXT", Lines: []*prDiffLine{{Line: " context"}}},
			},
		},
	}

	got := renderHunks(hunks)
	if !strings.HasPrefix(got, "@@ -10,5 +10,7 @@") {
		t.Fatalf("expected first hunk header, got:\n%s", got)
	}
	if !strings.Contains(got, "@@ -30,3 +32,4 @@") {
		t.Fatalf("expected second hunk header, got:\n%s", got)
	}
}

func TestRenderHunks_NewFile(t *testing.T) {
	t.Parallel()

	hunks := []*prHunk{
		{
			SourceLine: 0, SourceSpan: 0,
			DestinationLine: 1, DestinationSpan: 2,
			Segments: []*prSegment{
				{Type: "ADDED", Lines: []*prDiffLine{
					{Line: "hello"},
					{Line: "world"},
				}},
			},
		},
	}

	got := renderHunks(hunks)
	want := "@@ -0,0 +1,2 @@\n+hello\n+world"
	if got != want {
		t.Fatalf("unexpected patch\nwant:\n%s\ngot:\n%s", want, got)
	}
}
