# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

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