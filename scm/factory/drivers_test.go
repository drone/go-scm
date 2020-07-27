package factory

import (
	"os"
	"regexp"
	"testing"
)

func TestIdentify(t *testing.T) {
	urlTests := []struct {
		name     string
		hostname string
		want     string
		envVar   string
		wantErr  string
	}{
		{"simple github", "github.com", "github", "", ""},
		{"simple gitlab", "gitlab.com", "gitlab", "", ""},
		{"from environment mapping", "gl.example.com", "gitlab", "gl.example.com=gitlab,gh.github.com=github", ""},
		{"unknown host", "scm.example.com", "", "", "unable to identify driver"},
	}
	origEnv := os.Getenv("GIT_DRIVERS")
	defer func() {
		os.Setenv("GIT_DRIVERS", origEnv)
	}()

	for _, tt := range urlTests {
		t.Run(tt.hostname, func(rt *testing.T) {
			if tt.envVar != "" {
				os.Setenv("GIT_DRIVERS", tt.envVar)
			}
			identifier := NewDriverIdentifier()
			driver, err := identifier.Identify(tt.hostname)
			if !matchError(rt, tt.wantErr, err) {
				rt.Errorf("error failed to match, got %#v, want %s", err, tt.wantErr)
			}
			if driver != tt.want {
				rt.Errorf("got %s, want %s", driver, tt.want)
			}
		})
	}
}

func TestIdentifyWithExtras(t *testing.T) {
	identifier := NewDriverIdentifier(Mapping("test.example.com", "gitlab"))
	driver, err := identifier.Identify("test.example.com")
	if err != nil {
		t.Fatal(err)
	}
	if driver != "gitlab" {
		t.Fatalf("got %q, want %q", driver, "gitlab")
	}
}

func matchError(t *testing.T, s string, e error) bool {
	t.Helper()
	if s == "" && e == nil {
		return true
	}
	if s != "" && e == nil {
		return false
	}
	match, err := regexp.MatchString(s, e.Error())
	if err != nil {
		t.Fatal(err)
	}
	return match
}
