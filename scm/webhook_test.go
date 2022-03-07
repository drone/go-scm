package scm_test

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/stretchr/testify/require"
)

func TestWebhookWrapper(t *testing.T) {
	testCases := []struct {
		name   string
		verify func(name string, wh *scm.WebhookWrapper)
	}{
		{
			name: "push.json",
			verify: func(name string, wh *scm.WebhookWrapper) {
				require.NotNil(t, wh.PushHook, "no push hook for test %s", name)
			},
		},
	}

	dir := filepath.Join("test_data", "webhooks")
	for _, tc := range testCases {
		path := filepath.Join(dir, tc.name)
		require.FileExists(t, path)

		wh := &scm.WebhookWrapper{}

		data, err := ioutil.ReadFile(path)
		require.NoError(t, err, "failed to load file %s", path)

		err = json.Unmarshal(data, wh)
		require.NoError(t, err, "failed to unmarshal file %s", path)

		tc.verify(tc.name, wh)

		hook, err := wh.ToWebhook()
		require.NoError(t, err, "failed to parse hook")
		require.NotNil(t, hook, "nil wehhook returned")
	}
}
