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
	"time"

	"pkg.re/essentialkaos/ek.v9/fmtc"
	"pkg.re/essentialkaos/ek.v9/options"
	"pkg.re/essentialkaos/ek.v9/usage"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "Redis CLI Monitor"
	VER  = "1.4.0"
	DESC = "Tiny Redis client for renamed MONITOR commands"
)

const (
	OPT_HOST     = "H:host"
	OPT_PORT     = "P:port"
	OPT_AUTH     = "a:password"
	OPT_TIMEOUT  = "t:timeout"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_HOST:     {Value: "127.0.0.1"},
	OPT_PORT:     {Value: "6379"},
	OPT_TIMEOUT:  {Type: options.INT, Value: 3, Min: 1, Max: 300},
	OPT_AUTH:     {},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:      {Type: options.BOOL, Alias: "ver"},
}

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

	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if options.GetB(OPT_VER) {
		showAbout()
		return
	}

	if options.GetB(OPT_HELP) || len(args) == 0 {
		showUsage()
		return
	}

	connect(args[0])
}

// connect connect to Redis and print data
func connect(cmd string) {
	host := options.GetS(OPT_HOST) + ":" + options.GetS(OPT_PORT)
	timeout := time.Second * time.Duration(options.GetI(OPT_TIMEOUT))

	conn, err := net.DialTimeout("tcp", host, timeout)

	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	if options.GetS(OPT_AUTH) != "" {
		conn.Write([]byte("AUTH " + options.GetS(OPT_AUTH) + "\n"))
	}

	conn.Write([]byte(cmd + "\n"))
	connbuf := bufio.NewReader(conn)

	for {
		str, err := connbuf.ReadString('\n')

		if len(str) > 0 {
			if str == "+OK\r\n" {
				continue
			}

			fmt.Printf("%s", str[1:])
		}

		if err != nil {
			printError(err.Error())
			os.Exit(1)
		}
	}
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showUsage print usage info
func showUsage() {
	info := usage.NewInfo("", "command-name")

	info.AddOption(OPT_HOST, "Server hostname", "host")
	info.AddOption(OPT_PORT, "Server port", "port")
	info.AddOption(OPT_AUTH, "Password to use when connecting to the server", "password")
	info.AddOption(OPT_TIMEOUT, "Connection timeout in seconds", "1-300")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"--host 192.168.0.123 --password 6821 --timeout 15 RENAMED_MONITOR",
		"Execute \"RENAMED_MONITOR\" command on 192.168.0.123:6821 with 15 sec timeout",
	)

	info.AddExample(
		"-P 12345 -a MySuppaPassword1234 RENAMED_MONITOR",
		"Execute \"RENAMED_MONITOR\" command on 127.0.0.1:12345 with password \"MySuppaPassword1234\"",
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
