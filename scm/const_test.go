package scm

import (
	"encoding/json"
	"testing"
)

func TestStateJSON(t *testing.T) {
	for i := StateUnknown; i < StateExpected; i++ {
		in := i
		t.Run(in.String(), func(t *testing.T) {
			b, err := json.Marshal(in)
			if err != nil {
				t.Fatal(err)
			}

			var out State
			if err := json.Unmarshal(b, &out); err != nil {
				t.Fatal(err)
			}

			if in != out {
				t.Errorf("%s != %s", in, out)
			}
		})
	}
}

func TestActionJSON(t *testing.T) {
	for i := ActionCreate; i < ActionCompleted; i++ {
		in := i
		t.Run(in.String(), func(t *testing.T) {
			b, err := json.Marshal(in)
			if err != nil {
				t.Fatal(err)
			}

			var out Action
			if err := json.Unmarshal(b, &out); err != nil {
				t.Fatal(err)
			}

			if in != out {
				t.Errorf("%s != %s", in, out)
			}
		})
	}
}
