# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 1.7.0
### Added
- Improve status display text in new bitbucket pull request screen, from [@bradrydzewski](https://github.com/bradrydzewski). See [#27](https://github.com/drone/go-scm/issues/27).
- Implement timestamp value for GitHub push webhooks, from [@bradrydzewski](https://github.com/bradrydzewski).
- Implement deep link to branch.
- Implement git compare function to compare two separate commits, from [@chhsia0](https://github.com/chhsia0).
- Implement support for creating and updating GitLab and GitHub repository contents, from [@zhuxiaoyang](https://github.com/zhuxiaoyang).
- Capture Repository link for Gitea, Gogs and Gitlab, from [@chhsia0](https://github.com/chhsia0).

### Fixed
- Fix issue with GitHub enterprise deep link including API prefix, from [@bradrydzewski](https://github.com/bradrydzewski).
- Fix issue with GitHub deploy hooks for commits having an invalid reference, from [@bradrydzewski](https://github.com/bradrydzewski).
- Support for Skipping SSL verification for GitLab webhooks. See [#40](https://github.com/drone/go-scm/pull/40).
- Support for Skipping SSL verification for GitHub webhooks. See [#44](https://github.com/drone/go-scm/pull/40).
- Fix issue with handling slashes in Bitbucket branch names. See [#7](https://github.com/drone/go-scm/pull/47).
- Fix incorrect Gitea tag link. See [#52](https://github.com/drone/go-scm/pull/52).
- Encode ref when making Gitea API calls. See [#61](https://github.com/drone/go-scm/pull/61).

## [1.6.0]
### Added
- Support Head and Base sha for GitHub pull request, from [@bradrydzewski](https://github.com/bradrydzewski).
- Support Before sha for Bitbucket, from [@jkdev81](https://github.com/jkdev81).
- Support for creating GitHub deployment hooks, from [@bradrydzewski](https://github.com/bradrydzewski).
- Endpoint to get organization membership for GitHub, from [@bradrydzewski](https://github.com/bradrydzewski).
- Functions to generate deep links to git resources, from [@bradrydzewski](https://github.com/bradrydzewski).

### Fixed
- Fix issue getting a GitLab commit by ref, from [@bradrydzewski](https://github.com/bradrydzewski).

## [1.5.0]
### Added

- Fix missing sha for Gitea tag hooks, from [@techknowlogick](https://github.com/techknowlogick). See [#22](https://github.com/drone/go-scm/pull/22).
- Support for Gitea webhook signature verification, from [@techknowlogick](https://github.com/techknowlogick).

## [1.4.0]
### Added

- Fix issues base64 decoding GitLab content, from [@bradrydzewski](https://github.com/bradrydzewski).

## [1.3.0]
### Added

- Fix missing avatar in Gitea commit from [@jgeek1011](https://github.com/geek1011).
- Implement GET commit endpoint for Gogs from [@ogarcia](https://github.com/ogarcia).
