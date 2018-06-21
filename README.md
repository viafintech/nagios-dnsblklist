# Dnsblklist Check

Dnsblklist Check is a golang-based blacklist ip (mail-server ip) check.

## Installation
Just clone the repo and built it with

    make all # will build a linux and mac executable as long as docker is installed

### Resolver

Depending on your system it could be possible that you need to adjust your [cgo resolver](https://golang.org/pkg/net/).

    export GODEBUG=netdns=cgo

## Usage
    nagios-dnsblklist --help #or
    nagios-dnsblklist <subcommand> --help

## Configuration file

A default configuration file could look like:

```Yaml
blacklistServers:
  - 'black.list.server1'
  - 'black.list.server2'
timeout: 2
verbosity: 0
suppresscrit: false
```

## Known issues


## Contributing
This is an open source project and your contribution is very much appreciated.

1. Check for open issues or open a fresh issue to start a discussion around a feature idea or a bug.
2. Fork the repository on Github and make your changes on the **develop** branch (or branch off of it).
3. Send a pull request (with the **develop** branch as the target).


## Changelog
See [CHANGELOG.md](CHANGELOG.md)

## License
Dnsblklist is available under the MIT license. See the [LICENSE](LICENSE) file for more info.