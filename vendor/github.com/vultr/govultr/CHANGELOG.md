# Change Log

## [v0.4.2](https://github.com/vultr/govultr/compare/v0.4.1..v0.4.2) (2020-06-02)
### Bug Fix
*  DNS Records: Allow DNS Records on updates to have priority of 0 [#67](https://github.com/vultr/govultr/pull/67)

## [v0.4.1](https://github.com/vultr/govultr/compare/v0.4.0..v0.4.1) (2020-05-08)
### Bug Fix
*  LoadBalancers: Fix nil pointer in create call [#65](https://github.com/vultr/govultr/pull/65)

## [v0.4.0](https://github.com/vultr/govultr/compare/v0.3.3..v0.4.0) (2020-04-29)
### Enhancement
*  LoadBalancers: Proxy protocol is now available as an option [#62](https://github.com/vultr/govultr/pull/62)
*  LoadBalancers: Ability to attach instances during Create (note this adds a new param to create call) [#62](https://github.com/vultr/govultr/pull/62)

### Bug Fix
*  LoadBalancers: Create call will now properly pass your algorithm [#61](https://github.com/vultr/govultr/pull/61)

### CI/CD
* Drop go 1.12 and add 1.14 [#63](https://github.com/vultr/govultr/pull/63)

## [v0.3.3](https://github.com/vultr/govultr/compare/v0.3.2..v0.3.3) (2020-04-15)
### Dependencies
*  go-retryablehttp 0.6.4 -> 0.6.6 [#58](https://github.com/vultr/govultr/pull/58)

## [v0.3.2](https://github.com/vultr/govultr/compare/v0.3.1..v0.3.2) (2020-03-25)
### Enhancement
*  Added support to live attach/detach blockstorage [#55](https://github.com/vultr/govultr/pull/55)

## [v0.3.1](https://github.com/vultr/govultr/compare/v0.3.0..v0.3.1) (2020-03-11)
### Enhancement
*  Added support for Load Balancers SSL Calls [#53](https://github.com/vultr/govultr/pull/53)
### Bug Fix
* Fixed InstanceList type from string to int [#50](https://github.com/vultr/govultr/pull/50)
* Fixed struct json marshalling issue [#48](https://github.com/vultr/govultr/pull/48)
### Dependencies
* go-retryablehttp 0.6.3 -> 0.6.4 [#51](https://github.com/vultr/govultr/pull/51)

## [v0.3.0](https://github.com/vultr/govultr/compare/v0.2.0..v0.3.0) (2020-02-24)
### Enhancement
*  Added support for Load Balancers [#44](https://github.com/vultr/govultr/pull/44)
### Bug Fix
* Fixed Object Storage Get call [#46](https://github.com/vultr/govultr/pull/46)

## [v0.2.0](https://github.com/vultr/govultr/compare/v0.1.7..v0.2.0) (2020-02-06)
### Enhancement
*  Added support for Object Storage [#39](https://github.com/vultr/govultr/pull/39)

## [v0.1.7](https://github.com/vultr/govultr/compare/v0.1.6..v0.1.7) (2019-11-11)
### Enhancement
*  Version number was missing in v0.1.6 - Attempt was made to fix however it will not work. Cutting new release to remedy this.

## [v0.1.6](https://github.com/vultr/govultr/compare/v0.1.5..v0.1.6) (2019-11-07)
### Enhancement
*  Retry rate-limited requests with exponential backoff[#28](https://github.com/vultr/govultr/pull/28)

## [v0.1.5](https://github.com/vultr/govultr/compare/v0.1.4..v0.1.5) (2019-10-16)
### Enhancement
*  Whitelisting public endpoints that do not require the api key[#24](https://github.com/vultr/govultr/pull/24)

## [v0.1.4](https://github.com/vultr/govultr/compare/v0.1.3..v0.1.4) (2019-07-14)
### Bug Fixes
* Fix panic on request failure [#20](https://github.com/vultr/govultr/pull/20)

## [v0.1.3](https://github.com/vultr/govultr/compare/v0.1.2..v0.1.3) (2019-06-13)
### Features
* added `GetVc2zList` to Plans to retrieve `high-frequency compute` plans [#13](https://github.com/vultr/govultr/pull/13)

### Breaking Changes
* Renamed all variables named `vpsID` to `instanceID` [#14](https://github.com/vultr/govultr/pull/14)
* Server
    * Renamed Server struct field `VpsID` to `InstanceID` [#14](https://github.com/vultr/govultr/pull/14)
* Plans
    * Renamed Plan struct field `VpsID` to `PlanID` [#14](https://github.com/vultr/govultr/pull/14)
    * Renamed BareMetalPlan struct field `BareMetalID` to `PlanID` [#14](https://github.com/vultr/govultr/pull/14)
    * Renamed VCPlan struct field `VpsID` to `PlanID` [#14](https://github.com/vultr/govultr/pull/14)
    * Renamed Plan struct field `VCPUCount` to `vCPUs` [#13](https://github.com/vultr/govultr/pull/13)
    * Renamed BareMetalPlan struct field `CPUCount` to `CPUs` [#13](https://github.com/vultr/govultr/pull/13)
    * Renamed VCPlan struct field `VCPUCount` to `vCPUs` [#13](https://github.com/vultr/govultr/pull/13)
    * Renamed VCPlan struct field `Cost` to `Price` [#13](https://github.com/vultr/govultr/pull/13)

## [v0.1.2](https://github.com/vultr/govultr/compare/v0.1.1..v0.1.2) (2019-05-29)
### Fixes
* Fixed Server Option `NotifyActivate` bug that ignored a `false` value
* Fixed Bare Metal Server Option `UserData` to be based64encoded 
### Breaking Changes
* Renamed all methods named `GetList` to `List`
* Renamed all methods named `Destroy` to `Delete`
* Server Service
    * Renamed `GetListByLabel` to `ListByLabel`
    * Renamed `GetListByMainIP` to `ListByMainIP`
    * Renamed `GetListByTag` to `ListByTag`
* Bare Metal Server Service
    * Renamed `GetListByLabel` to `ListByLabel`
    * Renamed `GetListByMainIP` to `ListByMainIP`
    * Renamed `GetListByTag` to `ListByTag`

## [v0.1.1](https://github.com/vultr/govultr/compare/v0.1.0..v0.1.1) (2019-05-20)
### Features
* add `SnapshotID` to ServerOptions as an option during server creation
* bumped default RateLimit from `.2` to `.6` seconds
### Breaking Changes
* Iso
  * Renamed all instances of `Iso` to `ISO`.  
* BlockStorage
  * Renamed `Cost` to `CostPerMonth`
  * Renamed `Size` to `SizeGB` 
* BareMetal & Server 
  * Change `SSHKeyID` to `SSHKeyIDs` which are now `[]string` instead of `string`
  * Renamed `OS` to `Os`    

## v0.1.0 (2019-05-10)
### Features
* Initial release
* Supports all available API endpoints that Vultr has to offer
