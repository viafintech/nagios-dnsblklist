# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [2.0.7] - 2022-07-25
- Remove sbl./pbl./xbl.spamhaus.org DNSBLs in favor of [zen.spamhaus.org](https://www.spamhaus.org/zen/) which combines answers for all of them

## [2.0.6] - 2022-07-14
- Removed [cbl.abuseat.org](https://www.abuseat.org/)
- Updated Go version to 1.18.4

## [2.0.5] - 2020-09-30
- Removed [dynip.rothen.com](https://www.dnsbl.info/dnsbl-details.php?dnsbl=dynip.rothen.com)
- Updated Go version to 1.15.2

## [2.0.4] - 2019-11-04
- Removed [bl.emailbasura.org](https://www.dnsbl.info/emailbasura-offline.php)
- Updated Go version to 1.13.4

## [2.0.3] - 2019-10-25
- Updated Go version to 1.13.3
- Updated license and copyright information

## [2.0.2] - 2019-06-03
- Removed bl.spamcannibal.org due to unreachability/unrealiability
- Updated Go version from 1.12.4 to 1.12.5

## [2.0.1] - 2019-04-14
- Removed cdl.anti-spam.org.cn due to unreachability/unrealiability
- Updated Go version from 1.10.3 to 1.12.4

## [2.0.0] - 2018-06-21
- Using the cloudflare dns over https api
- Added docker support to build executables for mac and linux

## [1.0.1] - 2018-06-01
- Removed support for bl.spamcannibal.org (down) and rbl.suresupport.com (404)

## [1.0.0] - 2016-11-01
- Check an ip if it is blacklisted
- configure blacklistservers, timeout, verbosity and suppresscrit
- list all the blacklistservers in use
- supply a configuration file (default location or by flag)
- Initial release