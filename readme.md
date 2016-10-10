## Redis CLI Monitor

Tiny redis client for renamed MONITOR commands.

### Installation

````
go get github.com/essentialkaos/redis-cli-monitor
````

### Usage

```
Usage: redis-cli-monitor {options} command-name

Options

  --host, -H ip/host         Server hostname
  --port, -p port            Server port
  --password, -a password    Password to use when connecting to the server
  --timeout, -t 1-300        Connection timeout in seconds
  --help, -h                 Show this help message
  --version, -v              Show version

Examples:

  redis-cli-monitor -h 192.168.0.123 -p 6821 -t 15 RENAMED_MONITOR
  Execute "RENAMED_MONITOR" command on 192.168.0.123:6821 with 15 sec timeout

  redis-cli-monitor -p 12345 -a MySuppaPassword1234 RENAMED_MONITOR
  Execute "RENAMED_MONITOR" command on 127.0.0.1:12345 with password "MySuppaPassword1234"

```

### Prebuilt binaries

You can download prebuilt binaries for Linux and OS X from [EK Apps Repository](https://apps.kaos.io/redis-cli-monitor/).

### Build Status

| Repository | Status |
|------------|--------|
| Stable | [![Build Status](https://travis-ci.org/essentialkaos/redis-cli-monitor.svg?branch=master)](https://travis-ci.org/essentialkaos/redis-cli-monitor) |
| Unstable | [![Build Status](https://travis-ci.org/essentialkaos/redis-cli-monitor.svg?branch=develop)](https://travis-ci.org/essentialkaos/redis-cli-monitor) |

### License

[EKOL](https://essentialkaos.com/ekol)
