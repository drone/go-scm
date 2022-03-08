# Changelog

## [v1.19.1](https://github.com/drone/go-scm/tree/v1.19.1) (2022-02-23)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.19.0...v1.19.1)

**Fixed bugs:**

- Bitbucket list files fix [\#154](https://github.com/drone/go-scm/pull/154) ([mohitg0795](https://github.com/mohitg0795))
- GitHub list commits fix [\#152](https://github.com/drone/go-scm/pull/152) ([mohitg0795](https://github.com/mohitg0795))
- Bitbucket compare changes fix for rename and removed file ops [\#151](https://github.com/drone/go-scm/pull/151) ([mohitg0795](https://github.com/mohitg0795))

## [v1.19.0](https://github.com/drone/go-scm/tree/v1.19.0) (2022-02-09)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.18.0...v1.19.0)

**Implemented enhancements:**

- \(feat\) add path support for list commits on github and gitlab [\#149](https://github.com/drone/go-scm/pull/149) ([tphoney](https://github.com/tphoney))
- Extending bitbucket listCommits API to fetch commits for a given file [\#148](https://github.com/drone/go-scm/pull/148) ([mohitg0795](https://github.com/mohitg0795))
- Update GitHub signature header to use sha256 [\#123](https://github.com/drone/go-scm/pull/123) ([nlecoy](https://github.com/nlecoy))

**Merged pull requests:**

- v1.19.0 release prep [\#150](https://github.com/drone/go-scm/pull/150) ([tphoney](https://github.com/tphoney))

## [v1.18.0](https://github.com/drone/go-scm/tree/v1.18.0) (2022-01-18)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.17.0...v1.18.0)

**Implemented enhancements:**

- Added support for parsing prevFilePath field from github compare commits API response [\#143](https://github.com/drone/go-scm/pull/143) ([mohitg0795](https://github.com/mohitg0795))

**Fixed bugs:**

- Implement parsing/handling for missing pull request webhook events for BitBucket Server \(Stash\) driver [\#130](https://github.com/drone/go-scm/pull/130) ([raphendyr](https://github.com/raphendyr))

**Closed issues:**

- Bitbucket Stash driver doesn't handle event `pr:from_ref_updated` \(new commits / force push\) [\#116](https://github.com/drone/go-scm/issues/116)

**Merged pull requests:**

- release prep v1.18.0 [\#147](https://github.com/drone/go-scm/pull/147) ([eoinmcafee00](https://github.com/eoinmcafee00))

## [v1.17.0](https://github.com/drone/go-scm/tree/v1.17.0) (2022-01-07)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.16.3...v1.17.0)

**Implemented enhancements:**

- \(feat\) map archive flag to repo response [\#141](https://github.com/drone/go-scm/pull/141) ([eoinmcafee00](https://github.com/eoinmcafee00))
- Add the support for delete of the bitbucket file [\#139](https://github.com/drone/go-scm/pull/139) ([DeepakPatankar](https://github.com/DeepakPatankar))

**Fixed bugs:**

- Fix the syntax error of the example code [\#135](https://github.com/drone/go-scm/pull/135) ([LinuxSuRen](https://github.com/LinuxSuRen))

**Closed issues:**

- The deprecation of Bitbucket API endpoint /2.0/teams breaks user registration [\#136](https://github.com/drone/go-scm/issues/136)

**Merged pull requests:**

- release prep for v1.17.0 [\#142](https://github.com/drone/go-scm/pull/142) ([eoinmcafee00](https://github.com/eoinmcafee00))

## [v1.16.3](https://github.com/drone/go-scm/tree/v1.16.3) (2021-12-30)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.16.2...v1.16.3)

**Fixed bugs:**

- fix the deprecation of Bitbucket API endpoint /2.0/teams breaks user registration \(136\) [\#137](https://github.com/drone/go-scm/pull/137) ([eoinmcafee00](https://github.com/eoinmcafee00))

**Closed issues:**

- Any plans to support manage wehook [\#134](https://github.com/drone/go-scm/issues/134)

**Merged pull requests:**

- V1.16.3 [\#138](https://github.com/drone/go-scm/pull/138) ([eoinmcafee00](https://github.com/eoinmcafee00))

## [v1.16.2](https://github.com/drone/go-scm/tree/v1.16.2) (2021-11-30)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.16.1...v1.16.2)

**Merged pull requests:**

- release prep v1.16.2 [\#132](https://github.com/drone/go-scm/pull/132) ([marko-gacesa](https://github.com/marko-gacesa))
- fixbug: gitee webhook parse [\#131](https://github.com/drone/go-scm/pull/131) ([kit101](https://github.com/kit101))

## [v1.16.1](https://github.com/drone/go-scm/tree/v1.16.1) (2021-11-19)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.16.0...v1.16.1)

**Fixed bugs:**

- swap repo and target in bitbucket CompareChanges [\#127](https://github.com/drone/go-scm/pull/127) ([jimsheldon](https://github.com/jimsheldon))

**Merged pull requests:**

- release prep v1.16.1 [\#129](https://github.com/drone/go-scm/pull/129) ([eoinmcafee00](https://github.com/eoinmcafee00))

## [v1.16.0](https://github.com/drone/go-scm/tree/v1.16.0) (2021-11-19)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.15.2...v1.16.0)

**Implemented enhancements:**

- release prep for 1.16.0 [\#128](https://github.com/drone/go-scm/pull/128) ([eoinmcafee00](https://github.com/eoinmcafee00))
- Feat: implemented gitee provider [\#124](https://github.com/drone/go-scm/pull/124) ([kit101](https://github.com/kit101))
- add release & milestone functionality [\#121](https://github.com/drone/go-scm/pull/121) ([eoinmcafee00](https://github.com/eoinmcafee00))

**Fixed bugs:**

- Fix Gitea example code on README.md [\#126](https://github.com/drone/go-scm/pull/126) ([lunny](https://github.com/lunny))

## [v1.15.2](https://github.com/drone/go-scm/tree/v1.15.2) (2021-07-20)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.15.1...v1.15.2)

**Fixed bugs:**

- Fixing Gitea commit API in case `ref/heads/` prefix is added to ref [\#108](https://github.com/drone/go-scm/pull/108) ([Vici37](https://github.com/Vici37))
- use access json header / extend error message parsing for stash [\#89](https://github.com/drone/go-scm/pull/89) ([bakito](https://github.com/bakito))

**Closed issues:**

- Drone and Bitbucket broken for write permission detection for drone build restart permission. [\#87](https://github.com/drone/go-scm/issues/87)

**Merged pull requests:**

- \(maint\) prep for v.1.15.2 release [\#118](https://github.com/drone/go-scm/pull/118) ([tphoney](https://github.com/tphoney))
- Add a vet step to drone config [\#83](https://github.com/drone/go-scm/pull/83) ([tboerger](https://github.com/tboerger))

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 1.15.1
### Added
- (DRON-88) github fix pr ListChanges deleted/renamed status, from [@tphoney](https://github.com/tphoney). See [#113](https://github.com/drone/go-scm/pull/113).
- (DRON-84) github fix pr write permission issue with bitbucket server, from [@eoinmcafee0](https://github.com/eoinmcafee00). See [#114](https://github.com/drone/go-scm/pull/114).

## 1.15.0
### Added
- add delete file for github and gitlab, from [@tphoney](https://github.com/tphoney). See [#110](https://github.com/drone/go-scm/pull/110).

## 1.14.1
### Fixed
- fix gitlab repo encoding in commits from pr request, from [@aradisavljevic](https://github.com/aradisavljevic). See [#109](https://github.com/drone/go-scm/pull/109).

## 1.14.0
### Added
- Added ListCommits in pull request api, from [@aradisavljevic](https://github.com/aradisavljevic). See [#106](https://github.com/drone/go-scm/pull/106).

## 1.13.1
### Fixed
- gitlab, content.find return last_commit_id not commit_id, from [@tphoney](https://github.com/tphoney). See [#104](https://github.com/drone/go-scm/pull/104).

## 1.13.0
### Added
- Create branch functionality, from [@tphoney](https://github.com/tphoney). See [#103](https://github.com/drone/go-scm/pull/103).

## 1.12.0
### Added
- return sha/blob_id for content.list, from [@tphoney](https://github.com/tphoney). See [#102](https://github.com/drone/go-scm/pull/102).

## 1.11.0
### Added
- normalise sha in content, add bitbucket create/update, from [@tphoney](https://github.com/tphoney). See [#101](https://github.com/drone/go-scm/pull/101).

## 1.10.0
### Added
- return hash/object_id for files changed in github, from [@tphoney](https://github.com/tphoney). See [#99](https://github.com/drone/go-scm/pull/99).

## 1.9.0
### Added
- Added issue_comment parsing for github webhook, from [@aman-harness](https://github.com/aman-harness). See [#91](https://github.com/drone/go-scm/pull/91).
- Added Pr in issue, from [@aman-harness](https://github.com/aman-harness). See [#93](https://github.com/drone/go-scm/pull/93).
- gitlab contents. Find returns hash/blob, from [@tphoney](https://github.com/tphoney). See [#97](https://github.com/drone/go-scm/pull/97).
- add ListCommits for gitea and stash, from [@tphoney](https://github.com/tphoney). See [#98](https://github.com/drone/go-scm/pull/98).

### Changed
- retry with event subset for legacy stash versions, from [@bakito](https://github.com/bakito). See [#90](https://github.com/drone/go-scm/pull/90).

## 1.8.0
### Added
- Support for GitLab visibility attribute, from [@bradrydzewski](https://github.com/bradrydzewski). See [79951ad](https://github.com/drone/go-scm/commit/79951ad7a0d0b1989ea84d99be31fcb9320ae348).
- Support for GitHub visibility attribute, from [@bradrydzewski](https://github.com/bradrydzewski). See [5141b8e](https://github.com/drone/go-scm/commit/5141b8e1db921fe2101c12594c5159b9ffffebc3).

### Changed
- Support for parsing unknown pull request events, from [@bradrydzewski](https://github.com/bradrydzewski). See [ffa46d9](https://github.com/drone/go-scm/commit/ffa46d955454baa609975eebbe9fdfc4b0a9f7e9).

## 1.7.2
### Added
- Support for finding and listing repository tags in GitHub driver, from [@chhsia0](https://github.com/chhsia0). See [#79](https://github.com/drone/go-scm/pull/79).
- Support for finding and listing repository tags in Gitea driver, from [@bradyrdzewski](https://github.com/bradyrdzewski). See [427b8a8](https://github.com/drone/go-scm/commit/427b8a85897c892148801824760bc66d3a3cdcdb).
- Support for git object hashes in GitHub, from from [@bradyrdzewski](https://github.com/bradyrdzewski). See [5230330](https://github.com/drone/go-scm/commit/523033025a7ee875fcfb156f4c660b37e269b1a8).
- Support for before and after commit sha in Stash driver, from [@jlehtimaki](https://github.com/jlehtimaki). See [#82](https://github.com/drone/go-scm/pull/82).
- Support for before and after commit sha in GitLab and Bitbucket driver, from [@shubhag](https://github.com/shubhag). See [#85](https://github.com/drone/go-scm/pull/85).

## 1.7.1
### Added
- Support for skip verification in Bitbucket webhook creation, from [@chhsia0](https://github.com/chhsia0). See [#63](https://github.com/drone/go-scm/pull/63).
- Support for Gitea pagination, from [@CirnoT](https://github.com/CirnoT). See [#66](https://github.com/drone/go-scm/pull/66).
- Support for labels in pull request resources, from [@takirala](https://github.com/takirala). See [#67](https://github.com/drone/go-scm/pull/67).
- Support for updating webhooks, from [@chhsia0](https://github.com/chhsia0). See [#71](https://github.com/drone/go-scm/pull/71).

### Fixed
- Populate diff links in pull request resources, from [@shubhag](https://github.com/shubhag). See [#75](https://github.com/drone/go-scm/pull/75).
- Filter Bitbucket repository search by project, from [@bradrydzewski](https://github.com/bradrydzewski).

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


\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
