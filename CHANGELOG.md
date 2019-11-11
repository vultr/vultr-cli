# Change Log

## [v0.1.9](https://github.com/vultr/vultr-cli/compare/v0.1.8..v0.1.9) (2019-11-11)
### Bug fix
* Updating dependency versions  [PR 57](https://github.com/vultr/vultr-cli/pull/59)
* GoVultr v0.1.6 now supports retry [PR 57](https://github.com/vultr/vultr-cli/pull/59)

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
