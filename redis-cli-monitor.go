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

	"pkg.re/essentialkaos/ek.v7/arg"
	"pkg.re/essentialkaos/ek.v7/usage"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "Redis CLI Monitor"
	VER  = "1.1.0"
	DESC = "Tiny redis client for renamed MONITOR commands"
)

const (
	ARG_HOST    = "H:host"
	ARG_PORT    = "p:port"
	ARG_AUTH    = "a:password"
	ARG_TIMEOUT = "t:timeout"
	ARG_HELP    = "h:help"
	ARG_VER     = "v:version"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var argMap = arg.Map{
	ARG_HOST:    {Value: "127.0.0.1"},
	ARG_PORT:    {Value: "6379"},
	ARG_TIMEOUT: {Type: arg.INT, Value: 3, Min: 1, Max: 300},
	ARG_AUTH:    {},
	ARG_HELP:    {Type: arg.BOOL, Alias: "u:usage"},
	ARG_VER:     {Type: arg.BOOL, Alias: "ver"},
}

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	args, errs := arg.Parse(argMap)

	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err.Error())
		}

		os.Exit(1)
	}

	if arg.GetB(ARG_VER) == true {
		showAbout()
		return
	}

	if arg.GetB(ARG_HELP) == true || len(args) == 0 {
		showUsage()
		return
	}

	connect(args[0])
}

func connect(cmd string) {
	host := arg.GetS(ARG_HOST) + ":" + arg.GetS(ARG_PORT)
	timeout := time.Second * time.Duration(arg.GetI(ARG_TIMEOUT))

	conn, err := net.DialTimeout("tcp", host, timeout)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	if arg.GetS(ARG_AUTH) != "" {
		conn.Write([]byte("AUTH " + arg.GetS(ARG_AUTH) + "\n"))
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
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

func showUsage() {
	info := usage.NewInfo("", "command-name")

	info.AddOption(ARG_HOST, "Server hostname", "ip/host")
	info.AddOption(ARG_PORT, "Server port", "port")
	info.AddOption(ARG_AUTH, "Password to use when connecting to the server", "password")
	info.AddOption(ARG_TIMEOUT, "Connection timeout in seconds", "1-300")
	info.AddOption(ARG_HELP, "Show this help message")
	info.AddOption(ARG_VER, "Show version")

	info.AddExample(
		"-h 192.168.0.123 -p 6821 -t 15 RENAMED_MONITOR",
		"Execute \"RENAMED_MONITOR\" command on 192.168.0.123:6821 with 15 sec timeout",
	)

	info.AddExample(
		"-p 12345 -a MySuppaPassword1234 RENAMED_MONITOR",
		"Execute \"RENAMED_MONITOR\" command on 127.0.0.1:12345 with password \"MySuppaPassword1234\"",
	)

	info.Render()
}

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
