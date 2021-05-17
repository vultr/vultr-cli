# Change Log

## [v2.5.2](https://github.com/vultr/vultr-cli/compare/v2.5.1..v2.5.2) (2021-05-17)
### Enhancement
* Support config files in $XDG_CONFIG_HOME [PR 153](https://github.com/vultr/vultr-cli/pull/153)

### Documentation 
* Add Arch Linux install instructions [PR 154](https://github.com/vultr/vultr-cli/pull/154)

## [v2.5.1](https://github.com/vultr/vultr-cli/compare/v2.5.0..v2.5.1) (2021-05-12)
### Dependencies
* GoVultr v2.5.0 -> 2.5.1 (fixes issue with backup schedules) [PR 151](https://github.com/vultr/vultr-cli/pull/151)

## [v2.5.0](https://github.com/vultr/vultr-cli/compare/v2.4.1..v2.5.0) (2021-05-06)
### Enhancement
* LoadBalancers : New Features [149](https://github.com/vultr/vultr-cli/pull/149)
  * Ability to attach private networks
  * Ability to set firewalls
  * Get Firewall Rules
  * List Firewall Rules

## [v2.4.0](https://github.com/vultr/vultr-cli/compare/v2.3.0..v2.4.0) (2021-04-01)
### Enhancement
* Add `darwin_arm64` support and builds [PR 143](https://github.com/vultr/vultr-cli/pull/143)

## [v2.3.0](https://github.com/vultr/vultr-cli/compare/v2.2.0..v2.3.0) (2021-02-12)
### Enhancement
* Plans : add `disk count` field [PR 140](https://github.com/vultr/vultr-cli/pull/140)
* BlockStorage : add `mount ID` field [PR 140](https://github.com/vultr/vultr-cli/pull/140)

### Dependencies
* GoVultr v2.3.2 -> 2.4.0 [PR 140](https://github.com/vultr/vultr-cli/pull/140)
* Cobra v1.1.1 -> v1.1.3[PR 139](https://github.com/vultr/vultr-cli/pull/139)

## [v2.2.0](https://github.com/vultr/vultr-cli/compare/v2.1.0..v2.2.0) (2021-01-29)
### Enhancement
* BareMetal : add get command [PR 135](https://github.com/vultr/vultr-cli/pull/135)

### Bug Fixes
* BareMetal : typo in VNC commands [PR 135](https://github.com/vultr/vultr-cli/pull/135)

### Dependencies
* Various dependencies [PR 134](https://github.com/vultr/vultr-cli/pull/134)

## [v2.1.0](https://github.com/vultr/vultr-cli/compare/v2.0.1..v2.1.0) (2020-12-18)
### Enhancement
* Add Bare Metal Start command [127](https://github.com/vultr/vultr-cli/pull/127)
* Update paging information [127](https://github.com/vultr/vultr-cli/pull/127)

### Dependencies
* govultr 2.2.0 -> 2.3.0 [PR 127](https://github.com/vultr/vultr-cli/pull/127)

## [v2.0.1](https://github.com/vultr/vultr-cli/compare/v2.0.0..v2.0.1) (2020-12-08)
### Bug Fixes
* Adding paging support for DNS Records [PR 123](https://github.com/vultr/vultr-cli/pull/123)
* Cleaned up LB output to remove artifacts [PR 125](https://github.com/vultr/vultr-cli/pull/125)

### Dependencies
* govultr 2.0.0 -> 2.2.0 [PR 125](https://github.com/vultr/vultr-cli/pull/125)

## [v2.0.0](https://github.com/vultr/vultr-cli/compare/v1.0.0..v2.0.0) (2020-11-24)
### Enhancement
* Vultr-CLI v2.0.0 release 

### Changes
* Vultr-CLI v2.0.0 is running on API v2
* Server has been renamed to Instance to match with API v2

## [v1.0.0](https://github.com/vultr/vultr-cli/compare/v0.4.0..v1.0.0) (2020-11-19)
### Enhancement
* Vultr-CLI v1.0.0 release [PR 114](https://github.com/vultr/vultr-cli/pull/114) 

## [v0.4.0](https://github.com/vultr/vultr-cli/compare/v0.3.2..v0.4.0) (2020-09-03)
### Enhancement
* Improve error responses by adding a newline [PR 109](https://github.com/vultr/vultr-cli/pull/109) 
* Add Server User Data subcommands Get and Set [PR 105](https://github.com/vultr/vultr-cli/pull/105)

### Dependencies
* spf13/viper from 1.7.0 -> 1.7.1 [PR 104](https://github.com/vultr/vultr-cli/pull/104) [PR 108](https://github.com/vultr/vultr-cli/pull/108)

## [v0.3.2](https://github.com/vultr/vultr-cli/compare/v0.3.1..v0.3.2) (2020-06-04)
### Dependencies
* govultr 0.4.1 -> 0.4.2 [PR 98](https://github.com/vultr/vultr-cli/pull/98)

## [v0.3.1](https://github.com/vultr/vultr-cli/compare/v0.3.0..v0.3.1) (2020-05-18)
### Enhancement
* spf13/viper 1.6.3 -> 1.7.0 [PR 94](https://github.com/vultr/vultr-cli/pull/94)
* govultr 0.3.2 -> 0.4.1 [PR 95](https://github.com/vultr/vultr-cli/pull/95)


## [v0.3.0](https://github.com/vultr/vultr-cli/compare/v0.2.1..v0.3.0) (2020-04-14)
### Enhancement
* OS no longer required for App/Iso/Snapshot during server create [PR 88](https://github.com/vultr/vultr-cli/pull/88)
* Added in missing newline characters [PR 89](https://github.com/vultr/vultr-cli/pull/89)

### Dependencies
* spf13/viper 1.6.2 -> 1.6.3 [PR 87](https://github.com/vultr/vultr-cli/pull/87)
* spf13/cobra 0.0.6 -> 0.0.7 [PR 85](https://github.com/vultr/vultr-cli/pull/85)

## [v0.2.1](https://github.com/vultr/vultr-cli/compare/v0.2.0..v0.2.1) (2020-03-18)
### Dependencies
* govultr 0.3.0 -> 0.3.1 [PR 81](https://github.com/vultr/vultr-cli/pull/81)

## [v0.2.0](https://github.com/vultr/vultr-cli/compare/v0.1.11..v0.2.0) (2020-03-11)
### Enhancement
* Object Storage support [PR 79](https://github.com/vultr/vultr-cli/pull/79) [74](https://github.com/vultr/vultr-cli/pull/74)

### Bug Fix
* Server Create description [PR 75](https://github.com/vultr/vultr-cli/pull/75)

### Dependencies
* spf13/viper 1.5.0 -> 1.6.2 [PR 76](https://github.com/vultr/vultr-cli/pull/76)
* spf13/cobra 0.0.5 -> 0.0.6 [PR 78](https://github.com/vultr/vultr-cli/pull/78)
* govultr 0.1.7 -> 0.3.0 [PR 77](https://github.com/vultr/vultr-cli/pull/77)

## [v0.1.11](https://github.com/vultr/vultr-cli/compare/v0.1.10..v0.1.11) (2019-12-09)
### Bug fix
* Fix error message on network create [PR 65](https://github.com/vultr/vultr-cli/pull/65)

## [v0.1.10](https://github.com/vultr/vultr-cli/compare/v0.1.9..v0.1.10) (2019-11-12)
### Bug fix
* GoVultr v0.1.7 version fix [PR 61](https://github.com/vultr/vultr-cli/pull/61)

## [v0.1.9](https://github.com/vultr/vultr-cli/compare/v0.1.8..v0.1.9) (2019-11-11)
### Enhancements
* Updating dependency versions  [PR 59](https://github.com/vultr/vultr-cli/pull/59)
* GoVultr v0.1.6 now supports retry [PR 59](https://github.com/vultr/vultr-cli/pull/59)

## [v0.1.8](https://github.com/vultr/vultr-cli/compare/v0.1.7..v0.1.8) (2019-10-17)
### Bug fix
* Fix for goreleaser to release homebrew tap  [PR 57](https://github.com/vultr/vultr-cli/pull/57)

## [v0.1.7](https://github.com/vultr/vultr-cli/compare/v0.1.6..v0.1.7) (2019-10-17)
### Enhancements
* Bump GoVultr to v0.1.5 [PR 55](https://github.com/vultr/vultr-cli/pull/55)

  
## [v0.1.6](https://github.com/vultr/vultr-cli/compare/v0.1.5..v0.1.6) (2019-09-04)
### Enhancements
* Print the original API error messages in [PR 50](https://github.com/vultr/vultr-cli/pull/50) && [PR 52](https://github.com/vultr/vultr-cli/pull/52)
  * application
  * plans
  * regions
  * user

## [v0.1.5](https://github.com/vultr/vultr-cli/compare/v0.1.4..v0.1.5) (2019-09-03)
### Enhancements
* Add contextual instructions for Vultr API key setup [PR 47](https://github.com/vultr/vultr-cli/pull/47)

## [v0.1.4](https://github.com/vultr/vultr-cli/compare/v0.1.3..v0.1.4) (2019-08-26)
### Enhancements
* Makefile entry for gofmt [PR 44](https://github.com/vultr/vultr-cli/pull/44)
* New command `script contents` will display contents of a given script  [PR 43](https://github.com/vultr/vultr-cli/pull/43)

## [v0.1.3](https://github.com/vultr/vultr-cli/compare/v0.1.2..v0.1.3) (2019-08-21)
### Bug Fixes
* Quote handling on DNS Record Data [PR #41](https://github.com/vultr/vultr-cli/pull/41)
  
## [v0.1.2](https://github.com/vultr/vultr-cli/compare/v0.1.1..v0.1.2) (2019-07-15)
### Dependencies
* Updated dependencies [PR #35](https://github.com/vultr/vultr-cli/pull/35)
  * Govultr `v0.1.3` -> `v0.1.4`
  * Cobra `v0.0.4` -> `v0.0.5`
* Added vendor folder [PR #35](https://github.com/vultr/vultr-cli/pull/35)

## [v0.1.1](https://github.com/vultr/vultr-cli/compare/v0.1.0..v0.1.1) (2019-07-08)
### Enhancements
* Added `destroy` alias for all `delete` commands [PR #30](https://github.com/vultr/vultr-cli/pull/30)
* Added `description` field for `snapshot` command output [PR #29](https://github.com/vultr/vultr-cli/pull/29)
* Added GoReleaser to handle tagged releases [PR #31](https://github.com/vultr/vultr-cli/pull/31)
* Updated makefile to strip out GOPATH in build process [PR #28](https://github.com/vultr/vultr-cli/pull/28)
* Typo fixes

## v0.1.0 (2019-06-24)
### Features
* Initial release
