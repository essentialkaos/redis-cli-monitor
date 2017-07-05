package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"pkg.re/essentialkaos/ek.v9/timeutil"
	"strconv"
	"strings"
	"time"

	"pkg.re/essentialkaos/ek.v9/env"
	"pkg.re/essentialkaos/ek.v9/fmtc"
	"pkg.re/essentialkaos/ek.v9/fsutil"
	"pkg.re/essentialkaos/ek.v9/options"
	"pkg.re/essentialkaos/ek.v9/usage"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "Redis CLI Monitor"
	VER  = "2.0.1"
	DESC = "Tiny Redis client for renamed MONITOR commands"
)

const (
	OPT_HOST     = "h:host"
	OPT_PORT     = "p:port"
	OPT_RAW      = "r:raw"
	OPT_AUTH     = "a:password"
	OPT_TIMEOUT  = "t:timeout"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "help"
	OPT_VER      = "v:version"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_HOST:     {Value: "127.0.0.1"},
	OPT_PORT:     {Value: "6379"},
	OPT_TIMEOUT:  {Type: options.INT, Value: 3, Min: 1, Max: 300},
	OPT_RAW:      {Type: options.BOOL},
	OPT_AUTH:     {},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:      {Type: options.BOOL, Alias: "ver"},
}

var conn net.Conn

var useRawOutput bool

// ////////////////////////////////////////////////////////////////////////////////// //

// main is main function
func main() {
	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		for _, err := range errs {
			printError(err.Error())
		}

		os.Exit(1)
	}

	configureUI()

	if options.GetB(OPT_VER) {
		showAbout()
		return
	}

	if options.GetB(OPT_HELP) || len(args) == 0 {
		showUsage()
		return
	}

	connectToRedis()
	monitor(args[0])
}

// configureUI configure user interface
func configureUI() {
	envVars := env.Get()
	term := envVars.GetS("TERM")

	fmtc.DisableColors = true

	if term != "" {
		switch {
		case strings.Contains(term, "xterm"),
			strings.Contains(term, "color"),
			term == "screen":
			fmtc.DisableColors = false
		}
	}

	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if !fsutil.IsCharacterDevice("/dev/stdout") && envVars.GetS("FAKETTY") == "" {
		fmtc.DisableColors = true
		useRawOutput = true
	}

	if options.GetB(OPT_RAW) {
		useRawOutput = true
	}
}

// connectToRedis connect to Redis instance
func connectToRedis() {
	var err error

	host := options.GetS(OPT_HOST) + ":" + options.GetS(OPT_PORT)
	timeout := time.Second * time.Duration(options.GetI(OPT_TIMEOUT))

	conn, err = net.DialTimeout("tcp", host, timeout)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	if options.GetS(OPT_AUTH) == "" {
		return
	}

	_, err = conn.Write([]byte("AUTH " + options.GetS(OPT_AUTH) + "\r\n"))

	if err != nil {
		printErrorAndExit(err.Error())
	}
}

// monitor start outout commands in monitor
func monitor(cmd string) {
	buf := bufio.NewReader(conn)
	conn.Write([]byte(cmd + "\r\n"))

	for {
		str, err := buf.ReadString('\n')

		if len(str) > 0 {
			if str == "+OK\r\n" {
				continue
			}

			if strings.HasPrefix(str, "-ERR ") {
				printErrorAndExit("Redis return error message: " + strings.TrimRight(str[1:], "\r\n"))
			}

			if useRawOutput {
				fmt.Printf("%s", str[1:])
			} else {
				formatCommand(str[1:])
			}
		}

		if err != nil {
			printErrorAndExit(err.Error())
		}
	}
}

// formatCommand format command and add color codes
func formatCommand(cmd string) {
	sec, _ := strconv.ParseInt(cmd[:10], 10, 64)

	infoStart := strings.IndexRune(cmd, '[')
	infoEnd := strings.IndexRune(cmd, ']')

	fmtc.Printf(
		"{s}%s.%s{!} {s-}%s{!} %s",
		timeutil.Format(time.Unix(sec, 0), "%Y/%m/%d %H:%M:%S"), cmd[11:17],
		cmd[infoStart:infoEnd+1],
		cmd[infoEnd+2:],
	)
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// printErrorAndExit print error message and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showUsage print usage info
func showUsage() {
	info := usage.NewInfo("", "command-name")

	info.AddOption(OPT_HOST, "Server hostname {s-}(127.0.0.1 by default){!}", "ip/host")
	info.AddOption(OPT_PORT, "Server port {s-}(6379 by default){!}", "port")
	info.AddOption(OPT_RAW, "Print raw data")
	info.AddOption(OPT_AUTH, "Password to use when connecting to the server", "password")
	info.AddOption(OPT_TIMEOUT, "Connection timeout in seconds {s-}(3 by default){!}", "1-300")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"--host 192.168.0.123 --port 6821 --timeout 15 RENAMED_MONITOR",
		"Execute \"RENAMED_MONITOR\" command on 192.168.0.123:6821 with 15 sec timeout",
	)

	info.AddExample(
		"-p 6378 -a MySuppaPassword1234 RENAMED_MONITOR",
		"Execute \"RENAMED_MONITOR\" command on 127.0.0.1:6378 with password \"MySuppaPassword1234\"",
	)

	info.Render()
}

// showAbout print info about version
func showAbout() {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2006,
		Owner:   "ESSENTIAL KAOS",
		License: "Essential Kaos Open Source License <https://essentialkaos.com/ekol>",
	}

	about.Render()
}
