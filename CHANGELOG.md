# Changelog

## [0.4.1](https://github.com/trento-project/trento/tree/0.4.1) (2021-10-01)

[Full Changelog](https://github.com/trento-project/trento/compare/0.4.0...0.4.1)

### Added

- Add About page with subscription details [\#273](https://github.com/trento-project/trento/pull/273) (@arbulu89)
- Add a --rolling option to the install-agent script to use factory repos [\#270](https://github.com/trento-project/trento/pull/270) (@dottorblaster)
- Use only one GitHub runner instead of 1 per node using new install scripts [\#269](https://github.com/trento-project/trento/pull/269) (@rtorrero)
- Add warning and skipped states to the checks [\#266](https://github.com/trento-project/trento/pull/266) (@arbulu89)
- Discover SUSE subscription details [\#260](https://github.com/trento-project/trento/pull/260) (@arbulu89)
- Frontend tooling: introduce Prettier and ESLint [\#259](https://github.com/trento-project/trento/pull/259) (@dottorblaster)
- Add server installer [\#253](https://github.com/trento-project/trento/pull/253) (@fabriziosestito)
- Some improvements to the server installation on k3s through Helm [\#251](https://github.com/trento-project/trento/pull/251) (@dottorblaster)
- Build containers in CI/CD [\#250](https://github.com/trento-project/trento/pull/250) (@fabriziosestito)

### Fixed

- Trento RPM built in devel:sap:trento is not being injected the version number correctly [\#262](https://github.com/trento-project/trento/issues/262)
- Fix the liveness probe in the ARA chart due to sporadic sigterms [\#288](https://github.com/trento-project/trento/pull/288) (@rtorrero)
- Fix runner container image [\#287](https://github.com/trento-project/trento/pull/287) (@stefanotorresi)
- install-agent: fix typo in the script [\#280](https://github.com/trento-project/trento/pull/280) (@rtorrero)

### Removed

- Remove the discovery TTL consul health checks [\#208](https://github.com/trento-project/trento/issues/208)
- Remove docker-compose.yml [\#296](https://github.com/trento-project/trento/pull/296) (@stefanotorresi)
- remove unused gh action [\#294](https://github.com/trento-project/trento/pull/294) (@stefanotorresi)

### Other Changes

- Fix trento-server chart name [\#297](https://github.com/trento-project/trento/pull/297) (@fabriziosestito)
- Remove `$` from all the bash code examples [\#295](https://github.com/trento-project/trento/pull/295) (@stefanotorresi)
- Adjust scrips location [\#292](https://github.com/trento-project/trento/pull/292) (@nelsonkopliku)
- Fix test on subscription code, replacing Equal by ElementsMatch [\#289](https://github.com/trento-project/trento/pull/289) (@arbulu89)
- Minor improvements [\#286](https://github.com/trento-project/trento/pull/286) (@aleksei-burlakov)
- Use correct name for pre-release job [\#285](https://github.com/trento-project/trento/pull/285) (@rtorrero)
- halt deploy until images / packages are ready [\#284](https://github.com/trento-project/trento/pull/284) (@rtorrero)
- Rename sidebar entries to Pacemaker Clusters and SAP Systems [\#283](https://github.com/trento-project/trento/pull/283) (@arbulu89)
- fix wrong output path for the binaries [\#279](https://github.com/trento-project/trento/pull/279) (@rtorrero)
- Update subscription Load test to use ElementsMatch [\#276](https://github.com/trento-project/trento/pull/276) (@arbulu89)
- Adds an updated version of Trento Architecture Diagram [\#275](https://github.com/trento-project/trento/pull/275) (@nelsonkopliku)
- Run ansilble-lint in the CI process [\#272](https://github.com/trento-project/trento/pull/272) (@arbulu89)
- Set the TRENTO\_REPO variable properly inside install-agent.sh [\#271](https://github.com/trento-project/trento/pull/271) (@dottorblaster)
- Fix support for -e cluster\_selected\_checks= option [\#267](https://github.com/trento-project/trento/pull/267) (@brett060102)
- Fix: set the version in the Makefile explicitly [\#265](https://github.com/trento-project/trento/pull/265) (@aleksei-burlakov)
- Minor Helm chart updates [\#264](https://github.com/trento-project/trento/pull/264) (@stefanotorresi)
- add source label to dockerfile [\#263](https://github.com/trento-project/trento/pull/263) (@stefanotorresi)
- Add a default ansible configuration file usage [\#261](https://github.com/trento-project/trento/pull/261) (@arbulu89)
- Use new ghcr.io images [\#258](https://github.com/trento-project/trento/pull/258) (@rtorrero)
- Create the Trento ansible callback code [\#257](https://github.com/trento-project/trento/pull/257) (@arbulu89)
- Improve the metadata and check finding and move variables to defaults [\#255](https://github.com/trento-project/trento/pull/255) (@arbulu89)
- Fix: trento path is /usr/bin [\#252](https://github.com/trento-project/trento/pull/252) (@aleksei-burlakov)
- Fix obs workflow so that the submit job runs correctly on releases [\#248](https://github.com/trento-project/trento/pull/248) (@stefanotorresi)
- Remove discovery health checks; add trento-agent health check [\#240](https://github.com/trento-project/trento/pull/240) (@fabriziosestito)

## [0.4.0](https://github.com/trento-project/trento/releases/tag/0.4.0) 2021-09-15

### Added

- New Ansible-driven "Trento Runner" component, powering the main HA Checker feature (#150, #165, #187, #191, #204, #213)
- Add the checks catalog page (#159)
- Granular customization of the HA Checker rules (#181, #189, #211, #217)
- One-line installers and Continuously Delivered packages​ (#226)
- Helm Chart for the Trento Server deployment​​ (#206, #239, #242, #244)
- Add trento agent version visualization in the UI (#168, #198)
- Add discovery-period flag to the agent (#234)

### Changed

- Move the benchcommon based checks to ansible (#167, #172, #173, #174, #175, #180, #185, #190, #202, #223, #223, #225, #233)
- New SAP focused HANA Cluster and SAP Systems views​ (#169, #170, #171, #179, #192, #193, #222, #228)
- Tagging, filtering, and general navigation enhancements (#160, #221)
- Update the Dockerfile to use distroless containers (#241)
- Update the agent systemd unit file (#158)
- Rename the "Checks" section in the single host page (#238)
- Use full version of trento server in the web ui (#245)

### Fixed

- Fix and improve CI process steps (#182, #184, #212, #243)
- Fix Azure metadata discovery for different hypervisors (#229, #232)
- Fix assets upload in the CI process (#147, #156)
- Fix SID retrieval on the cluster context (#149)
- Fix the CIB Groups usage (it was not present before) (#151)

### Removed

- Remove Environment and Landscapes pages in favor of the new tagging system (#188)
- Remove automatic cluster name generation (#203)
- Remove HANA role badge (#210)
- Remove benchcommo checker for HA checks (#220)

## [0.3.0](https://github.com/trento-project/trento/releases/tag/0.3.0) 2021-07-14

### Added

- Check that the HANA and HANA SPS versions are compatible (#100)
- Check NTP time synchronization is configured (#105)
- Get system replication and landscape data for HANA (#106)
- Add self-hosted actions runner continuous deployment (#108)
- Add Check the hacluster user's password is not linux (#109)
- Add version command to CLI (#111)
- Add: Set log timestamp format (#113)
- Generate new consistent ID and human readable name for clusters (#114)
- Collect hdbnsutil -sr_state output (#115)
- Add resources members happy path tests (#119)
- Add: Check that HANA's autostart is disabled (#121)
- Add cloud metadata details - Azure (#126)
- Add pagination to the hosts view table (#132)
- Implement new styles for the host list view (#139)
- List clusters template revamp (#142)

### Changed

- Change: standard logging --> github.com/sirupsen/logrus (#107)
- Update hana-scale-up-perf-optimized-azure.yaml (#110)
- Single cluster template revamp (#120)
- Separate tags in the host list view (#123)
- Update generic layout styles to make the view wider (#138)

### Fixed

- Fix empty tables with a user-friendly message in case of no records (#102)
- Fix minor text polishings (#116)
- Lock mockery dependency version (#121)
- Fix broken test due to cluster-id introduction (#125)
- Add templates sanity check (#129)
- Fix TestSAPSystemHandler test (#130)
- Side menu: show the right title on hover (#131)
- Fix missing `</li>` in menu sidebar (#133)

## [0.2.0](https://github.com/trento-project/trento/releases/tag/0.2.0) 2021-06-16

### Added

- Add SAP Systems to default environment and landscape in absence of one (#70)
- Check that /etc/hosts contains all cluster nodes (#98)
- Check the UCAST is used by corosync with at least 2 com-n rings (#91, #96)
- Add project logo to the header (#90)
- Check that 2-nodes cluster must either have disk-based sbd or qdevice (#87)
- Landing page update with scope documentation (#82)
- Add this changelog ;) (#80)
- SBD configuration and service discovery (#72)

### Changed

- README updates
- Side bar and Home landing page improvements
- azure-rules check 1.3.5 was splitted into two checks
- Improve sidebar template (#84)
- Copy sapcontrol webservices from the exporter library instead of importing them (#81)
- Change how some checks are grouped together (#73, #74, #94)
- Reorganize the SAP System domain model structures (#75)

### Fixed

- Fix SAP system layout rendering
- Don't let the app crash on 404s (#97)
- Use the correct path for the global.ini config file of the SAPHanaSR check (#95)

## [0.1.0](https://github.com/trento-project/trento/releases/tag/0.1.0) 2021-05-26

### Added

- first release of Trento
- Automated discovery of SAP HANA HA clusters
- SAP Systems and Instances overview
- Grouping by Landscapes and Environments
- Configuration validation for Pacemaker, Corosync, SBD, SAPHanaSR and other + generic SUSE Linux Enterprise for SAP Application OS settings
- Specific configuration audits for SAP HANA Scale-Up Performance-Optimized
- scenarios deployed on MS Azure cloud.

---

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
