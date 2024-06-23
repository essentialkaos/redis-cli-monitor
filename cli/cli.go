package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/strutil"
	"github.com/essentialkaos/ek/v12/support"
	"github.com/essentialkaos/ek/v12/support/deps"
	"github.com/essentialkaos/ek/v12/terminal"
	"github.com/essentialkaos/ek/v12/terminal/tty"
	"github.com/essentialkaos/ek/v12/timeutil"
	"github.com/essentialkaos/ek/v12/usage"
	"github.com/essentialkaos/ek/v12/usage/completion/bash"
	"github.com/essentialkaos/ek/v12/usage/completion/fish"
	"github.com/essentialkaos/ek/v12/usage/completion/zsh"
	"github.com/essentialkaos/ek/v12/usage/man"
	"github.com/essentialkaos/ek/v12/usage/update"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Application info
const (
	APP  = "Redis CLI Monitor"
	VER  = "2.2.3"
	DESC = "Tiny Redis client for renamed MONITOR commands"
)

// Supported command line options
const (
	OPT_HOST     = "h:host"
	OPT_PORT     = "p:port"
	OPT_DB       = "n:db"
	OPT_RAW      = "r:raw"
	OPT_AUTH     = "a:password"
	OPT_TIMEOUT  = "t:timeout"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "help"
	OPT_VER      = "v:version"

	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_HOST:     {Type: options.MIXED, Value: "127.0.0.1"},
	OPT_PORT:     {Value: "6379"},
	OPT_DB:       {Type: options.INT, Min: 0, Max: 999},
	OPT_TIMEOUT:  {Type: options.INT, Value: 3, Min: 1, Max: 300},
	OPT_RAW:      {Type: options.BOOL},
	OPT_AUTH:     {},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.MIXED},

	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// colorTagApp contains color tag for app name
var colorTagApp string

// colorTagVer contains color tag for app version
var colorTagVer string

// conn is connection to Redis
var conn net.Conn

// useRawOutput is raw output flag
var useRawOutput bool

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main application function
func Run(gitRev string, gomod []byte) {
	preConfigureUI()

	args, errs := options.Parse(optMap)

	if !errs.IsEmpty() {
		terminal.Error("Options parsing errors:")
		terminal.Error(errs.String())
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Collect(APP, VER).
			WithRevision(gitRev).
			WithDeps(deps.Extract(gomod)).
			WithApps(getRedisVersionInfo()).
			Print()
		os.Exit(0)
	case options.GetB(OPT_HELP), options.GetS(OPT_HOST) == "true", len(args) == 0:
		genUsage().Print()
		os.Exit(0)
	}

	connectToRedis()
	monitor(args.Get(0).String())
}

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if !tty.IsTTY() {
		fmtc.DisableColors = true
	}

	switch {
	case fmtc.IsTrueColorSupported():
		colorTagApp, colorTagVer = "{*}{#DC382C}", "{#A32422}"
	case fmtc.Is256ColorsSupported():
		colorTagApp, colorTagVer = "{*}{#160}", "{#124}"
	default:
		colorTagApp, colorTagVer = "{r*}", "{r}"
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if options.GetB(OPT_RAW) {
		useRawOutput = true
	}
}

// connectToRedis connects to Redis instance
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

// monitor starts output commands in monitor
func monitor(cmd string) {
	buf := bufio.NewReader(conn)
	dbNum := options.GetS(OPT_DB)

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

			if dbNum != "" && !isDBMatch(str[1:], dbNum) {
				continue
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

// formatCommand formats command and add color codes
func formatCommand(cmd string) {
	sec, _ := strconv.ParseInt(cmd[:10], 10, 64)

	infoStart := strings.IndexRune(cmd, '[')
	infoEnd := strings.IndexRune(cmd, ']')

	if infoStart == -1 || infoEnd == -1 {
		return
	}

	fmtc.Printf(
		"{s-}%s.%s{!} {s}%-26s{!} %s",
		timeutil.Format(time.Unix(sec, 0), "%Y/%m/%d %H:%M:%S"), cmd[11:17],
		cmd[infoStart:infoEnd+1],
		cmd[infoEnd+2:],
	)
}

// isDBMatch returns true if given command executed over DB with given number
func isDBMatch(cmd, dbNum string) bool {
	start := strings.IndexRune(cmd, '[')

	if start == -1 {
		return false
	}

	end := strings.IndexRune(cmd[start:], ' ')

	if end == -1 {
		return false
	}

	end += start

	return dbNum == cmd[start+1:end]
}

// printErrorAndExit print error message and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	terminal.Error(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getRedisVersionInfo returns info about Redis version
func getRedisVersionInfo() support.App {
	cmd := exec.Command("redis-server", "--version")
	output, err := cmd.Output()

	if err != nil {
		return support.App{"Redis", ""}
	}

	ver := strutil.ReadField(string(output), 2, false, ' ')
	ver = strings.TrimLeft(ver, "v=")

	return support.App{"Redis", ver}
}

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Print(bash.Generate(info, "redis-cli-monitor"))
	case "fish":
		fmt.Print(fish.Generate(info, "redis-cli-monitor"))
	case "zsh":
		fmt.Print(zsh.Generate(info, optMap, "redis-cli-monitor"))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(
		man.Generate(
			genUsage(),
			genAbout(""),
		),
	)
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "command-name")

	info.AppNameColorTag = colorTagApp

	info.AddOption(OPT_HOST, "Server hostname {s-}(127.0.0.1 by default){!}", "ip/host")
	info.AddOption(OPT_PORT, "Server port {s-}(6379 by default){!}", "port")
	info.AddOption(OPT_DB, "Database number", "db")
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

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2006,
		Owner:   "ESSENTIAL KAOS",
		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",

		AppNameColorTag: colorTagApp,
		VersionColorTag: colorTagVer,
		DescSeparator:   "{s}—{!}",
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
		about.UpdateChecker = usage.UpdateChecker{
			"essentialkaos/redis-cli-monitor",
			update.GitHubChecker,
		}
	}

	return about
}
