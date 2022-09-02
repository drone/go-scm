# Changelog

## [v1.15.2](https://github.com/drone/go-scm/tree/v1.15.2) (2022-09-02)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.27.0...v1.15.2)

**Implemented enhancements:**

- Added support for branch in list commits bb onprem API [\#215](https://github.com/drone/go-scm/pull/215) ([mohitg0795](https://github.com/mohitg0795))

**Closed issues:**

- file naming conventions [\#208](https://github.com/drone/go-scm/issues/208)
- Support for Azure Devops git repos? [\#53](https://github.com/drone/go-scm/issues/53)

**Merged pull requests:**

- \[PL-26239\]: fix for list response [\#218](https://github.com/drone/go-scm/pull/218) ([bhavya181](https://github.com/bhavya181))
- \[PL-26239\]: added api to list installation for github app [\#213](https://github.com/drone/go-scm/pull/213) ([bhavya181](https://github.com/bhavya181))
- \(maint\) fixing naming and add more go best practice [\#211](https://github.com/drone/go-scm/pull/211) ([tphoney](https://github.com/tphoney))

## [v1.27.0](https://github.com/drone/go-scm/tree/v1.27.0) (2022-07-19)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.26.0...v1.27.0)

**Merged pull requests:**

- Update scm version 1.27.0 [\#206](https://github.com/drone/go-scm/pull/206) ([raghavharness](https://github.com/raghavharness))
- Using resource version 2.0 for Azure [\#205](https://github.com/drone/go-scm/pull/205) ([raghavharness](https://github.com/raghavharness))

## [v1.26.0](https://github.com/drone/go-scm/tree/v1.26.0) (2022-07-01)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.25.0...v1.26.0)

**Implemented enhancements:**

- Support parsing PR comment events for Bitbucket Cloud [\#202](https://github.com/drone/go-scm/pull/202) ([rutvijmehta-harness](https://github.com/rutvijmehta-harness))
- added issue comment hook support for Azure [\#200](https://github.com/drone/go-scm/pull/200) ([raghavharness](https://github.com/raghavharness))

**Fixed bugs:**

- \[CI-4623\] - Azure webhook parseAPI changes [\#198](https://github.com/drone/go-scm/pull/198) ([raghavharness](https://github.com/raghavharness))

**Merged pull requests:**

- Fixed formatting in README.md [\#199](https://github.com/drone/go-scm/pull/199) ([hemanthmantri](https://github.com/hemanthmantri))

## [v1.25.0](https://github.com/drone/go-scm/tree/v1.25.0) (2022-06-16)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.24.0...v1.25.0)

**Implemented enhancements:**

- Support parsing Gitlab Note Hook event [\#194](https://github.com/drone/go-scm/pull/194) ([rutvijmehta-harness](https://github.com/rutvijmehta-harness))

**Fixed bugs:**

- \[PL-25889\]: fix list branches Azure API [\#195](https://github.com/drone/go-scm/pull/195) ([bhavya181](https://github.com/bhavya181))
- Return project specific hooks only in ListHooks API for Azure. [\#192](https://github.com/drone/go-scm/pull/192) ([raghavharness](https://github.com/raghavharness))

**Merged pull requests:**

- Update scm version 1.25.0 [\#197](https://github.com/drone/go-scm/pull/197) ([rutvijmehta-harness](https://github.com/rutvijmehta-harness))

## [v1.24.0](https://github.com/drone/go-scm/tree/v1.24.0) (2022-06-07)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.23.0...v1.24.0)

**Implemented enhancements:**

- Added PR find and listCommit API support for Azure [\#188](https://github.com/drone/go-scm/pull/188) ([raghavharness](https://github.com/raghavharness))

**Fixed bugs:**

- remove redundant slash from list commits api [\#190](https://github.com/drone/go-scm/pull/190) ([aman-harness](https://github.com/aman-harness))
- Using target commit instead of source in base info for azure [\#189](https://github.com/drone/go-scm/pull/189) ([raghavharness](https://github.com/raghavharness))

**Closed issues:**

- gitee client pagination bug [\#187](https://github.com/drone/go-scm/issues/187)

**Merged pull requests:**

- release\_prep\_v1.24.0 [\#191](https://github.com/drone/go-scm/pull/191) ([tphoney](https://github.com/tphoney))

## [v1.23.0](https://github.com/drone/go-scm/tree/v1.23.0) (2022-05-23)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.22.0...v1.23.0)

**Implemented enhancements:**

- Add the support to fetch commit of a particular file [\#182](https://github.com/drone/go-scm/pull/182) ([DeepakPatankar](https://github.com/DeepakPatankar))

**Fixed bugs:**

- Remove the null value de-reference issue when the bitbucket server url is nil [\#183](https://github.com/drone/go-scm/pull/183) ([DeepakPatankar](https://github.com/DeepakPatankar))
- \[PL-24913\]: Handle the error raised while creating a multipart input [\#181](https://github.com/drone/go-scm/pull/181) ([DeepakPatankar](https://github.com/DeepakPatankar))

**Merged pull requests:**

- Upgrade the scm version [\#185](https://github.com/drone/go-scm/pull/185) ([DeepakPatankar](https://github.com/DeepakPatankar))

## [v1.22.0](https://github.com/drone/go-scm/tree/v1.22.0) (2022-05-10)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.21.1...v1.22.0)

**Implemented enhancements:**

- \[feat\]: \[PL-24913\]: Add the support for create and update file in bitbucket server [\#177](https://github.com/drone/go-scm/pull/177) ([DeepakPatankar](https://github.com/DeepakPatankar))
- \[feat\]: \[PL-24915\]: Add the support for create branches in bitbucket server [\#174](https://github.com/drone/go-scm/pull/174) ([DeepakPatankar](https://github.com/DeepakPatankar))
- \[PL-24911\]: Make project name as optional param in Azure Repo APIs [\#173](https://github.com/drone/go-scm/pull/173) ([DeepakPatankar](https://github.com/DeepakPatankar))

**Fixed bugs:**

- \[feat\]: \[PL-25025\]: Updated Project validation for Azure API [\#179](https://github.com/drone/go-scm/pull/179) ([mankrit-singh](https://github.com/mankrit-singh))
- \[fix\]: \[PL-24880\]: Trim the ref when fetching default branch in get Repo API [\#172](https://github.com/drone/go-scm/pull/172) ([DeepakPatankar](https://github.com/DeepakPatankar))
- fixbug: gitee populatePageValues [\#167](https://github.com/drone/go-scm/pull/167) ([kit101](https://github.com/kit101))

**Closed issues:**

- gitea find commit   [\#125](https://github.com/drone/go-scm/issues/125)

**Merged pull requests:**

- \[feat\]: \[PL-25025\]: Changelog Updated/New Version [\#180](https://github.com/drone/go-scm/pull/180) ([mankrit-singh](https://github.com/mankrit-singh))

## [v1.21.1](https://github.com/drone/go-scm/tree/v1.21.1) (2022-04-22)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.21.0...v1.21.1)

**Fixed bugs:**

- remove double invocation of convertRepository [\#170](https://github.com/drone/go-scm/pull/170) ([d1wilko](https://github.com/d1wilko))

**Merged pull requests:**

- \(maint\) release prep for 1.21.1 [\#171](https://github.com/drone/go-scm/pull/171) ([d1wilko](https://github.com/d1wilko))

## [v1.21.0](https://github.com/drone/go-scm/tree/v1.21.0) (2022-04-22)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.20.0...v1.21.0)

**Implemented enhancements:**

- Add support for repository find in azure [\#164](https://github.com/drone/go-scm/pull/164) ([goelsatyam2](https://github.com/goelsatyam2))
- \(feat\) add azure webhook parsing, creation deletion & list [\#163](https://github.com/drone/go-scm/pull/163) ([tphoney](https://github.com/tphoney))
- \(DRON-242\) azure add compare commits,get commit,list repos [\#162](https://github.com/drone/go-scm/pull/162) ([tphoney](https://github.com/tphoney))

**Fixed bugs:**

- \(fix\) handle nil repos in github responses [\#168](https://github.com/drone/go-scm/pull/168) ([tphoney](https://github.com/tphoney))

**Closed issues:**

- When attempting to clone my git repo from GitLab drone hangs on git fetch. [\#161](https://github.com/drone/go-scm/issues/161)
- Fix dump response [\#119](https://github.com/drone/go-scm/issues/119)

**Merged pull requests:**

- \(maint\) release prep for 1.21.0 [\#169](https://github.com/drone/go-scm/pull/169) ([d1wilko](https://github.com/d1wilko))

## [v1.20.0](https://github.com/drone/go-scm/tree/v1.20.0) (2022-03-08)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.19.1...v1.20.0)

**Implemented enhancements:**

- \(DRON-242\) initial implementation of azure devops support [\#158](https://github.com/drone/go-scm/pull/158) ([tphoney](https://github.com/tphoney))

**Fixed bugs:**

- fixed raw response dumping in client [\#159](https://github.com/drone/go-scm/pull/159) ([marko-gacesa](https://github.com/marko-gacesa))

**Merged pull requests:**

- \(maint\) release prep for 1.20.0 [\#160](https://github.com/drone/go-scm/pull/160) ([d1wilko](https://github.com/d1wilko))

## [v1.19.1](https://github.com/drone/go-scm/tree/v1.19.1) (2022-02-23)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.19.0...v1.19.1)

**Fixed bugs:**

- Bitbucket list files fix [\#154](https://github.com/drone/go-scm/pull/154) ([mohitg0795](https://github.com/mohitg0795))
- GitHub list commits fix [\#152](https://github.com/drone/go-scm/pull/152) ([mohitg0795](https://github.com/mohitg0795))
- Bitbucket compare changes fix for rename and removed file ops [\#151](https://github.com/drone/go-scm/pull/151) ([mohitg0795](https://github.com/mohitg0795))

**Merged pull requests:**

- prep for v1.19.1 [\#155](https://github.com/drone/go-scm/pull/155) ([tphoney](https://github.com/tphoney))

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

## [v1.15.1](https://github.com/drone/go-scm/tree/v1.15.1) (2021-06-17)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.15.0...v1.15.1)

**Merged pull requests:**

- update changelog for 1.15.1 [\#115](https://github.com/drone/go-scm/pull/115) ([eoinmcafee00](https://github.com/eoinmcafee00))
- fix issue where write permission on bitbucket server wasn't working a… [\#114](https://github.com/drone/go-scm/pull/114) ([eoinmcafee00](https://github.com/eoinmcafee00))
- \(DRON-88\) github fix pr ListChanges deleted/renamed status [\#113](https://github.com/drone/go-scm/pull/113) ([tphoney](https://github.com/tphoney))

## [v1.15.0](https://github.com/drone/go-scm/tree/v1.15.0) (2021-05-27)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.14.1...v1.15.0)

**Merged pull requests:**

- \(feat\) add delete file for github and gitlab [\#110](https://github.com/drone/go-scm/pull/110) ([tphoney](https://github.com/tphoney))

## [v1.14.1](https://github.com/drone/go-scm/tree/v1.14.1) (2021-05-19)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.14.0...v1.14.1)

**Merged pull requests:**

- Integration tests for pr commits [\#109](https://github.com/drone/go-scm/pull/109) ([aradisavljevic](https://github.com/aradisavljevic))

## [v1.14.0](https://github.com/drone/go-scm/tree/v1.14.0) (2021-05-12)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.13.1...v1.14.0)

**Closed issues:**

- Gitea API ref causing build failure [\#105](https://github.com/drone/go-scm/issues/105)

**Merged pull requests:**

- Update changelog for 1.14.0 [\#107](https://github.com/drone/go-scm/pull/107) ([tphoney](https://github.com/tphoney))
- Added api call to get commit details list for given pull request [\#106](https://github.com/drone/go-scm/pull/106) ([aradisavljevic](https://github.com/aradisavljevic))

## [v1.13.1](https://github.com/drone/go-scm/tree/v1.13.1) (2021-04-15)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.13.0...v1.13.1)

**Merged pull requests:**

- \(fix\) gitlab, content.find return last\_commit\_id not commit\_id [\#104](https://github.com/drone/go-scm/pull/104) ([tphoney](https://github.com/tphoney))

## [v1.13.0](https://github.com/drone/go-scm/tree/v1.13.0) (2021-04-14)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.12.0...v1.13.0)

**Merged pull requests:**

- \(feat\) add create branch functionality [\#103](https://github.com/drone/go-scm/pull/103) ([tphoney](https://github.com/tphoney))

## [v1.12.0](https://github.com/drone/go-scm/tree/v1.12.0) (2021-04-09)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.11.0...v1.12.0)

**Merged pull requests:**

- \(feat\) return sha/blob\_id for content.list [\#102](https://github.com/drone/go-scm/pull/102) ([tphoney](https://github.com/tphoney))

## [v1.11.0](https://github.com/drone/go-scm/tree/v1.11.0) (2021-04-07)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.10.0...v1.11.0)

**Merged pull requests:**

- \(feat\) normalise sha in content, add bitbucket create/update [\#101](https://github.com/drone/go-scm/pull/101) ([tphoney](https://github.com/tphoney))

## [v1.10.0](https://github.com/drone/go-scm/tree/v1.10.0) (2021-03-23)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.9.0...v1.10.0)

**Merged pull requests:**

- \(change\) return hash/object\_id for files changed in github. [\#99](https://github.com/drone/go-scm/pull/99) ([tphoney](https://github.com/tphoney))

## [v1.9.0](https://github.com/drone/go-scm/tree/v1.9.0) (2021-03-23)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.8.0...v1.9.0)

**Closed issues:**

- the event pr:from\_ref\_updated is unknown with Bitbucket Server v6.10.5 [\#88](https://github.com/drone/go-scm/issues/88)

**Merged pull requests:**

- \(feat\) add ListCommits for gitea and stash [\#98](https://github.com/drone/go-scm/pull/98) ([tphoney](https://github.com/tphoney))
- \(feat\) gitlab contents.Find returns hash/blob [\#97](https://github.com/drone/go-scm/pull/97) ([tphoney](https://github.com/tphoney))
- \[CI-0\]: Added Pr in issue [\#93](https://github.com/drone/go-scm/pull/93) ([aman-harness](https://github.com/aman-harness))
- ENH: Added issue\_comment parsing for github webhook [\#91](https://github.com/drone/go-scm/pull/91) ([aman-harness](https://github.com/aman-harness))
- retry with event subset for legacy stash versions [\#90](https://github.com/drone/go-scm/pull/90) ([bakito](https://github.com/bakito))

## [v1.8.0](https://github.com/drone/go-scm/tree/v1.8.0) (2020-12-17)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.7.2...v1.8.0)

**Closed issues:**

- Failure when PR labels are queried on Bitbucket Cloud [\#86](https://github.com/drone/go-scm/issues/86)

## [v1.7.2](https://github.com/drone/go-scm/tree/v1.7.2) (2020-11-30)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.7.1...v1.7.2)

**Closed issues:**

- Adding code for files changed in the commit range? [\#80](https://github.com/drone/go-scm/issues/80)

**Merged pull requests:**

- Set before & after attribute for gitlab & bitbucket push event [\#85](https://github.com/drone/go-scm/pull/85) ([shubham149](https://github.com/shubham149))
- Added before and after values to show before and after commitids on push [\#82](https://github.com/drone/go-scm/pull/82) ([jlehtimaki](https://github.com/jlehtimaki))
- Implemented FindTag for github driver. [\#79](https://github.com/drone/go-scm/pull/79) ([chhsia0](https://github.com/chhsia0))

## [v1.7.1](https://github.com/drone/go-scm/tree/v1.7.1) (2020-09-02)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.7.0...v1.7.1)

**Implemented enhancements:**

- Support fetching labels from pull requests in github,gitlab & gitea [\#67](https://github.com/drone/go-scm/pull/67) ([takirala](https://github.com/takirala))

**Closed issues:**

- paginate Gitea repository list [\#64](https://github.com/drone/go-scm/issues/64)
- support for create and update content? [\#55](https://github.com/drone/go-scm/issues/55)

**Merged pull requests:**

- Fix github pull request link from diff url to html url [\#75](https://github.com/drone/go-scm/pull/75) ([shubhag](https://github.com/shubhag))
- Add support for commit list in push webhook [\#74](https://github.com/drone/go-scm/pull/74) ([shubhag](https://github.com/shubhag))
- Fixed github integration tests. [\#72](https://github.com/drone/go-scm/pull/72) ([chhsia0](https://github.com/chhsia0))
- Added UpdateHook method to RepositoryService. [\#71](https://github.com/drone/go-scm/pull/71) ([chhsia0](https://github.com/chhsia0))
- Fix pagination being ignored for Gitea [\#66](https://github.com/drone/go-scm/pull/66) ([CirnoT](https://github.com/CirnoT))
- Supported `skip_cert_verification` for bitbucket webhook registration. [\#63](https://github.com/drone/go-scm/pull/63) ([chhsia0](https://github.com/chhsia0))
- Support PR creation. [\#60](https://github.com/drone/go-scm/pull/60) ([chhsia0](https://github.com/chhsia0))
- Cleaned up content service and added missing tests. [\#42](https://github.com/drone/go-scm/pull/42) ([chhsia0](https://github.com/chhsia0))

## [v1.7.0](https://github.com/drone/go-scm/tree/v1.7.0) (2020-05-20)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.6.0...v1.7.0)

**Fixed bugs:**

- Bitbucket Build Status Empty [\#27](https://github.com/drone/go-scm/issues/27)

**Closed issues:**

- what is idea about creating the Request with header of `"Content-Type": "application/x-www-form-urlencoded"`? [\#56](https://github.com/drone/go-scm/issues/56)
- Gitlab push with multiple commits \(wrong information\) [\#45](https://github.com/drone/go-scm/issues/45)
- Gitlab Hook sets incorrect value for parameter for skipping ssl verification [\#39](https://github.com/drone/go-scm/issues/39)
- Support for listing the contents of a directory [\#36](https://github.com/drone/go-scm/issues/36)
- Stash driver doesn't work for domain \(@\) users [\#31](https://github.com/drone/go-scm/issues/31)

**Merged pull requests:**

- Ensure git ref is URI encoded \(%2F\) for Gitea call [\#61](https://github.com/drone/go-scm/pull/61) ([CirnoT](https://github.com/CirnoT))
- Set `Repository.Link` in various SCM drivers. [\#58](https://github.com/drone/go-scm/pull/58) ([chhsia0](https://github.com/chhsia0))
- Change gitlab force\_remove\_source\_branch to bool [\#57](https://github.com/drone/go-scm/pull/57) ([jezhodges](https://github.com/jezhodges))
- Support content create and update [\#54](https://github.com/drone/go-scm/pull/54) ([soulseen](https://github.com/soulseen))
- Fix gitea tag link [\#52](https://github.com/drone/go-scm/pull/52) ([lunny](https://github.com/lunny))
- gitea: added "Before" field with the last commit before the push is made [\#51](https://github.com/drone/go-scm/pull/51) ([daniel-meister](https://github.com/daniel-meister))
- Fix Bitbucket link handling for branches with slashes [\#47](https://github.com/drone/go-scm/pull/47) ([christianruhstaller](https://github.com/christianruhstaller))
- Fix Gitlab webhook commit details [\#46](https://github.com/drone/go-scm/pull/46) ([marcotuna](https://github.com/marcotuna))
- scm/driver/github: allow SkipVerify [\#44](https://github.com/drone/go-scm/pull/44) ([gpaul](https://github.com/gpaul))
- Implemented `List` method for the content service. [\#41](https://github.com/drone/go-scm/pull/41) ([chhsia0](https://github.com/chhsia0))
- fixes gitlab hook ssl verification [\#40](https://github.com/drone/go-scm/pull/40) ([ConorNevin](https://github.com/ConorNevin))
- Added `CompareChanges` function to `GitService`. [\#38](https://github.com/drone/go-scm/pull/38) ([chhsia0](https://github.com/chhsia0))

## [v1.6.0](https://github.com/drone/go-scm/tree/v1.6.0) (2019-09-20)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.5.0...v1.6.0)

**Closed issues:**

- Bitbucket Author Username empty because of API changes [\#28](https://github.com/drone/go-scm/issues/28)
- Track Before Sha for Bitbucket Server [\#23](https://github.com/drone/go-scm/issues/23)
- Enable pluggable drivers via RPC [\#4](https://github.com/drone/go-scm/issues/4)

**Merged pull requests:**

- Support SHA of Git Tag from Gogs Webhook [\#37](https://github.com/drone/go-scm/pull/37) ([marcotuna](https://github.com/marcotuna))
- bitbucket: added "Before" field with the last commit before the push … [\#32](https://github.com/drone/go-scm/pull/32) ([jkdev81](https://github.com/jkdev81))

## [v1.5.0](https://github.com/drone/go-scm/tree/v1.5.0) (2019-06-06)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.4.0...v1.5.0)

**Closed issues:**

- Use StdEncoding to decode GitLab files [\#21](https://github.com/drone/go-scm/issues/21)
- include base\_ref in push hooks [\#9](https://github.com/drone/go-scm/issues/9)
- Update GitLab Status API Endpoint [\#6](https://github.com/drone/go-scm/issues/6)

**Merged pull requests:**

- Validate webhook using signature header in Gitea [\#24](https://github.com/drone/go-scm/pull/24) ([techknowlogick](https://github.com/techknowlogick))
- get sha of tag in gitea [\#22](https://github.com/drone/go-scm/pull/22) ([techknowlogick](https://github.com/techknowlogick))

## [v1.4.0](https://github.com/drone/go-scm/tree/v1.4.0) (2019-04-16)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.3.0...v1.4.0)

**Closed issues:**

- implement commit by sha endpoint for gogs [\#7](https://github.com/drone/go-scm/issues/7)

## [v1.3.0](https://github.com/drone/go-scm/tree/v1.3.0) (2019-04-10)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.2.0...v1.3.0)

**Closed issues:**

- Gitea: No avatar, full name, or commit message set for cron builds [\#18](https://github.com/drone/go-scm/issues/18)

**Merged pull requests:**

- Fetch specific commit by sha for gogs, closes \#7 [\#20](https://github.com/drone/go-scm/pull/20) ([ogarcia](https://github.com/ogarcia))
- Fix Gitea commit info \(fixes \#18\) [\#19](https://github.com/drone/go-scm/pull/19) ([pgaskin](https://github.com/pgaskin))

## [v1.2.0](https://github.com/drone/go-scm/tree/v1.2.0) (2019-02-24)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.1.0...v1.2.0)

**Closed issues:**

- Rebuild Drone 1.0.0 RC6\(?\) [\#16](https://github.com/drone/go-scm/issues/16)
- Bitbucket Server \(Stash\) Write Access [\#15](https://github.com/drone/go-scm/issues/15)

## [v1.1.0](https://github.com/drone/go-scm/tree/v1.1.0) (2019-02-22)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.0.9...v1.1.0)

## [v1.0.9](https://github.com/drone/go-scm/tree/v1.0.9) (2019-02-15)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.0.8...v1.0.9)

## [v1.0.8](https://github.com/drone/go-scm/tree/v1.0.8) (2019-02-15)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.0.7...v1.0.8)

## [v1.0.7](https://github.com/drone/go-scm/tree/v1.0.7) (2019-02-12)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.0.6...v1.0.7)

## [v1.0.6](https://github.com/drone/go-scm/tree/v1.0.6) (2019-02-07)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.0.5...v1.0.6)

**Closed issues:**

- DRONE\_COMMIT\_BRANCH refers to target, not source branch in merge requests [\#11](https://github.com/drone/go-scm/issues/11)
- use head\_commit for GitHub Tag if possible [\#8](https://github.com/drone/go-scm/issues/8)

**Merged pull requests:**

- Fetch specific commit for Gitea [\#10](https://github.com/drone/go-scm/pull/10) ([techknowlogick](https://github.com/techknowlogick))

## [v1.0.5](https://github.com/drone/go-scm/tree/v1.0.5) (2018-12-13)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.0.4...v1.0.5)

**Closed issues:**

- State not set for Gitea commit status [\#5](https://github.com/drone/go-scm/issues/5)

**Merged pull requests:**

- Paginate by URL [\#3](https://github.com/drone/go-scm/pull/3) ([nathan-fps](https://github.com/nathan-fps))

## [v1.0.4](https://github.com/drone/go-scm/tree/v1.0.4) (2018-11-13)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.0.3...v1.0.4)

## [v1.0.3](https://github.com/drone/go-scm/tree/v1.0.3) (2018-11-04)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.0.2...v1.0.3)

## [v1.0.2](https://github.com/drone/go-scm/tree/v1.0.2) (2018-11-04)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.0.1...v1.0.2)

**Merged pull requests:**

- Fix out-of-bounds error within Gitea integration [\#2](https://github.com/drone/go-scm/pull/2) ([tboerger](https://github.com/tboerger))
- Use correct driver name for Gitea [\#1](https://github.com/drone/go-scm/pull/1) ([tboerger](https://github.com/tboerger))

## [v1.0.1](https://github.com/drone/go-scm/tree/v1.0.1) (2018-08-30)

[Full Changelog](https://github.com/drone/go-scm/compare/v1.0.0...v1.0.1)

## [v1.0.0](https://github.com/drone/go-scm/tree/v1.0.0) (2018-08-25)

[Full Changelog](https://github.com/drone/go-scm/compare/6c26457e9596c5f82726624adb6a4ab1b5d82376...v1.0.0)



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
