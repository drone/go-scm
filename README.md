# go-scm

[![Go Doc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](http://godoc.org/github.com/drone/go-scm/scm)

Package scm provides a unified interface to multiple source code management systems including GitHub, GitHub Enterprise, Bitbucket, Bitbucket Server, Gitea and Gogs.

## Getting Started

Create a GitHub client:

```Go
package main

import (
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
)

func main() {
	client := github.NewDefault()
}
```

Create a GitHub Enterprise client:

```Go
import (
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
)

func main() {
    client, err := github.New("https://github.company.com/api/v3")
}
```

Create a Bitbucket client:

```Go
import (
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/bitbucket"
)

func main() {
    client, err := bitbucket.New()
}
```

Create a Bitbucket Server (Stash) client:

```Go
import (
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/stash"
)

func main() {
    client, err := stash.New("https://stash.company.com")
}
```

Create a Gitea client:

```Go
import (
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitea"
)

func main() {
    client, err := gitea.New("https://gitea.company.com")
}
```

## Authentication

The scm client does not directly handle authentication. Instead, when creating a new client, provide a `http.Client` that can handle authentication for you. For convenience, this library includes oauth1 and oauth2 implementations that can be used to authenticate requests.

```Go
package main

import (
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/drone/go-scm/scm/driver/transport/oauth2"
)

func main() {
	client := github.NewDefault()

	// provide a custom http.Client with a transport
	// that injects the oauth2 token.
	client.Client := &http.Client{
		Transport: &Transport{
			Source: StaticTokenSource(
				&scm.Token{
					Token: "ecf4c1f9869f59758e679ab54b4",
				},
			),
		},
	}

	// provide a custom http.Client with a transport
	// that injects the private GitLab token through
	// the PRIVATE_TOKEN header variable.
	client.Client := &http.Client{
		Transport: &PrivateToken{
			Token: "ecf4c1f9869f59758e679ab54b4",
		},
	}
}
```

## Release procedure

Run the changelog generator.

```BASH
docker run -it --rm -v "$(pwd)":/usr/local/src/your-app githubchangeloggenerator/github-changelog-generator -u drone -p go-scm -t <secret github token>
```

You can generate a token by logging into your GitHub account and going to Settings -> Personal access tokens.

Next we tag the PR's with the fixes or enhancements labels. If the PR does not fufil the requirements, do not add a label.

Run the changelog generator again with the future version according to semver.

```BASH
docker run -it --rm -v "$(pwd)":/usr/local/src/your-app githubchangeloggenerator/github-changelog-generator -u drone -p go-scm -t <secret token> --future-release v1.15.2
```

Create your pull request for the release. Get it merged then tag the release.
