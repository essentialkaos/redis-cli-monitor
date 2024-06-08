<p align="center"><a href="#readme"><img src="https://gh.kaos.st/redis-cli-monitor.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/w/redis-cli-monitor/ci"><img src="https://kaos.sh/w/redis-cli-monitor/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/r/redis-cli-monitor"><img src="https://kaos.sh/r/redis-cli-monitor.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/b/redis-cli-monitor"><img src="https://kaos.sh/b/9503a6f8-c9da-4057-ae44-b079686bcc13.svg" alt="codebeat badge" /></a>
  <a href="https://kaos.sh/w/redis-cli-monitor/codeql"><img src="https://kaos.sh/w/redis-cli-monitor/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#usage-demo">Usage demo</a> • <a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>

<br/>

Tiny Redis client for renamed `MONITOR` commands.

### Usage demo

[![demo](https://gh.kaos.st/redis-cli-monitor-200.gif)](#usage-demo)

### Installation

#### From source

To build the `redis-cli-monitor` from scratch, make sure you have a working Go 1.18+ workspace ([instructions](https://go.dev/doc/install)), then:

```
go install github.com/essentialkaos/redis-cli-monitor@latest
```

#### From [ESSENTIAL KAOS Public Repository](https://yum.kaos.st)

```bash
sudo yum install -y https://yum.kaos.st/kaos-repo-latest.el$(grep 'CPE_NAME' /etc/os-release | tr -d '"' | cut -d':' -f5).noarch.rpm
sudo yum install redis-cli-monitor
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and macOS from [EK Apps Repository](https://apps.kaos.st/redis-cli-monitor/latest).

To install the latest prebuilt version, do:

```bash
bash <(curl -fsSL https://apps.kaos.st/get) redis-cli-monitor
```

### Usage

<img src=".github/images/usage.svg" />

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![CI](https://kaos.sh/w/redis-cli-monitor/ci.svg?branch=master)](https://kaos.sh/w/redis-cli-monitor/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/redis-cli-monitor/ci.svg?branch=master)](https://kaos.sh/w/redis-cli-monitor/ci?query=branch:develop) |

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
