# Change Log
## [3.6.0](https://github.com/vultr/vultr-cli/compare/v3.5.0...v3.6.0) (2025-07-17)
### Enhancements
* Instance: Add user and password display [PR 541](https://github.com/vultr/vultr-cli/pull/541)
* Bare Metal: Add user and password display [PR 541](https://github.com/vultr/vultr-cli/pull/541)
* Database: Add support for additional Kafka features  [PR 543](https://github.com/vultr/vultr-cli/pull/543)

### Dependencies
* Update govultr from v3.20.0 to v3.21.0 [PR 542](https://github.com/vultr/vultr-cli/pull/542)
* Update govultr from v3.21.0 to v3.21.1 [PR 545](https://github.com/vultr/vultr-cli/pull/545)
* Bump github.com/go-viper/mapstructure/v2 from 2.2.1 to 2.3.0 [PR 544](https://github.com/vultr/vultr-cli/pull/544)

## [3.5.0](https://github.com/vultr/vultr-cli/compare/v3.4.0...v3.5.0) (2025-06-17)
### Enhancements
* Object Storage: add tier flag to create [PR 520](https://github.com/vultr/vultr-cli/pull/520)
* Database: adjust plans printer for new supported engines [PR 505](https://github.com/vultr/vultr-cli/pull/505)
* Database: add backup schedule [PR 534](https://github.com/vultr/vultr-cli/pull/534)
* Account: add bandwidth command [PR 535](https://github.com/vultr/vultr-cli/pull/535)
* Object Storage: update cluster display and add tiers [PR 536](https://github.com/vultr/vultr-cli/pull/536)
* Load Balancer: add auto-ssl [PR 521](https://github.com/vultr/vultr-cli/pull/521)
* Load Balancer: make SSL options read from files [PR 521](https://github.com/vultr/vultr-cli/pull/521)

### Bug Fixes
* Account: correctly set global command options [PR 517](https://github.com/vultr/vultr-cli/pull/517)
* Container Registry: add missing args check [PR 506](https://github.com/vultr/vultr-cli/pull/506)
* Kubernetes: fix nil error checks [PR 526](https://github.com/vultr/vultr-cli/pull/526)
* DNS: fix nil error checks [PR 526](https://github.com/vultr/vultr-cli/pull/526)
* Bare Metal: fix nil error checks [PR 526](https://github.com/vultr/vultr-cli/pull/526)
* Object Storage: fix printer output for tiers [PR 539](https://github.com/vultr/vultr-cli/pull/539)

### Documentation
* Instance: clarify userdata flag usage [PR 499](https://github.com/vultr/vultr-cli/pull/499)
* Kubernetes: fix whitespace on upgrades helptext [PR 522](https://github.com/vultr/vultr-cli/pull/522)
* Object Storage: add example help text for all commands [PR 538](https://github.com/vultr/vultr-cli/pull/538)

### Automation
* Remove unused lint config rules exclusion [PR 511](https://github.com/vultr/vultr-cli/pull/511)
* Remove coverage workflow [PR 518](https://github.com/vultr/vultr-cli/pull/518)
* Set goreleaser builds info to root [PR 527](https://github.com/vultr/vultr-cli/pull/527)
* Replace set-output with GITHUB_OUTPUT [PR 531](https://github.com/vultr/vultr-cli/pull/531)
* Update Golangci-lint to use v2 [PR 537](https://github.com/vultr/vultr-cli/pull/537)

### Deprecations
* VPC2: all vpc2 commands are deprecated [PR 525](https://github.com/vultr/vultr-cli/pull/525)
* Database: redis-named fields are deprecated [PR 495](https://github.com/vultr/vultr-cli/pull/495)

### Dependencies 
* Update govultr from v3.11.2 to v3.12.0 [PR 494](https://github.com/vultr/vultr-cli/pull/494)
* Update govultr from v3.12.0 to v3.14.1 [PR 503](https://github.com/vultr/vultr-cli/pull/503)
* Update govultr from v3.14.1 to v3.15.0 [PR 516](https://github.com/vultr/vultr-cli/pull/516)
* Update govultr from v3.15.0 to v3.17.0 [PR 524](https://github.com/vultr/vultr-cli/pull/524)
* Update govultr from v3.17.0 to v3.20.0 [PR 533](https://github.com/vultr/vultr-cli/pull/533)
* Update oauth from 0.23.0 to 0.24.0 [PR 496](https://github.com/vultr/vultr-cli/pull/496)
* Update go to v1.24 [PR 515](https://github.com/vultr/vultr-cli/pull/515)
* Bump github.com/spf13/cobra from 1.8.1 to 1.9.1 [PR 509](https://github.com/vultr/vultr-cli/pull/509)
* Bump golang.org/x/oauth2 from 0.25.0 to 0.27.0 [PR 510](https://github.com/vultr/vultr-cli/pull/510)
* Bump golang.org/x/oauth2 from 0.27.0 to 0.28.0 [PR 519](https://github.com/vultr/vultr-cli/pull/519)
* Bump golang.org/x/oauth2 from 0.24.0 to 0.25.0 [PR 500](https://github.com/vultr/vultr-cli/pull/500)
* Bump github.com/spf13/viper from 1.19.0 to 1.20.0 [PR 523](https://github.com/vultr/vultr-cli/pull/523)
* Bump golang.org/x/oauth2 from 0.28.0 to 0.29.0 [PR 529](https://github.com/vultr/vultr-cli/pull/529)
* Bump github.com/spf13/viper from 1.20.0 to 1.20.1 [PR 528](https://github.com/vultr/vultr-cli/pull/528)
* Bump golang.org/x/oauth2 from 0.29.0 to 0.30.0 [PR 532](https://github.com/vultr/vultr-cli/pull/532)

### New Contributors
* @mdspur made their first contribution in [PR 499](https://github.com/vultr/vultr-cli/pull/499)
* @kanha-gupta made their first contribution in [PR 506](https://github.com/vultr/vultr-cli/pull/506)

## [3.4.0](https://github.com/vultr/vultr-cli/compare/v3.3.1...v3.4.0) (2024-10-31)
### Enhancements
* Database: Add MySQL advanced config [PR 472](https://github.com/vultr/vultr-cli/pull/472)
* Load Balancer: Add nodes flag and printer display [PR 474](https://github.com/vultr/vultr-cli/pull/474)
* Database: Add support for Kafka [PR 487](https://github.com/vultr/vultr-cli/pull/487)

### Dependencies
* Bump golang.org/x/oauth2 from 0.21.0 to 0.22.0 [PR 470](https://github.com/vultr/vultr-cli/pull/470)
* Bump golang.org/x/oauth2 from 0.22.0 to 0.23.0 [PR 473](https://github.com/vultr/vultr-cli/pull/473)
* Update govultr from v3.9.0 to v3.9.1 [PR 471](https://github.com/vultr/vultr-cli/pull/471)
* Update govultr from v3.9.1 to v3.11.0 [PR 481](https://github.com/vultr/vultr-cli/pull/481)
* Update govultr from v3.11.0 to v3.11.1 [PR 485](https://github.com/vultr/vultr-cli/pull/485)
* Update govultr from v3.11.1 to v3.11.2 [PR 486](https://github.com/vultr/vultr-cli/pull/486)

### Documentation
* Update autocompletion instructions [PR 477](https://github.com/vultr/vultr-cli/pull/477)
* Add github CODEOWNERS file [PR 484](https://github.com/vultr/vultr-cli/pull/484)
* Fix instructions to install vultr-cli via Go [PR 475](https://github.com/vultr/vultr-cli/pull/475)
* Container Registry: Clarify public flag options [PR 479](https://github.com/vultr/vultr-cli/pull/479)

### Automation
* Remove deprecated exportloopref linter [PR 482](https://github.com/vultr/vultr-cli/pull/482)

### New Contributors
* @zikalino made their first contribution in [PR 475](https://github.com/vultr/vultr-cli/pull/475)
* @cjones-vultr made their first contribution in [PR 477](https://github.com/vultr/vultr-cli/pull/477)

## [3.3.1](https://github.com/vultr/vultr-cli/compare/v3.3.0...v3.3.1) (2024-07-29)
### Bug Fixes
* DNS: Make record create priority and ttl flags optional [PR 468](https://github.com/vultr/vultr-cli/pull/468)

### Dependencies
* Bump github.com/spf13/cobra from 1.8.0 to 1.8.1 [PR 460](https://github.com/vultr/vultr-cli/pull/460)

## [3.3.0](https://github.com/vultr/vultr-cli/compare/v3.2.0...v3.2.1) (2024-07-05)
### Enhancements
* CDN: Add support for CDN functions [PR 462](https://github.com/vultr/vultr-cli/pull/462)

### Bug Fixes
* Change mysql long run flag to int [PR 465](https://github.com/vultr/vultr-cli/pull/465)

### Documentation
* Add CDN command to README [PR 463](https://github.com/vultr/vultr-cli/pull/463)

### Dependencies
* Update govultr from v3.8.1 to v3.9.0 [PR 461](https://github.com/vultr/vultr-cli/pull/461)

## [3.2.0](https://github.com/vultr/vultr-cli/compare/v3.1.0...v3.2.0) (2024-06-11)
### Enhancements
* Inference: Add commands [PR 453](https://github.com/vultr/vultr-cli/pull/453)
* Bare Metal: Add mdisk_mode flag on create command [PR 449](https://github.com/vultr/vultr-cli/pull/449)

### Dependencies
* Bump golang.org/x/oauth2 from 0.20.0 to 0.21.0 [PR 451](https://github.com/vultr/vultr-cli/pull/451)
* Bump github.com/vultr/govultr/v3 from 3.6.4 to 3.7.0 [PR 448](https://github.com/vultr/vultr-cli/pull/448)
* Bump github.com/spf13/viper from 1.18.2 to 1.19.0 [PR 450](https://github.com/vultr/vultr-cli/pull/450)
* Update govultr from v2.7.0 to v2.8.1 [PR 452](https://github.com/vultr/vultr-cli/pull/452)

### Automation
* Linter updates and fixes [PR 457](https://github.com/vultr/vultr-cli/pull/457)
* Update goreleaser config and action [PR 456](https://github.com/vultr/vultr-cli/pull/456)
* Ignore linux binary [PR 455](https://github.com/vultr/vultr-cli/pull/455)

### New Contributors
* @fjoenichols made their first contribution in [PR 449](https://github.com/vultr/vultr-cli/pull/449)

## [v3.1.0](https://github.com/vultr/vultr-cli/compare/v3.0.3...v3.1.0) (2024-05-17)
### Bug Fixes
* Bare Metal: fix mistyped persistent_pxe baremetal flag [PR 432](https://github.com/vultr/vultr-cli/pull/432)
* All: fix error when a config file is not present [PR 434](https://github.com/vultr/vultr-cli/pull/434)
* Load Balancer: set correct flag type on LB proxy protocol [PR 436](https://github.com/vultr/vultr-cli/pull/436)

### Enhancements
* All: unify float precision on all printers [PR 437](https://github.com/vultr/vultr-cli/pull/437)
* Kubernetes: add node-pool list summarize flag [PR 442](https://github.com/vultr/vultr-cli/pull/442)

### Automation
* Update mattermost notify action [PR 441](https://github.com/vultr/vultr-cli/pull/441)
* Make mattermost notifications generic [PR 445](https://github.com/vultr/vultr-cli/pull/445)

### Dependencies 
* Bump golang.org/x/oauth2 from 0.18.0 to 0.19.0 [PR 440](https://github.com/vultr/vultr-cli/pull/440)
* Bump golang.org/x/oauth2 from 0.19.0 to 0.20.0 [PR 443](https://github.com/vultr/vultr-cli/pull/443)

### New Contributors
* @jasites made their first contribution in [PR 432](https://github.com/vultr/vultr-cli/pull/432)

## [v3.0.3](https://github.com/vultr/vultr-cli/compare/v3.0.2...v3.0.3) (2024-03-15) 
### Bug Fixes
* All: fix missing oauth token when a config file is used [PR 430](https://github.com/vultr/vultr-cli/pull/430)

### Dependencies
* Bump google.golang.org/protobuf from 1.31.0 to 1.33.0 [PR 428](https://github.com/vultr/vultr-cli/pull/428)

## [v3.0.2](https://github.com/vultr/vultr-cli/compare/v3.0.1...v3.0.2) (2024-03-12)
### Enhancements
* Kubernetes: add node labels flag to node pool create and update [PR 422](https://github.com/vultr/vultr-cli/pull/422)

### Bug Fixes
* VPC2: fix incorrect govultr delete method [PR 426](https://github.com/vultr/vultr-cli/pull/426)
* Plans: fix inadvertent short-circuiting bare metal printer [PR 416](https://github.com/vultr/vultr-cli/pull/416)
* All: make cobra command run error return consistent [PR 419](https://github.com/vultr/vultr-cli/pull/419)

### Documentation
* Remove extraneous readme details [PR 425](https://github.com/vultr/vultr-cli/pull/425)
* Tidy up help documentation [PR 418](https://github.com/vultr/vultr-cli/pull/418)

### Dependencies
* Update govultr from v3.6.3 to v3.6.4 [PR 421](https://github.com/vultr/vultr-cli/pull/421)
* Update govultr from v3.6.2 to v3.6.3 [PR 417](https://github.com/vultr/vultr-cli/pull/417)
* Bump golang.org/x/oauth2 from 0.17.0 to 0.18.0 [PR 420](https://github.com/vultr/vultr-cli/pull/420)

## [v3.0.1](https://github.com/vultr/vultr-cli/compare/v3.0.0...v3.0.1) (2024-02-27)
### Enhancements
* Kubernetes: add `enable-firewall` flag on create [PR 413](https://github.com/vultr/vultr-cli/pull/413)
* All: consolidate pagination metadata nil pointer checks [PR 410](https://github.com/vultr/vultr-cli/pull/410)

### Bug Fixes
* Kubernetes: remove shorthand flag conflict with output [PR 406](https://github.com/vultr/vultr-cli/pull/406)
* Regions: make printer consistent with returned data [PR 412](https://github.com/vultr/vultr-cli/pull/412)
* Plans: make printer consistent with returned data [PR 412](https://github.com/vultr/vultr-cli/pull/412)

### Dependencies
* Bump github.com/vultr/govultr/v3 from 3.6.1 to 3.6.2 [PR 407](https://github.com/vultr/vultr-cli/pull/407)

### Documentation
* Update README to use v3 in `go get` command [PR 405](https://github.com/vultr/vultr-cli/pull/405)
* Firewall: add correct usage for firewall rules [PR 414](https://github.com/vultr/vultr-cli/pull/414)

### New Contributors
* @PaulSonOfLars made their first contribution in [PR 410](https://github.com/vultr/vultr-cli/pull/410)

## [v3.0.0](https://github.com/vultr/vultr-cli/compare/v2.22.0...v3.0.0) (2024-02-15)
### Enhancements
* [Complete refactor](https://github.com/vultr/vultr-cli/pull/402) of the CLI commands and project packages. All commands have been restructured and standardized with these goals in mind:
  * Move commands into separate packages
  * All output through a common interface.  Now supporting JSON and YAML on all commands using the `--output` flag.
  * Auth checks only happen when appropriate to the API endpoint
  * Some generically useful stuff like printers for IPs or User Data have been moved out to their own packages
  * Base functionality should all be the same at this point, but there is room for improvement:
    * Common error formatting for API messages
    * More configuration options
    * Better testability

### Dependencies
* Bump golang.org/x/oauth2 from 0.16.0 to 0.17.0 [PR 397](https://github.com/vultr/vultr-cli/pull/397)

## [v2.22.0](https://github.com/vultr/vultr-cli/compare/v2.21.0...v2.22.0) (2024-02-01)
### Enhancements
* Database: add user access control for redis [PR 383](https://github.com/vultr/vultr-cli/pull/383)
* Marketplace: add support for app variables [PR 389](https://github.com/vultr/vultr-cli/pull/389)

### Bug Fixes
* Conatiner Registry: fix read-write flag on docker credentials [PR 395](https://github.com/vultr/vultr-cli/pull/395)

### Dependencies
* Bump github.com/vultr/govultr/v3 from 3.4.1 to 3.5.0 [PR 382](https://github.com/vultr/vultr-cli/pull/382)
* Update govultr from v3.5.0 to v3.6.0 [PR 388](https://github.com/vultr/vultr-cli/pull/388)
* Bump github.com/vultr/govultr/v3 from 3.6.0 to 3.6.1 [PR 393](https://github.com/vultr/vultr-cli/pull/393)
* Bump golang.org/x/oauth2 from 0.15.0 to 0.16.0 [PR 392](https://github.com/vultr/vultr-cli/pull/392)
* Bump github.com/spf13/viper from 1.17.0 to 1.18.2 [PR 390](https://github.com/vultr/vultr-cli/pull/390)

### New Contributors
* @biondizzle made their first contribution in [PR 395](https://github.com/vultr/vultr-cli/pull/395)

## [v2.21.0](https://github.com/vultr/vultr-cli/compare/v2.20.0...v2.21.0) (2023-11-29)
### Enhancements
* Database: Add usage commands [PR 378](https://github.com/vultr/vultr-cli/pull/378)
* Container Registry: Implemented [PR 380](https://github.com/vultr/vultr-cli/pull/380)
* Bare Metal: Update tags display to use delimiters [PR 372](https://github.com/vultr/vultr-cli/pull/372)
* Instance: Update tags display to use delimiters [PR 372](https://github.com/vultr/vultr-cli/pull/372)
* Database: Add read replica promotion [PR 375](https://github.com/vultr/vultr-cli/pull/375)
* Kubernetes: Add kubeconfig filepath export [PR 361](https://github.com/vultr/vultr-cli/pull/361)

### Dependencies
* Update govultr to v3.4.1 [PR 376](https://github.com/vultr/vultr-cli/pull/376)
* Bump golang.org/x/oauth2 from 0.14.0 to 0.15.0 [PR 379](https://github.com/vultr/vultr-cli/pull/379)
* Bump github.com/spf13/cobra from 1.7.0 to 1.8.0 [PR 371](https://github.com/vultr/vultr-cli/pull/371)
* Bump govultr to v3.4.0 [PR 374](https://github.com/vultr/vultr-cli/pull/374)
* Bump golang.org/x/oauth2 from 0.13.0 to 0.14.0 [PR 373](https://github.com/vultr/vultr-cli/pull/373)

## [v2.20.0](https://github.com/vultr/vultr-cli/compare/v2.19.0...v2.20.0) (2023-11-01)
### Enhancements
* Managed Database public/private hostnames, cleanup summarize view [PR 363](https://github.com/vultr/vultr-cli/pull/363)
* Allow some commands to be run without authenticating against the API [PR 364](https://github.com/vultr/vultr-cli/pull/364)
* Add support for the VKE HA control plane option [PR 368](https://github.com/vultr/vultr-cli/pull/368)
* Add Support for DBaaS FerretDB Subscriptions [PR 369](https://github.com/vultr/vultr-cli/pull/369)

### Bug Fixes
* Adjust DBaaS VPC pointer to detect changes [PR 366](https://github.com/vultr/vultr-cli/pull/366)

### Dependencies
* Bump golang.org/x/net from 0.15.0 to 0.17.0 [PR 358](https://github.com/vultr/vultr-cli/pull/358)
* Update govultr to v3.3.2 [PR 362](https://github.com/vultr/vultr-cli/pull/362)
* Update govultr to v3.3.3 [PR 365](https://github.com/vultr/vultr-cli/pull/365)
* Update govultr to v3.3.4 [PR 367](https://github.com/vultr/vultr-cli/pull/367)
* Bump golang.org/x/oauth2 from 0.12.0 to 0.13.0 [PR 356](https://github.com/vultr/vultr-cli/pull/356)
* Bump github.com/spf13/viper from 1.16.0 to 1.17.0 [PR 357](https://github.com/vultr/vultr-cli/pull/357)

## [v2.19.0](https://github.com/vultr/vultr-cli/compare/v2.18.2...v2.19.0) (2023-10-18)
### Enhancements
* Kubernetes: Add summarize list options [PR 348](https://github.com/vultr/vultr-cli/pull/348)
* Database: Add summarize list options [PR 348](https://github.com/vultr/vultr-cli/pull/348)
* Load Balancer: Add summarize list options [PR 348](https://github.com/vultr/vultr-cli/pull/348)
* Rework the printer output code [PR 355](https://github.com/vultr/vultr-cli/pull/355)

### Documentation
* VPC2: Correct create command example [PR 350](https://github.com/vultr/vultr-cli/pull/350)

### Bug Fixes
* Remove the useless cobra init help toggle flag [PR 349](https://github.com/vultr/vultr-cli/pull/349)

### Dependencies
* Bump golang.org/x/oauth2 from 0.11.0 to 0.12.0 [PR 351](https://github.com/vultr/vultr-cli/pull/351)
* Update to go v1.21 [PR 347](https://github.com/vultr/vultr-cli/pull/347)

### Automation
* Add project name back to the archive file names in goreleaser [PR 346](https://github.com/vultr/vultr-cli/pull/346)
* Add golangci-lint and fix linter errors [PR 353](https://github.com/vultr/vultr-cli/pull/353)
* Remove unnecessary Go dependency from .goreleaser [PR 97](https://github.com/vultr/vultr-cli/pull/97)

### New Contributors
* @0az made their first contribution in [PR 97](https://github.com/vultr/vultr-cli/pull/97)
* @resmo made their first contribution in [PR 350](https://github.com/vultr/vultr-cli/pull/350)

## [v2.18.2](https://github.com/vultr/vultr-cli/compare/v2.18.0...v2.18.2) (2023-08-24)
### Automation
* Update how archive names are generated by goreleaser [PR 342](https://github.com/vultr/vultr-cli/pull/342)
* Remove deprecated brews tap command in goreleaser [PR_344](https://github.com/vultr/vultr-cli/pull/344)

## [v2.18.0](https://github.com/vultr/vultr-cli/compare/v2.17.0...v2.18.0) (2023-08-23)
### Enhancements
* Database: Add VPC support for DBaaS instances [PR 331](https://github.com/vultr/vultr-cli/pull/331)
* Bare Metal: Add support for VPC 2.0 [PR 335](https://github.com/vultr/vultr-cli/pull/335)
* Instance: Add support for VPC 2.0 [PR 335](https://github.com/vultr/vultr-cli/pull/335)
* Application: Add more aliases for the apps command [PR 336](https://github.com/vultr/vultr-cli/pull/336)
* VPC2: Add Nodes Endpoints [PR 339](https://github.com/vultr/vultr-cli/pull/339)
* Database: Managed Database Nesting Refactor [PR 340](https://github.com/vultr/vultr-cli/pull/340)

### Bug Fixes
* Instance: Fix reserved IPv4 flag docs [PR 337](https://github.com/vultr/vultr-cli/pull/337)

### Dependencies
* Update govultr to v3.3.0 [PR 334](https://github.com/vultr/vultr-cli/pull/334)
* Update govultr to v3.3.1 [PR 338](https://github.com/vultr/vultr-cli/pull/338)
* Update govultr to v3.1.0 [PR 329](https://github.com/vultr/vultr-cli/pull/329)
* Bump github.com/vultr/govultr/v3 from 3.1.0 to 3.2.0 [PR 330](https://github.com/vultr/vultr-cli/pull/330)
* Bump golang.org/x/oauth2 from 0.9.0 to 0.10.0 [PR 328](https://github.com/vultr/vultr-cli/pull/328)
* Bump golang.org/x/oauth2 from 0.10.0 to 0.11.0 [PR 333](https://github.com/vultr/vultr-cli/pull/333)

### New Contributors
* @nhooyr made their first contribution in [PR 337](https://github.com/vultr/vultr-cli/pull/337)

## [v2.17.0](https://github.com/vultr/vultr-cli/compare/v2.16.2...v2.17.0) (2023-06-14)
### Enhancements
* Instances: Add support for attaching and detaching VPC networks [PR 318](https://github.com/vultr/vultr-cli/pull/318)

### Bug Fixes
* Database: Fix database update errors and remove db engine/version [PR 314](https://github.com/vultr/vultr-cli/pull/314)

### Documentation
* README: Use a more succinct Homebrew command to tap-and-install [PR 315](https://github.com/vultr/vultr-cli/pull/315)
* README: Fix spelling [PR 324](https://github.com/vultr/vultr-cli/pull/324)
* README: Add docker install/usage instructions [PR 322](https://github.com/vultr/vultr-cli/pull/322)
* README: Mention default config yaml location [PR 325](https://github.com/vultr/vultr-cli/pull/325)

### Dependencies
* Bump github.com/vultr/govultr/v3 from 3.0.2 to 3.0.3 [PR 320](https://github.com/vultr/vultr-cli/pull/320)
* Bump github.com/spf13/cobra from 1.6.1 to 1.7.0 [PR 310](https://github.com/vultr/vultr-cli/pull/310)
* Bump github.com/spf13/viper from 1.15.0 to 1.16.0 [PR 319](https://github.com/vultr/vultr-cli/pull/319)
* Bump golang.org/x/oauth2 from 0.6.0 to 0.7.0 [PR 312](https://github.com/vultr/vultr-cli/pull/312)
* Bump golang.org/x/oauth2 from 0.7.0 to 0.8.0 [PR 316](https://github.com/vultr/vultr-cli/pull/316)
* Bump golang.org/x/oauth2 from 0.8.0 to 0.9.0 [PR 323](https://github.com/vultr/vultr-cli/pull/323)
* Update Github workflows to go v1.20 [PR 311](https://github.com/vultr/vultr-cli/pull/311)

### New Contributors
* @ELLIOTTCABLE made their first contribution in [PR 315](https://github.com/vultr/vultr-cli/pull/315)

## [v2.16.2](https://github.com/vultr/vultr-cli/compare/v2.15.1...v2.16.2) (2023-03-31)
### Enhancements
* Database: Add DBaaS Support [PR 302](https://github.com/vultr/vultr-cli/pull/302)

### Dependencies
* Update go to 1.20 [PR 303](https://github.com/vultr/vultr-cli/pull/303)
* Update govultr to v3.0.1 [PR 301](https://github.com/vultr/vultr-cli/pull/301)
* Update govultr to v3.0.2 [PR 304](https://github.com/vultr/vultr-cli/pull/304)
* Fix goreleaser configurations [PR 306](https://github.com/vultr/vultr-cli/pull/306)
* Fix github automatic release configurations [PR 308](https://github.com/vultr/vultr-cli/pull/308)

### New Contributors
* @christhemorse made their first contribution in [PR 302](https://github.com/vultr/vultr-cli/pull/302)

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
