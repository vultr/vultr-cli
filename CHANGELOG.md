# Change Log
## [v2.15.1](https://github.com/vultr/vultr-cli/compare/v2.15.0...v2.15.1) (2023-03-09)
### Enhancements
* Update goreleaser to add latest docker image tag [PR 287](https://github.com/vultr/vultr-cli/pull/287)

### Documentation
* Block Storage: Make cli param examples consistently use = [PR 291](https://github.com/vultr/vultr-cli/pull/291)
* Instances: Make cli param examples consistently use = [PR 291](https://github.com/vultr/vultr-cli/pull/291)
* Regions: Add vcg plan options to docstrings [PR 299](https://github.com/vultr/vultr-cli/pull/299)
* Plans: Add vcg plan options to docstrings [PR 299](https://github.com/vultr/vultr-cli/pull/299)

### Dependencies
* Bump github.com/spf13/cobra from 1.5.0 to 1.6.0 by [PR 288](https://github.com/vultr/vultr-cli/pull/288)
* Bump github.com/spf13/cobra from 1.6.0 to 1.6.1 by [PR 289](https://github.com/vultr/vultr-cli/pull/289)
* Bump github.com/spf13/viper from 1.13.0 to 1.14.0 [PR 290](https://github.com/vultr/vultr-cli/pull/290)
* Bump golang.org/x/oauth2 from 0.0.0-20221014153046-6fdb5e3db783 to 0.5.0 [PR 296](https://github.com/vultr/vultr-cli/pull/296)
* Bump golang.org/x/oauth2 from 0.5.0 to 0.6.0 [PR 298](https://github.com/vultr/vultr-cli/pull/298)
* Bump github.com/spf13/viper from 1.14.0 to 1.15.0 [PR 293](https://github.com/vultr/vultr-cli/pull/293)
* Bump golang.org/x/net from 0.6.0 to 0.7.0 [PR 297](https://github.com/vultr/vultr-cli/pull/297)

### New Contributors
* @happytreees made their first contribution in [PR 287](https://github.com/vultr/vultr-cli/pull/287)

## [v2.15.0](https://github.com/vultr/vultr-cli/compare/v2.14.2...v2.15.0) (2022-10-04)
### Enhancements
* Add arm builds [PR 283](https://github.com/vultr/vultr-cli/pull/283)

### Dependencies
* Bump github.com/spf13/cobra from 1.4.0 to 1.5.0 [PR 274](https://github.com/vultr/vultr-cli/pull/274)
* Bump go from 1.17 to 1.19 [PR 284]( https://github.com/vultr/vultr-cli/pull/284)

### Documentation
* Remove extraneous dash from example command-line [PR 279](https://github.com/vultr/vultr-cli/pull/279)

### New Contributors
* @uplime made their first contribution in [PR 279](https://github.com/vultr/vultr-cli/pull/279)
* @mondragonfx made their first contribution in [PR 284](https://github.com/vultr/vultr-cli/pull/284)

## [v2.14.2](https://github.com/vultr/vultr-cli/compare/v2.14.1...v2.14.2) (2022-06-14)
### Enhancements
* Reserved IP: Add support for reserved IP label updates [PR 272](https://github.com/vultr/vultr-cli/pull/272)

### Dependencies
* Bump govultr version from 2.17.1 to 2.17.2 [PR 272](https://github.com/vultr/vultr-cli/pull/272)

## [v2.14.1](https://github.com/vultr/vultr-cli/compare/v2.14.0...v2.14.1) (2022-06-03)
### Enhancements
* Plans: Add GPU fields [PR 269](https://github.com/vultr/vultr-cli/pull/269)
* Instances: Update `tag` to string pointer [PR 268](https://github.com/vultr/vultr-cli/pull/268)
* Kuberneted: Update `tag` to string pointer [PR 268](https://github.com/vultr/vultr-cli/pull/268)

### Dependencies
* Bump github.com/spf13/viper from 1.11.0 to 1.12.0 [PR 266](https://github.com/vultr/vultr-cli/pull/266)
* Bump govultr version from 2.16.0 to 2.17.1 [PR 267](https://github.com/vultr/vultr-cli/pull/267)

## [v2.14.0](https://github.com/vultr/vultr-cli/compare/v2.13.0..v2.14.0) (2022-05-09)
### Enhancements
* Kubernetes : Add support for kubernetes version upgrades on individual clusters [PR 263](https://github.com/vultr/vultr-cli/pull/263)
* Kubernetes : Add support for node pool auto scaler options [PR 261](https://github.com/vultr/vultr-cli/pull/261)
* Firewall Rule : Update IP type option to match API verbiage for firewall rules [PR 262](https://github.com/vultr/vultr-cli/pull/262)
* Baremetal : Add support for multiple tags via the `tags` field [PR 259](https://github.com/vultr/vultr-cli/pull/259)
* Instances : Add support for multiple tags via the `tags` field [PR 259](https://github.com/vultr/vultr-cli/pull/259)

### Deprecations
* Firewall Rule : The `type` option on firewall rules has been replaced by `ip-type` [PR 262](https://github.com/vultr/vultr-cli/pull/262)
* Baremetal : the `tag` field has been replaced by `tags` which supports multiple tags [PR 259](https://github.com/vultr/vultr-cli/pull/259)
* Instances : the `tag` field has been replaced by `tags` which supports multiple tags [PR 259](https://github.com/vultr/vultr-cli/pull/259)

### Dependencies
* Bump github.com/vultr/govultr/v2 from 2.15.1 to 2.16.0 [PR 260](https://github.com/vultr/vultr-cli/pull/260)

### Documentation
* Update BSD install instructions [PR 258](https://github.com/vultr/vultr-cli/pull/258)
* Update README and improve verbiage for snapshots [PR 257](https://github.com/vultr/vultr-cli/pull/257)

## [v2.13.0](https://github.com/vultr/vultr-cli/compare/v2.12.2..v2.13.0) (2022-04-15)
### Enhancements
* VPC : new commands which will be replacing `network` (private networks) [PR 251](https://github.com/vultr/vultr-cli/pull/251)
* BlockStorage : adding support for new `block_type` field [PR 249](https://github.com/vultr/vultr-cli/pull/249/)
* LoadBalancer : Updating `vpc` functionality added [PR 251](https://github.com/vultr/vultr-cli/pull/251)

### Deprecations
* Network : These commands have been replaced by `vpc` [PR 251](https://github.com/vultr/vultr-cli/pull/251)
* Instance : The following fields have been deprecated on the `create` command `private-network` and `network`. Please use `vpc-enable` or `vpc-ids` [PR 251](https://github.com/vultr/vultr-cli/pull/251)
* LoadBalancer : The following fields have been deprecated on the `create` command `private-network`. Please use `vpc` instead [PR 251](https://github.com/vultr/vultr-cli/pull/251)

### Dependencies
* Bump github.com/vultr/govultr/v2 from 2.14.1 to 2.15.1 [PR 249](https://github.com/vultr/vultr-cli/pull/249)
* Bump github.com/spf13/viper from 1.10.1 to 1.11.0 [PR 252](https://github.com/vultr/vultr-cli/pull/252)

### Documentation
* Add Fedora installation instructions [PR 246](https://github.com/vultr/vultr-cli/pull/246)

## [v2.12.2](https://github.com/vultr/vultr-cli/compare/v2.12.1..v2.12.2) (2022-04-01)
### Enhancements
* Instances : fix csv flags `ssh-keys` and `network` [PR 244](https://github.com/vultr/vultr-cli/pull/244) @optik-aper
* Plans + Regions: Add new plan types in examples [PR 241](https://github.com/vultr/vultr-cli/pull/241/) @AFatalErrror
* Plans Metal : new command to retrieve just bare metal plans [PR 240](https://github.com/vultr/vultr-cli/pull/240) @optik-aper
* Readme : fix command example [PR 239](https://github.com/vultr/vultr-cli/pull/239) @travispaul

### Dependencies
* Bump github.com/vultr/govultr/v2 from 2.14.1 to 2.14.2 [PR 238](https://github.com/vultr/vultr-cli/pull/238)
* Bump github.com/spf13/cobra from 1.3.0 to 1.4.0 [PR 236](https://github.com/vultr/vultr-cli/pull/236)
* Bump builds from go 1.16 -> 1.17 [PR 243](https://github.com/vultr/vultr-cli/pull/243)

## [v2.12.1](https://github.com/vultr/vultr-cli/compare/v2.12.0..v2.12.1) (2022-02-07)
### Dependencies
* Bump github.com/vultr/govultr/v2 from 2.14.0 to 2.14.1 [PR 232](https://github.com/vultr/vultr-cli/pull/232)

### Enhancements
* Firewall Rule : Add ip type, source and subnet size to firewall rule printer [PR 234](https://github.com/vultr/vultr-cli/pull/234)

## [v2.12.0](https://github.com/vultr/vultr-cli/compare/v2.11.3..v2.12.0) (2022-01-21)
### Dependencies
* Bump github.com/vultr/govultr/v2 from 2.12.0 to 2.14.0 [PR 230](https://github.com/vultr/vultr-cli/pull/230)
* Bump github.com/spf13/viper from 1.10.0 to 1.10.1 [PR 224](https://github.com/vultr/vultr-cli/pull/224)
* Bump github.com/spf13/cobra from 1.2.1 to 1.3.0 [PR 223](https://github.com/vultr/vultr-cli/pull/223)

### Enhancements
* Script : Return b64 script when getting script by id [PR 229](https://github.com/vultr/vultr-cli/pull/229)

### Breaking Changes
* Script : get command will display data vertically now instead of horizontal [PR 229](https://github.com/vultr/vultr-cli/pull/229)

### Bug Fixes
* Firewalls : change source from int to string [PR 228](https://github.com/vultr/vultr-cli/pull/228)

## [v2.11.3](https://github.com/vultr/vultr-cli/compare/v2.11.2..v2.11.3) (2021-12-13)
### Dependencies
* Bump github.com/spf13/viper from 1.9.0 to 1.10.0 [PR 219](https://github.com/vultr/vultr-cli/pull/219)

### Enhancements
* Add OpenBSD install instructions [PR 218](https://github.com/vultr/vultr-cli/pull/218)

## [v2.11.2](https://github.com/vultr/vultr-cli/compare/v2.11.1..v2.11.2) (2021-12-01)
### Dependencies
* Update GoVultr from 2.11.1 to 2.12.0 [PR 215](https://github.com/vultr/vultr-cli/pull/215)

## [v2.11.1](https://github.com/vultr/vultr-cli/compare/v2.11.0..v2.11.1) (2021-11-29)
### Dependencies
* Bump github.com/vultr/govultr/v2 from 2.11.0 to 2.11.1 [PR 213](https://github.com/vultr/vultr-cli/pull/213)

### Bug Fixes
* Load Balancers : Allow SSL certificates to be passed in on Create and Update [PR 213](https://github.com/vultr/vultr-cli/pull/213)

## [v2.11.0](https://github.com/vultr/vultr-cli/compare/v2.10.0..v2.11.0) (2021-11-23)
### Enhancements
* DNS: Add support for getting a domains dns sec status [PR 211](https://github.com/vultr/vultr-cli/pull/211)
* Instance : Support changing hostname on reinstall [PR 209](https://github.com/vultr/vultr-cli/pull/209) [PR 210](https://github.com/vultr/vultr-cli/pull/210)

### Dependencies
* Update GoVultr from 2.10.0 to 2.11.0 [PR 209](https://github.com/vultr/vultr-cli/pull/209)

## [v2.10.0](https://github.com/vultr/vultr-cli/compare/v2.9.0..v2.10.0) (2021-11-04)
### Enhancements
* Billing: Add support for retrieving billing information [PR 203](https://github.com/vultr/vultr-cli/pull/203)

### Dependencies
* Update GoVultr from 2.9.2 to 2.10.0 [PR 203](https://github.com/vultr/vultr-cli/pull/203)

## [v2.9.0](https://github.com/vultr/vultr-cli/compare/v2.8.5..v2.9.0) (2021-10-27)
### Bug Fixes
* Allow `go get` and `go install` to work with `github.com/vultr/vultr-cli/v2` [PR 199](https://github.com/vultr/vultr-cli/pull/199)

## [v2.8.5](https://github.com/vultr/vultr-cli/compare/v2.8.4..v2.8.5) (2021-10-20)
### Dependencies
* Update GoVultr from 2.9.0 to 2.9.1 and update necessary fields [PR 196](https://github.com/vultr/vultr-cli/pull/196)
* Update GoVultr from 2.9.1 to 2.9.2 [PR 197](https://github.com/vultr/vultr-cli/pull/197)

### Enhancements
* Kubernetes: Add support for adding/modifying tags on Node Pools [PR 196](https://github.com/vultr/vultr-cli/pull/196)

## [v2.8.4](https://github.com/vultr/vultr-cli/compare/v2.8.3..v2.8.4) (2021-09-28)
### Dependencies
* Update GoVultr from 2.8.1 to 2.9.0 and update necessary fields [PR 192](https://github.com/vultr/vultr-cli/pull/192)

### Enhancements
* Snapshots: `COMPRESSED SIZE` has been added to printer output [PR 192](https://github.com/vultr/vultr-cli/pull/192)
* Kubernetes: `COUNT` has changed to `NODE QUANTITY` and `PLAN ID` has changed to `PLAN` for kubernetes printer output [PR 192](https://github.com/vultr/vultr-cli/pull/192)

## [v2.8.3](https://github.com/vultr/vultr-cli/compare/v2.8.2..v2.8.3) (2021-09-20)
### Dependencies
* Bump github.com/spf13/viper from 1.8.1 to 1.9.0 [PR 189](https://github.com/vultr/vultr-cli/pull/189)

### Bug Fixes
* Backups: Fix typo in backups alias [PR 188](https://github.com/vultr/vultr-cli/pull/188). Thanks @rmorey for your contribution

## [v2.8.2](https://github.com/vultr/vultr-cli/compare/v2.8.1..v2.8.2) (2021-09-07)
### Enhancements
* Instances: change default value for notify flag [PR 185](https://github.com/vultr/vultr-cli/pull/185)
* README: add example using boolean flag [PR 186](https://github.com/vultr/vultr-cli/pull/186)

## [v2.8.1](https://github.com/vultr/vultr-cli/compare/v2.8.0..v2.8.1) (2021-09-01)
### Dependencies
* GoVultr 2.8.0 -> 2.8.1 (added more kubernetes support)[PR 181](https://github.com/vultr/vultr-cli/pull/181)

### Enhancements
* Kubernetes: Add support for new Kubernetes calls [PR 181](https://github.com/vultr/vultr-cli/pull/181)
* Add User-Agent: [PR 182](https://github.com/vultr/vultr-cli/pull/182)

## [v2.8.0](https://github.com/vultr/vultr-cli/compare/v2.7.0..v2.8.0) (2021-08-23)
### Dependencies
* GoVultr 2.7.1 -> 2.8.0 (added kubernetes support)[PR 177](https://github.com/vultr/vultr-cli/pull/177)

### Enhancements
* Kubernetes: Add support for Kubernetes (VKE) [PR 178](https://github.com/vultr/vultr-cli/pull/178)
* README: update commands needed for building from source [PR 173](https://github.com/vultr/vultr-cli/pull/173)
* README: update examples  [PR 174](https://github.com/vultr/vultr-cli/pull/174)

## [v2.7.0](https://github.com/vultr/vultr-cli/compare/v2.6.0..v2.7.0) (2021-07-16)
### Dependencies
* GoVultr 2.6.0 -> 2.7.1 (added image_id support for instance and bare metal updates) [PR 169](https://github.com/vultr/vultr-cli/pull/169)

### Enhancements
* Instances: Add image_id support [PR 169](https://github.com/vultr/vultr-cli/pull/169)
* Bare-metal: Add image_id support [PR 169](https://github.com/vultr/vultr-cli/pull/169)
* Add documentation for autocompletions in README

## [v2.6.0](https://github.com/vultr/vultr-cli/compare/v2.5.3..v2.6.0) (2021-07-07)
### Dependencies
* Bump github.com/spf13/viper from 1.7.1 to 1.8.1 [PR 163](https://github.com/vultr/vultr-cli/pull/163)
* GoVultr v2.5.1 -> 2.6.0 (added support for persistent_pxe) [PR 164](https://github.com/vultr/vultr-cli/pull/164)

### Enhancements
* Bare-metal : Support `persistent_pxe` on create [PR 164](https://github.com/vultr/vultr-cli/pull/164)

## [v2.5.3](https://github.com/vultr/vultr-cli/compare/v2.5.2..v2.5.3) (2021-06-28)
### Dependencies
* Bump github.com/spf13/viper from 1.7.1 to 1.8.1 [PR 160](https://github.com/vultr/vultr-cli/pull/160)

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
