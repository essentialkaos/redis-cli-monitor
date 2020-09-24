<p align="center"><a href="#readme"><img src="https://gh.kaos.st/redis-cli-monitor.svg"/></a></p>

<p align="center">
  <a href="https://travis-ci.com/essentialkaos/redis-cli-monitor"><img src="https://travis-ci.com/essentialkaos/redis-cli-monitor.svg"></a>
  <a href="https://github.com/essentialkaos/redis-cli-monitor/actions?query=workflow%3ACodeQL"><img src="https://github.com/essentialkaos/redis-cli-monitor/workflows/CodeQL/badge.svg" /></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/redis-cli-monitor"><img src="https://goreportcard.com/badge/github.com/essentialkaos/redis-cli-monitor"></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-redis-cli-monitor-master"><img alt="codebeat badge" src="https://codebeat.co/badges/9503a6f8-c9da-4057-ae44-b079686bcc13" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#usage-demo">Usage demo</a> • <a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>

<br/>

Tiny Redis client for renamed `MONITOR` commands.

### Usage demo

[![demo](https://gh.kaos.st/redis-cli-monitor-200.gif)](#usage-demo)

### Installation

#### From source

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)):

```
git config --global http.https://pkg.re.followRedirects true
```

To build the `redis-cli-monitor` from scratch, make sure you have a working Go 1.14+ workspace ([instructions](https://golang.org/doc/install)), then:

```
go get github.com/essentialkaos/redis-cli-monitor
```

If you want to update `redis-cli-monitor` to latest stable release, do:

```
go get -u github.com/essentialkaos/redis-cli-monitor
```

#### From [ESSENTIAL KAOS Public Repository](https://yum.kaos.st)

```bash
sudo yum install -y https://yum.kaos.st/get/$(uname -r).rpm
sudo yum install redis-cli-monitor
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and OS X from [EK Apps Repository](https://apps.kaos.st/redis-cli-monitor/latest).

To install the latest prebuilt version, do:

```bash
bash <(curl -fsSL https://apps.kaos.st/get) redis-cli-monitor
```

### Usage

```
Usage: redis-cli-monitor {options} command-name

Options

  --host, -h ip/host         Server hostname (127.0.0.1 by default)
  --port, -p port            Server port (6379 by default)
  --db, -n db                Database number
  --raw, -r                  Print raw data
  --password, -a password    Password to use when connecting to the server
  --timeout, -t 1-300        Connection timeout in seconds (3 by default)
  --no-color, -nc            Disable colors in output
  --help                     Show this help message
  --version, -v              Show version

Examples

  redis-cli-monitor --host 192.168.0.123 --port 6821 --timeout 15 RENAMED_MONITOR
  Execute "RENAMED_MONITOR" command on 192.168.0.123:6821 with 15 sec timeout

  redis-cli-monitor -p 6378 -a MySuppaPassword1234 RENAMED_MONITOR
  Execute "RENAMED_MONITOR" command on 127.0.0.1:6378 with password "MySuppaPassword1234"

```

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![Build Status](https://travis-ci.com/essentialkaos/redis-cli-monitor.svg?branch=master)](https://travis-ci.com/essentialkaos/redis-cli-monitor) |
| `develop` | [![Build Status](https://travis-ci.com/essentialkaos/redis-cli-monitor.svg?branch=develop)](https://travis-ci.com/essentialkaos/redis-cli-monitor) |

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
