# Change Log

## GoVultr v1 changelog is located [here](https://github.com/vultr/govultr/blob/v1/CHANGELOG.md)

## [v2.2.0](https://github.com/vultr/govultr/compare/v2.1.0..v2.2.0) (2020-12-07)
### Breaking Change
* All bools have been updated to pointers to avoid issues where false values not being sent in request. Thanks @Static-Flow [115](https://github.com/vultr/govultr/pull/115)

## [v2.1.0](https://github.com/vultr/govultr/compare/v2.0.0..v2.1.0) (2020-11-30)
### Bug fixes
* ReservedIP - Attach call creates proper json. [112](https://github.com/vultr/govultr/pull/112)
* User - APIEnabled takes pointer of bool [112](https://github.com/vultr/govultr/pull/112)

## v2.0.0 (2020-11-20)
### Initial Release
* GoVultr v2.0.0 Release - Uses Vultr API v2.
* GoVultr v1.0.0 is now on [branch v1](https://github.com/vultr/govultr/tree/v1)