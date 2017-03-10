########################################################################################

.PHONY = fmt all clean deps

########################################################################################

all: redis-cli-monitor

redis-cli-monitor:
	go build redis-cli-monitor.go

deps:
	go get -v pkg.re/essentialkaos/ek.v7

fmt:
	find . -name "*.go" -exec gofmt -s -w {} \;

clean:
	rm -f redis-cli-monitor

########################################################################################

