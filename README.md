# go-scm

A small library with minimal depenencies for working with Webhooks, Commits, Issues, Pull Requests, Comments, Reviews, Teams and more on multiple git provider:

* [GitHub](https://github.com/jenkins-x/go-scm/blob/master/scm/driver/github/github.go#L46)
* [GitHub Enterprise](https://github.com/jenkins-x/go-scm/blob/master/scm/driver/github/github.go#L19) (you specify a server URL)
* [BitBucket Server](https://github.com/jenkins-x/go-scm/blob/master/scm/driver/stash/stash.go#L24)
* [BitBucket Cloud](https://github.com/jenkins-x/go-scm/blob/master/scm/driver/bitbucket/bitbucket.go#L20)
* [GitLab](https://github.com/jenkins-x/go-scm/blob/master/scm/driver/gitlab/gitlab.go#L19)
* [Gitea](https://github.com/jenkins-x/go-scm/blob/master/scm/driver/gitea/gitea.go#L22)
* [Gogs](https://github.com/jenkins-x/go-scm/blob/master/scm/driver/gogs/gogs.go#L22)
* [Fake](https://github.com/jenkins-x/go-scm/blob/master/scm/driver/fake/fake.go)

## Working on the code

Clone this repository and use go test...

``` 
git clone https://github.com/jenkins-x/go-scm.git
cd go-scm
go test ./...
```

## Community

We have a [kanban board](https://github.com/jenkins-x/go-scm/projects/1?add_cards_query=is%3Aopen) of stuff to work on if you fancy contributing!

You can also find us [on Slack](http://slack.k8s.io/) at [kubernetes.slack.com](https://kubernetes.slack.com/):

* [\#jenkins-x-dev](https://kubernetes.slack.com/messages/C9LTHT2BB) for developers of Jenkins X and related OSS projects
* [\#jenkins-x-user](https://kubernetes.slack.com/messages/C9MBGQJRH) for users of Jenkins X


## Building

See the [guide to building and running the code](BUILDING.md)

## Writing tests

There are lots of tests for each driver; using sample JSON that comes from the git provider together with the expected canonical JSON.

e.g. I added this test for ListTeams on github: https://github.com/jenkins-x/go-scm/blob/master/scm/driver/github/org_test.go#L83-116

you then add some real json from the git provider: https://github.com/jenkins-x/go-scm/blob/master/scm/driver/github/testdata/teams.json and provide the expected json: https://github.com/jenkins-x/go-scm/blob/master/scm/driver/github/testdata/teams.json.golden

## Git API Reference docs

To help hack on the different drivers here's a list of docs which outline the git providers REST APIs

### GitHub

* REST API reference: https://developer.github.com/v3/
* WebHooks: https://developer.github.com/v3/activity/events/types/

### Bitbucket Server

* REST API reference:  https://docs.atlassian.com/bitbucket-server/rest/6.5.1/bitbucket-rest.htm
* Webhooks: https://confluence.atlassian.com/bitbucketserver/event-payload-938025882.html

## Fake driver for testing

When testing the use of go-scm its really handy to use the [fake](https://github.com/jenkins-x/go-scm/blob/master/scm/driver/fake/fake.go) provider which lets you populate the in memory resources inside the driver or query resources after a test has run.

```go 
client, data := fake.NewDefault()
```    
