# Changelog

## [0.7.1](https://github.com/trento-project/trento/tree/0.7.1) (2021-12-21)

[Full Changelog](https://github.com/trento-project/trento/compare/0.7.0...0.7.1)

### Added

- Missing info about HANA instance when trento agent is stopped in primary node [\#536](https://github.com/trento-project/trento/issues/536)
- Add the cluster modal user interaction in the new REST api [\#219](https://github.com/trento-project/trento/issues/219)
- Add dependabot.yml config for GH [\#612](https://github.com/trento-project/trento/pull/612) (@rtorrero)

### Fixed

- Pacemaker Cluster - created tags don't showed up if other filter selected [\#602](https://github.com/trento-project/trento/issues/602)
- Make the table cell alignment consistent in the cluster detail page [\#535](https://github.com/trento-project/trento/issues/535)
- Information missing or wrong information displayed in console after full shutdown [\#515](https://github.com/trento-project/trento/issues/515)
- Wrong HANA cluster info after failover [\#513](https://github.com/trento-project/trento/issues/513)
- HANA primary info missing after failover [\#512](https://github.com/trento-project/trento/issues/512)
- Status of SAP instance doesn't get updated after trento-agent got stopped  [\#491](https://github.com/trento-project/trento/issues/491)
- The fencing timeout check fails if the time unit is set [\#447](https://github.com/trento-project/trento/issues/447)
- Point the telemetry service to telemetry.trento.suse.com [\#640](https://github.com/trento-project/trento/pull/640) (@dottorblaster)
- Fix test 373DB8 to pass if timeout unit is set [\#639](https://github.com/trento-project/trento/pull/639) (@arbulu89)
- Align properly the cluster sites tables columns [\#628](https://github.com/trento-project/trento/pull/628) (@arbulu89)
- Fix cluster list filters [\#627](https://github.com/trento-project/trento/pull/627) (@fabriziosestito)
- Paginate correctly sap systems [\#626](https://github.com/trento-project/trento/pull/626) (@arbulu89)
- Update gin gonic to the latest version due to major breaking bug in the router mechanism [\#610](https://github.com/trento-project/trento/pull/610) (@fabriziosestito)
- Fix cluster details resource issues [\#604](https://github.com/trento-project/trento/pull/604) (@arbulu89)

### Closed Issues

- SAP systems view: PAS instance \(01\) not listed under NWP [\#609](https://github.com/trento-project/trento/issues/609)
- Cluster setting: no field to enter connection user [\#608](https://github.com/trento-project/trento/issues/608)
- Bad Gateway or always get an error message about fetching checks date for hana\_cluster [\#606](https://github.com/trento-project/trento/issues/606)
- Pacemaker Clusters details - Health - 'Show check results' is misleading [\#605](https://github.com/trento-project/trento/issues/605)
- Tilde is not correctly resolved on installation script [\#492](https://github.com/trento-project/trento/issues/492)
- Trento Agent version name should be consistent [\#490](https://github.com/trento-project/trento/issues/490)
- Port the SAP systems list to the new architecture [\#339](https://github.com/trento-project/trento/issues/339)
- Cluster Detail page displays an error even if the trento-runner is running [\#330](https://github.com/trento-project/trento/issues/330)

### Other Changes

- Bump eslint from 7.32.0 to 8.5.0 in /web/frontend [\#632](https://github.com/trento-project/trento/pull/632) (@dependabot[bot])
- Bump prettier from 2.4.1 to 2.5.1 in /web/frontend [\#630](https://github.com/trento-project/trento/pull/630) (@dependabot[bot])
- Bump sass from 1.45.0 to 1.45.1 in /web/frontend [\#629](https://github.com/trento-project/trento/pull/629) (@dependabot[bot])
- Bump github.com/lib/pq from 1.10.2 to 1.10.4 [\#625](https://github.com/trento-project/trento/pull/625) (@dependabot[bot])
- Bump github.com/tdewolff/minify/v2 from 2.9.16 to 2.9.24 [\#624](https://github.com/trento-project/trento/pull/624) (@dependabot[bot])
- Bump github.com/swaggo/swag from 1.7.4 to 1.7.6 [\#622](https://github.com/trento-project/trento/pull/622) (@dependabot[bot])
- Bump github.com/spf13/afero from 1.1.2 to 1.6.0 [\#621](https://github.com/trento-project/trento/pull/621) (@dependabot[bot])
- Bump @babel/preset-react from 7.14.5 to 7.16.5 in /web/frontend [\#620](https://github.com/trento-project/trento/pull/620) (@dependabot[bot])
- Bump axios from 0.21.4 to 0.24.0 in /web/frontend [\#618](https://github.com/trento-project/trento/pull/618) (@dependabot[bot])
- Bump @yaireo/tagify from 4.7.2 to 4.9.4 in /web/frontend [\#617](https://github.com/trento-project/trento/pull/617) (@dependabot[bot])
- Bump docker/login-action from 1.10.0 to 1.11.0 [\#616](https://github.com/trento-project/trento/pull/616) (@dependabot[bot])
- Bump sass from 1.32.8 to 1.45.0 in /web/frontend [\#615](https://github.com/trento-project/trento/pull/615) (@dependabot[bot])
- Bump docker/metadata-action from 3.3.0 to 3.6.2 [\#614](https://github.com/trento-project/trento/pull/614) (@dependabot[bot])
- Bump actions/setup-python from 1 to 2.3.1 [\#613](https://github.com/trento-project/trento/pull/613) (@dependabot[bot])

## [0.7.0](https://github.com/trento-project/trento/tree/0.7.0) (2021-12-17)

[Full Changelog](https://github.com/trento-project/trento/compare/0.6.0...0.7.0)

### Added

- Premium stuff. [\#582](https://github.com/trento-project/trento/issues/582)
- Provide direct navigation from the check results modal to the catalog [\#532](https://github.com/trento-project/trento/issues/532)
- Use the host identifier as host details page reference instead of the hostname [\#521](https://github.com/trento-project/trento/issues/521)
- Improvements in the runner logging [\#469](https://github.com/trento-project/trento/issues/469)
- Introduce Helm chart tests in the CI pipeline [\#405](https://github.com/trento-project/trento/issues/405)
- Wrap all the GetSelectedChecksById calls in the runner, to a unique API call [\#383](https://github.com/trento-project/trento/issues/383)
- Add PremiumDetection capabilities [\#586](https://github.com/trento-project/trento/pull/586) (@nelsonkopliku)
- Reduce discovery times to 10 seconds and do not store duplicated events [\#581](https://github.com/trento-project/trento/pull/581) (@fabriziosestito)
- Refactor cluster settings fe [\#561](https://github.com/trento-project/trento/pull/561) (@fabriziosestito)
- Add live logging ansible [\#530](https://github.com/trento-project/trento/pull/530) (@arbulu89)
- Add flavor to version package [\#505](https://github.com/trento-project/trento/pull/505) (@rtorrero)
- Migrate about page and sles subscriptions data to the database [\#461](https://github.com/trento-project/trento/pull/461) (@arbulu89)
- Split the runner environment vars in directories [\#440](https://github.com/trento-project/trento/pull/440) (@rtorrero)

### Fixed

- Remove column 'System' from HANA Databases - Hosts [\#549](https://github.com/trento-project/trento/issues/549)
- SAP Systems - Filters -SID doesn't work or works only with big delay [\#545](https://github.com/trento-project/trento/issues/545)
- Filters tags added at page of Hosts shows up on page Peacemake Clusters [\#544](https://github.com/trento-project/trento/issues/544)
- List in "Settings \> Checks catalog" wrongly rendered [\#543](https://github.com/trento-project/trento/issues/543)
- Instance numbers with one digit in SAP system detail view [\#508](https://github.com/trento-project/trento/issues/508)
- Cluster detail view: checks selected but not showing in Health section [\#507](https://github.com/trento-project/trento/issues/507)
- Cluster view: Health section not showing clusters in any status [\#506](https://github.com/trento-project/trento/issues/506)
- Unknown host listed in Hosts view [\#503](https://github.com/trento-project/trento/issues/503)
- Bad Gateway error when navigating through the different views in the console [\#502](https://github.com/trento-project/trento/issues/502)
- SAP System link in Host Detail view takes you to a Not Found page \(The requested URL doesn't exist\) [\#501](https://github.com/trento-project/trento/issues/501)
- "SAP System details" instead of "HANA Database details" shows up if open a SID  [\#495](https://github.com/trento-project/trento/issues/495)
- cannot detect hosts status correctly or returns 500 error code [\#482](https://github.com/trento-project/trento/issues/482)
- "Warning" value stays 0 even there is a problem with duplicated cluster name [\#479](https://github.com/trento-project/trento/issues/479)
- Trento runner reports all checks in green if the `ssh` command is not available [\#277](https://github.com/trento-project/trento/issues/277)
- Cluster nodes disappearing from the cluster list when HA discovery loop fails [\#205](https://github.com/trento-project/trento/issues/205)
- HANA secondary sync state should display a message when replication is not working [\#154](https://github.com/trento-project/trento/issues/154)
- Fix null return in check settings [\#603](https://github.com/trento-project/trento/pull/603) (@fabriziosestito)
- fix bogus docker build makefile error [\#592](https://github.com/trento-project/trento/pull/592) (@stefanotorresi)
- Fix SBD checks 0B6DB2 and 49591F remediation rendering [\#587](https://github.com/trento-project/trento/pull/587) (@arbulu89)
- Fix Cluster hana detail host link [\#555](https://github.com/trento-project/trento/pull/555) (@fabriziosestito)
- Fix Pacemaker Site Details displaying name-related host hrefs [\#554](https://github.com/trento-project/trento/pull/554) (@dottorblaster)
- Fix sap systems template [\#547](https://github.com/trento-project/trento/pull/547) (@fabriziosestito)
- Set a timeout in ansible tasks [\#534](https://github.com/trento-project/trento/pull/534) (@arbulu89)
- Show SAP instance number using 2 digits [\#527](https://github.com/trento-project/trento/pull/527) (@arbulu89)
- Create properly the SAP system url in the host page [\#526](https://github.com/trento-project/trento/pull/526) (@arbulu89)
- Show HANA details in its details page [\#525](https://github.com/trento-project/trento/pull/525) (@arbulu89)

### Removed

- It's the final consul cleanup [\#569](https://github.com/trento-project/trento/pull/569) (@nelsonkopliku)
- Removed consul config dir option [\#567](https://github.com/trento-project/trento/pull/567) (@nelsonkopliku)
- Cleanup agent from consul [\#559](https://github.com/trento-project/trento/pull/559) (@nelsonkopliku)
- Remove cluster generic [\#551](https://github.com/trento-project/trento/pull/551) (@fabriziosestito)
- Remove premium checks and their variables [\#500](https://github.com/trento-project/trento/pull/500) (@arbulu89)

### Closed Issues

- Pacemaker Clusters -  Health status always has value '0' for Passing, Warning, Critical all the time [\#601](https://github.com/trento-project/trento/issues/601)
- Bad Gateway got displayed for trento server - Hosts [\#598](https://github.com/trento-project/trento/issues/598)
- All hosts got red icon '!'  regardless trento agent is running  [\#596](https://github.com/trento-project/trento/issues/596)
- The icon of magnifying glass is misleading [\#540](https://github.com/trento-project/trento/issues/540)
- Refactor cluster entity/model to have just one SID [\#493](https://github.com/trento-project/trento/issues/493)
- Add a make target to build API documentation [\#321](https://github.com/trento-project/trento/issues/321)

### Other Changes

- Revert collector deduplication [\#599](https://github.com/trento-project/trento/pull/599) (@fabriziosestito)
- Move About under Settings [\#595](https://github.com/trento-project/trento/pull/595) (@stefanotorresi)
- remove repetition in the checks description [\#594](https://github.com/trento-project/trento/pull/594) (@stefanotorresi)
- follow official product guidelines [\#593](https://github.com/trento-project/trento/pull/593) (@stefanotorresi)
- "About premium" improvements [\#591](https://github.com/trento-project/trento/pull/591) (@rtorrero)
- Add missing hosts preload in attached database retrieval [\#590](https://github.com/trento-project/trento/pull/590) (@fabriziosestito)
- Remove the extra space in the BuiltTag comment [\#589](https://github.com/trento-project/trento/pull/589) (@arbulu89)
- Fix attached database instances hydration [\#588](https://github.com/trento-project/trento/pull/588) (@fabriziosestito)
- Remove hardcoded constant for the flavor to a Makefile variable [\#585](https://github.com/trento-project/trento/pull/585) (@rtorrero)
- Add suse registry tag labels in the helm chart [\#584](https://github.com/trento-project/trento/pull/584) (@arbulu89)
- Hide system column in SAP System detail hosts table [\#583](https://github.com/trento-project/trento/pull/583) (@fabriziosestito)
- Disable test parallelism [\#580](https://github.com/trento-project/trento/pull/580) (@fabriziosestito)
- Reduce the agent discovery interval default value to 30 seconds [\#578](https://github.com/trento-project/trento/pull/578) (@fabriziosestito)
- Further consul cleanup [\#577](https://github.com/trento-project/trento/pull/577) (@nelsonkopliku)
- Updated architecture diagram to the consul-free version [\#575](https://github.com/trento-project/trento/pull/575) (@nelsonkopliku)
- Bump postgresql version [\#574](https://github.com/trento-project/trento/pull/574) (@fabriziosestito)
- Disable host\_key\_checking when running ansible playbook [\#572](https://github.com/trento-project/trento/pull/572) (@nelsonkopliku)
- Fixed clusters settings endpoint leftover [\#571](https://github.com/trento-project/trento/pull/571) (@nelsonkopliku)
- Add error information about inventory content creation [\#570](https://github.com/trento-project/trento/pull/570) (@nelsonkopliku)
- Stop trento agent before installing the new one [\#568](https://github.com/trento-project/trento/pull/568) (@fabriziosestito)
- Leftover cleanup [\#566](https://github.com/trento-project/trento/pull/566) (@nelsonkopliku)
- Fix cluster type field in hosts projection [\#565](https://github.com/trento-project/trento/pull/565) (@fabriziosestito)
- Fixes ssh-address required option [\#564](https://github.com/trento-project/trento/pull/564) (@nelsonkopliku)
- Remove consul references from the README [\#562](https://github.com/trento-project/trento/pull/562) (@dottorblaster)
- Clean up helm chart form consul references [\#560](https://github.com/trento-project/trento/pull/560) (@nelsonkopliku)
- Make runner use clusters settings API instead of Consul [\#558](https://github.com/trento-project/trento/pull/558) (@nelsonkopliku)
- Change the license icon inside the sidebar to assignment one [\#557](https://github.com/trento-project/trento/pull/557) (@dottorblaster)
- Use the %{name} macro in Provides [\#556](https://github.com/trento-project/trento/pull/556) (@arbulu89)
- Hosts UI revamp [\#552](https://github.com/trento-project/trento/pull/552) (@dottorblaster)
- Refactor sap system detail view [\#548](https://github.com/trento-project/trento/pull/548) (@fabriziosestito)
- Add cluster settings api [\#546](https://github.com/trento-project/trento/pull/546) (@nelsonkopliku)
- Add support to scroll to anchors in the catalog [\#539](https://github.com/trento-project/trento/pull/539) (@rtorrero)
- Switch to projected sapsystems list [\#538](https://github.com/trento-project/trento/pull/538) (@fabriziosestito)
- Refactor cluster tags api with the new clusters service [\#533](https://github.com/trento-project/trento/pull/533) (@rtorrero)
- Refactor cluster SIDs in SID \(model, entity, projector, service, handlers\) [\#531](https://github.com/trento-project/trento/pull/531) (@fabriziosestito)
- Cleanup cluster detail leftovers [\#529](https://github.com/trento-project/trento/pull/529) (@fabriziosestito)
- Fix hosts next handler/service naming leftovers [\#528](https://github.com/trento-project/trento/pull/528) (@fabriziosestito)
- Switch to projected HANA cluster detail view [\#524](https://github.com/trento-project/trento/pull/524) (@fabriziosestito)
- Change the Telemetry apiHost to a dummy service [\#523](https://github.com/trento-project/trento/pull/523) (@nelsonkopliku)
- Adds a link to a piece of documentation on SSL certificates [\#522](https://github.com/trento-project/trento/pull/522) (@nelsonkopliku)
- Removes skipping certificate check from collector client [\#520](https://github.com/trento-project/trento/pull/520) (@nelsonkopliku)
- Improve changes file generation [\#519](https://github.com/trento-project/trento/pull/519) (@arbulu89)
- Switch to hosts next [\#517](https://github.com/trento-project/trento/pull/517) (@fabriziosestito)
- Add SAPSystems to host \(next\) [\#511](https://github.com/trento-project/trento/pull/511) (@fabriziosestito)
- Move checks structs to entities and remove old checks in catalog [\#509](https://github.com/trento-project/trento/pull/509) (@arbulu89)
- Add premium badge to the premium checks in catalog [\#504](https://github.com/trento-project/trento/pull/504) (@arbulu89)
- Include missing references on some checks.  [\#499](https://github.com/trento-project/trento/pull/499) (@diegoakechi)
- Fixed flaky test on telemetry engine [\#498](https://github.com/trento-project/trento/pull/498) (@nelsonkopliku)
- Add EULA acceptance middleware and UI [\#497](https://github.com/trento-project/trento/pull/497) (@dottorblaster)
- Project HANA Cluster Details [\#488](https://github.com/trento-project/trento/pull/488) (@fabriziosestito)
- Publish host telemetry to telemetry service [\#487](https://github.com/trento-project/trento/pull/487) (@nelsonkopliku)
- remove the 'trento-' prefix from container artifacts [\#486](https://github.com/trento-project/trento/pull/486) (@stefanotorresi)
- Identify trento installation [\#472](https://github.com/trento-project/trento/pull/472) (@nelsonkopliku)
- Enable SAPSystems projection [\#462](https://github.com/trento-project/trento/pull/462) (@dottorblaster)
- Ci goodies 2: The Comeback [\#459](https://github.com/trento-project/trento/pull/459) (@rtorrero)
- Add heartbeat to agent [\#444](https://github.com/trento-project/trento/pull/444) (@fabriziosestito)

## [0.6.0](https://github.com/trento-project/trento/tree/0.6.0) (2021-11-18)

[Full Changelog](https://github.com/trento-project/trento/compare/0.5.0...0.6.0)

### Added

- Introducing config files for trento \(web|agent|runner\) [\#423](https://github.com/trento-project/trento/issues/423)
- Prune data collection events older than X days [\#399](https://github.com/trento-project/trento/issues/399)
- Refactor runner config [\#376](https://github.com/trento-project/trento/issues/376)
- Add context to the web app and projector worker goroutine to handle graceful stop [\#351](https://github.com/trento-project/trento/issues/351)
- Include Ansible output in the Runner console logs [\#322](https://github.com/trento-project/trento/issues/322)
- Detect aws and gcp clouds in the agent [\#466](https://github.com/trento-project/trento/pull/466) (@arbulu89)
- Order check groups by name in catalog endpoint [\#465](https://github.com/trento-project/trento/pull/465) (@dottorblaster)
- Address HANA cluster settings modal quirks [\#463](https://github.com/trento-project/trento/pull/463) (@dottorblaster)
- Introduce toasts and use them in the checks settings [\#439](https://github.com/trento-project/trento/pull/439) (@dottorblaster)
- Allow Agent rpm to run with config file [\#434](https://github.com/trento-project/trento/pull/434) (@nelsonkopliku)
- Project telemetry data [\#418](https://github.com/trento-project/trento/pull/418) (@nelsonkopliku)
- Refactored Host Discovery to publish more extensive information [\#403](https://github.com/trento-project/trento/pull/403) (@nelsonkopliku)
- Store checks metadata in the DB instead of ARA [\#402](https://github.com/trento-project/trento/pull/402) (@arbulu89)
- Refactor runner cmd config [\#393](https://github.com/trento-project/trento/pull/393) (@arbulu89)
- Implement new API for the checks connection data, storing the data in the DB [\#391](https://github.com/trento-project/trento/pull/391) (@arbulu89)
- Cluster checks selection implemented using the DB [\#375](https://github.com/trento-project/trento/pull/375) (@arbulu89)
- Create the client side api code for check selection [\#369](https://github.com/trento-project/trento/pull/369) (@arbulu89)
- Split Web api code in files [\#368](https://github.com/trento-project/trento/pull/368) (@arbulu89)
- Refactor version checks [\#366](https://github.com/trento-project/trento/pull/366) (@aleksei-burlakov)
- Agent publishes cluster discovery [\#361](https://github.com/trento-project/trento/pull/361) (@nelsonkopliku)
- Add check selection api to the server [\#357](https://github.com/trento-project/trento/pull/357) (@arbulu89)
- Refactor ansible inventory creation removing consultemplate [\#347](https://github.com/trento-project/trento/pull/347) (@arbulu89)
- Add secure data collector endpoint [\#341](https://github.com/trento-project/trento/pull/341) (@fabriziosestito)

### Fixed

- Clusters order in the clusters list page changes over time [\#455](https://github.com/trento-project/trento/issues/455)
- Swagger page calls are broken [\#453](https://github.com/trento-project/trento/issues/453)
- Check Catalog code blocks overflow too much and break the collapsable UX [\#365](https://github.com/trento-project/trento/issues/365)
- The links on the low part of the Home page point to outdated markdown files [\#358](https://github.com/trento-project/trento/issues/358)
- Added scrollbars to the codeblocks [\#471](https://github.com/trento-project/trento/pull/471) (@MMuschner)
- Include the consul-config-dir init in the agent install script [\#470](https://github.com/trento-project/trento/pull/470) (@arbulu89)
- Pin pyparsing to ~2.0 version to avoid issues in runner container [\#468](https://github.com/trento-project/trento/pull/468) (@arbulu89)
- Add conditional in the spec file to detect TW and otherwise avoid missing macro [\#464](https://github.com/trento-project/trento/pull/464) (@rtorrero)
- Remove /api prefix from swagger api docstrings [\#457](https://github.com/trento-project/trento/pull/457) (@arbulu89)
- Fix cloud os user name retrieval in the runner side [\#417](https://github.com/trento-project/trento/pull/417) (@arbulu89)
- Restablish consul-config-dir usage [\#409](https://github.com/trento-project/trento/pull/409) (@arbulu89)
- check for having elements in a slice before accessing those [\#407](https://github.com/trento-project/trento/pull/407) (@nelsonkopliku)
- Typo fixes in home.html.tmpl [\#397](https://github.com/trento-project/trento/pull/397) (@MMuschner)
- Fixed outdated links [\#390](https://github.com/trento-project/trento/pull/390) (@MMuschner)
- Use ElementsMatch to avoid randomly ordered maps in test [\#362](https://github.com/trento-project/trento/pull/362) (@arbulu89)

### Closed issues

- Limit concurrency to 1 in the CI [\#389](https://github.com/trento-project/trento/issues/389)
- Make database-dependant tests skippable. [\#364](https://github.com/trento-project/trento/issues/364)
- Restore manual triggering of the CI/CD [\#354](https://github.com/trento-project/trento/issues/354)
- Refactor the ProjectorRegistry in a separate file [\#353](https://github.com/trento-project/trento/issues/353)

### Other Changes

- Add cluster list smoke test [\#460](https://github.com/trento-project/trento/pull/460) (@fabriziosestito)
- Add a log telling the configuration file being used [\#458](https://github.com/trento-project/trento/pull/458) (@nelsonkopliku)
- Move concurrency to workflow-level [\#456](https://github.com/trento-project/trento/pull/456) (@rtorrero)
- Uniform runner config loading to config files [\#452](https://github.com/trento-project/trento/pull/452) (@nelsonkopliku)
- Cleanup filtering & pagination [\#451](https://github.com/trento-project/trento/pull/451) (@fabriziosestito)
- Cleanup clusters service [\#450](https://github.com/trento-project/trento/pull/450) (@fabriziosestito)
- fixed agent config file creation [\#445](https://github.com/trento-project/trento/pull/445) (@nelsonkopliku)
- Revert "More ci goodies" [\#443](https://github.com/trento-project/trento/pull/443) (@fabriziosestito)
- Add Cypress and add a first smoke test [\#442](https://github.com/trento-project/trento/pull/442) (@dottorblaster)
- Add heartbeat endpoint [\#441](https://github.com/trento-project/trento/pull/441) (@fabriziosestito)
- More ci goodies [\#438](https://github.com/trento-project/trento/pull/438) (@stefanotorresi)
- Fix selected checks in settings endpoint being deserialized to null [\#437](https://github.com/trento-project/trento/pull/437) (@dottorblaster)
- Revert 435 [\#436](https://github.com/trento-project/trento/pull/436) (@stefanotorresi)
- Fix docker build and makefile introducing a new go-build target [\#435](https://github.com/trento-project/trento/pull/435) (@dottorblaster)
- Add hosts projector [\#433](https://github.com/trento-project/trento/pull/433) (@fabriziosestito)
- Use Agent UUID from machine id [\#432](https://github.com/trento-project/trento/pull/432) (@nelsonkopliku)
- Opened Resource cleanup in fixtures [\#431](https://github.com/trento-project/trento/pull/431) (@nelsonkopliku)
- Forcing refreshing updated\_at information of the HostTelemetry [\#429](https://github.com/trento-project/trento/pull/429) (@nelsonkopliku)
- Use a different package to extract system information during host discovery [\#428](https://github.com/trento-project/trento/pull/428) (@nelsonkopliku)
- Add host health aggregation matrix to the dev notes [\#426](https://github.com/trento-project/trento/pull/426) (@stefanotorresi)
- Using correct mocked value for discovered cloud [\#424](https://github.com/trento-project/trento/pull/424) (@nelsonkopliku)
- Switch clusters page on projected read models [\#422](https://github.com/trento-project/trento/pull/422) (@nelsonkopliku)
- Fix `rpm` package version in Dockerfile to 0.0.2 [\#421](https://github.com/trento-project/trento/pull/421) (@dottorblaster)
- More helm configuration settings + testing [\#420](https://github.com/trento-project/trento/pull/420) (@stefanotorresi)
- Use the correct test request constructor [\#416](https://github.com/trento-project/trento/pull/416) (@stefanotorresi)
- Prune old events [\#414](https://github.com/trento-project/trento/pull/414) (@fabriziosestito)
- Makefile updates [\#413](https://github.com/trento-project/trento/pull/413) (@stefanotorresi)
- Update Swagger usage [\#411](https://github.com/trento-project/trento/pull/411) (@stefanotorresi)
- fixed contuing on empty attachedDatabases [\#408](https://github.com/trento-project/trento/pull/408) (@nelsonkopliku)
- fix CI issues [\#406](https://github.com/trento-project/trento/pull/406) (@stefanotorresi)
- Fix runner deployment [\#404](https://github.com/trento-project/trento/pull/404) (@fabriziosestito)
- Cleanup projectors and handlers [\#401](https://github.com/trento-project/trento/pull/401) (@fabriziosestito)
- Change postgresql trento dev default port from 32432 to 5432 [\#398](https://github.com/trento-project/trento/pull/398) (@fabriziosestito)
- add env vars prefix [\#395](https://github.com/trento-project/trento/pull/395) (@stefanotorresi)
- skip tests instead of panicking when db is not available [\#392](https://github.com/trento-project/trento/pull/392) (@stefanotorresi)
- Add collector host/port to the agent config loading function [\#388](https://github.com/trento-project/trento/pull/388) (@fabriziosestito)
- Small improvement of CI [\#387](https://github.com/trento-project/trento/pull/387) (@nelsonkopliku)
- added DATA\_COLLECTOR\_ENABLED=true to the CI [\#384](https://github.com/trento-project/trento/pull/384) (@nelsonkopliku)
- Expose collector service in helm chart [\#379](https://github.com/trento-project/trento/pull/379) (@fabriziosestito)
- Add default collector port configuration [\#378](https://github.com/trento-project/trento/pull/378) (@fabriziosestito)
- Add enable mtls condition before building tls config [\#377](https://github.com/trento-project/trento/pull/377) (@fabriziosestito)
- Refactor web/agent configuration [\#373](https://github.com/trento-project/trento/pull/373) (@fabriziosestito)
- Add annotation to always roll deployments [\#372](https://github.com/trento-project/trento/pull/372) (@fabriziosestito)
- Make all discoveries able to publish data to the Collector [\#371](https://github.com/trento-project/trento/pull/371) (@nelsonkopliku)
- update install-agent to get agent source from fork [\#370](https://github.com/trento-project/trento/pull/370) (@nelsonkopliku)
- Refactor web app/projectors pool to handle graceful shutdown and drain [\#367](https://github.com/trento-project/trento/pull/367) (@fabriziosestito)
- Migrate cluster settings modal to React [\#363](https://github.com/trento-project/trento/pull/363) (@dottorblaster)
- Allow CI to install forked versions of the agent when running from a fork [\#359](https://github.com/trento-project/trento/pull/359) (@rtorrero)
- Added how-to about adding checks [\#290](https://github.com/trento-project/trento/pull/290) (@MMuschner)

## [0.5.0](https://github.com/trento-project/trento/tree/0.5.0) (2021-10-20)

[Full Changelog](https://github.com/trento-project/trento/compare/0.4.1...0.5.0)

### Added

- Add a test for ApiClusterCheckResultsHandler [\#304](https://github.com/trento-project/trento/issues/304)
- Allow install-server script to fetch from different repo owners [\#342](https://github.com/trento-project/trento/pull/342) (@rtorrero)
- Add HANA replication state in the Databases list view [\#338](https://github.com/trento-project/trento/pull/338) (@arbulu89)
- Add DB information in the SAP systems list page [\#334](https://github.com/trento-project/trento/pull/334) (@arbulu89)
- Add the possibility to filter the checks table [\#333](https://github.com/trento-project/trento/pull/333) (@dottorblaster)
- Compare corosync.conf across the nodes [\#331](https://github.com/trento-project/trento/pull/331) (@aleksei-burlakov)
- Cluster checks table makeover: achieve a hierarchical view [\#329](https://github.com/trento-project/trento/pull/329) (@dottorblaster)
- Prevent installing the server if firewalld is detected [\#324](https://github.com/trento-project/trento/pull/324) (@fabriziosestito)
- Refactor and cleanup web tests [\#323](https://github.com/trento-project/trento/pull/323) (@fabriziosestito)
- Discover a globally unique SAP system ID [\#311](https://github.com/trento-project/trento/pull/311) (@arbulu89)
- Add PostgreSQL [\#306](https://github.com/trento-project/trento/pull/306) (@fabriziosestito)
- Add hana database entry sidebar [\#303](https://github.com/trento-project/trento/pull/303) (@arbulu89)
- Cluster health details view makeover [\#291](https://github.com/trento-project/trento/pull/291) (@dottorblaster)
- Runner - Create new check ids system [\#282](https://github.com/trento-project/trento/pull/282) (@arbulu89)

### Fixed

- Checks result modal node column empty on unreachable [\#340](https://github.com/trento-project/trento/issues/340)
- Trento runner container in K3S is leaving ssh defunc processes [\#326](https://github.com/trento-project/trento/issues/326)
- The check 2C2D43 \(2.2.8\) is blocking the whole runner's ansible execution [\#325](https://github.com/trento-project/trento/issues/325)
- Tags autocomplete options shows duplicated tags [\#317](https://github.com/trento-project/trento/issues/317)
- HANA Databases page misses tags API, so tags are broken [\#308](https://github.com/trento-project/trento/issues/308)
- Make the new cluster health details play well with unreachable hosts [\#305](https://github.com/trento-project/trento/issues/305)
- container images are not built with the correct version constant [\#301](https://github.com/trento-project/trento/issues/301)
- Fix apparmor\_parser requirement check [\#355](https://github.com/trento-project/trento/pull/355) (@fabriziosestito)
- Show properly the checks when some node is unreachable by ansible [\#343](https://github.com/trento-project/trento/pull/343) (@arbulu89)
- Add a runner\_on\_skipped hook into the callback module of the runner [\#337](https://github.com/trento-project/trento/pull/337) (@dottorblaster)
- Fix check 2.2.8 updating the sudo calls to not block the runner in k3s [\#328](https://github.com/trento-project/trento/pull/328) (@arbulu89)
- Use tini in the runner container to remove ssh zombie processes [\#327](https://github.com/trento-project/trento/pull/327) (@arbulu89)
- Store hosts reachable state during the runner execution [\#320](https://github.com/trento-project/trento/pull/320) (@arbulu89)
- Remove duplicated tags in getTags function [\#319](https://github.com/trento-project/trento/pull/319) (@arbulu89)
- Fix tags usage in the databases view [\#318](https://github.com/trento-project/trento/pull/318) (@arbulu89)
- Get correct version to OBS using git release tags [\#313](https://github.com/trento-project/trento/pull/313) (@rtorrero)
- Fetch tags before building container images [\#302](https://github.com/trento-project/trento/pull/302) (@fabriziosestito)

### Other Changes

- Add minimum ansible version in the documentation and dockerfile [\#336](https://github.com/trento-project/trento/pull/336) (@arbulu89)
- update readme according to latest changes and revisit its structure [\#312](https://github.com/trento-project/trento/pull/312) (@stefanotorresi)
- Trigger `obs-commit` job also for releases [\#310](https://github.com/trento-project/trento/pull/310) (@rtorrero)
- Avoid using "rolling" as version and use the previous version instead [\#309](https://github.com/trento-project/trento/pull/309) (@rtorrero)
- updates env documentation [\#307](https://github.com/trento-project/trento/pull/307) (@nelsonkopliku)
- add more guidelines to the release how-to [\#300](https://github.com/trento-project/trento/pull/300) (@stefanotorresi)
- Disable gin debug logging on demand [\#299](https://github.com/trento-project/trento/pull/299) (@dottorblaster)
- Add apparmor pre-requirement in the install-server script [\#268](https://github.com/trento-project/trento/pull/268) (@dottorblaster)

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
