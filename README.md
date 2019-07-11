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

## Fake driver for testing

When testing the use of go-scm its really handy to use the [fake](https://github.com/jenkins-x/go-scm/blob/master/scm/driver/fake/fake.go) provider which lets you populate the in memory resources inside the driver or query resources after a test has run.

```go 
client, data := fake.NewDefault()
```    
