package scm

import "time"

type (
	Execution struct {
		ID      string
		Status  string
		Created time.Time
		Updated time.Time
		URL     string
	}
)
