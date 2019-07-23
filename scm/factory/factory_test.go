package factory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("fake", "", "")
	if client == nil {
		t.Errorf("no client created")
	}
	if err != nil {
		t.Errorf("failed to create client %s", err.Error())
	}
}

func TestGHEEndpoint(t *testing.T) {
	assert.Equal(t, "https://my.ghe.com/custom/api/v5", ensureGHEEndpoint("https://my.ghe.com/custom/api/v5"))
	assert.Equal(t, "https://my.ghe.com/custom/api/v3", ensureGHEEndpoint("https://my.ghe.com/custom"))
	assert.Equal(t, "https://my.ghe.com/api/v3", ensureGHEEndpoint("https://my.ghe.com"))
}

func TestBitBucketEndpoint(t *testing.T) {
	assert.Equal(t, "https://my.bitbucket.com/custom/rest", ensureBitBucketServerEndpoint("https://my.bitbucket.com/custom"))
	assert.Equal(t, "https://my.bitbucket.com/rest", ensureBitBucketServerEndpoint("https://my.bitbucket.com"))
}
