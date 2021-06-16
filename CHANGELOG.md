# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0](https://github.com/trento-project/trento/releases/tag/0.2.0) 2021-05-26

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
### Fixed
  - Fix SAP system layout rendering
  - Don't let the app crash on 404s (#97)
  - use the correct path for the global.ini config file of the SAPHanaSR check (#95)

## [0.1.0](https://github.com/trento-project/trento/releases/tag/0.1.0) 2021-05-26

### Added
  - first release of Trento
  - Automated discovery of SAP HANA HA clusters
  - SAP Systems and Instances overview
  - Grouping by Landscapes and Environments
  - Configuration validation for Pacemaker, Corosync, SBD, SAPHanaSR and other + generic SUSE Linux Enterprise for SAP Application OS settings
  - Specific configuration audits for SAP HANA Scale-Up Performance-Optimized
  - scenarios deployed on MS Azure cloud.
