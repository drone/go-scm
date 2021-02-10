package codecommit

import (
	"github.com/drone/go-scm/scm"
)

// NewDefault returns a new codecommit client.
func NewDefault() *scm.Client {
	client := &wrapper{}
	client.Webhooks = &webhookService{}
	return client.Client
}

// wraper wraps the Client to provide high level helper functions.
type wrapper struct {
	*scm.Client
}
