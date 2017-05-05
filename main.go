package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func checkedPrintln(a ...interface{}) {
	if _, err := fmt.Fprintln(os.Stderr, a...); err != nil {
		panic(err)
	}
}

func init() {
	flag.Usage = func() {
		checkedPrintln("SYNOPSIS")
		checkedPrintln("    ", os.Args[0], "[flags] command filters ...")
		checkedPrintln("    ", os.Args[0], "[flags] command commandargs ... -- filters ...")
		checkedPrintln()
		checkedPrintln("DESCRIPTION")
		checkedPrintln("     Execute command on the list of files that match the given filter.")
		checkedPrintln()
		checkedPrintln("     command  The command to be executed with the matches files.")
		checkedPrintln()
		checkedPrintln("     filters  There are 2 kinds of filters are supported, name filter and content filter.")
		checkedPrintln("              Simple string will be treated as file name filter.")
		checkedPrintln("              String with leading '?' will be treat as content filter.")
		checkedPrintln("              Multiple filters will be treated as AND.")
		checkedPrintln()
		checkedPrintln("EXAMPLES")
		checkedPrintln("     To execute vim on any file that has go in their name. ")
		checkedPrintln()
		checkedPrintln("         $", os.Args[0], " vim go")
		checkedPrintln()
		checkedPrintln("     To execute vim -p on any file that has go in their name. ")
		checkedPrintln()
		checkedPrintln("         $", os.Args[0], " vim -p -- go")
		checkedPrintln()
		checkedPrintln("     To execute vim -p on any file that has go in their name and has newGrepper in their contents. ")
		checkedPrintln()
		checkedPrintln("         $", os.Args[0], " vim -p -- go ?newGrepper")
		flag.PrintDefaults()
	}
}

func parseArgs() ([]string, []string) {
	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	dashNdx := -1
	for i, arg := range args {
		if arg == "--" {
			dashNdx = i
			break
		}
	}

	var commandArgs []string
	var queryArgs []string

	if dashNdx == -1 {
		commandArgs = []string{args[0]}
		queryArgs = args[1:]
	} else {
		commandArgs = args[0:dashNdx]
		queryArgs = args[dashNdx+1:]
	}

	return commandArgs, queryArgs
}

func main() {
	flag.Parse()
	commandArgs, queryArgs := parseArgs()

	console := newConsole()
	results := console.update(strings.Join(queryArgs, " "))

	if len(results) == 0 {
		checkedPrintln("no files matched")
		os.Exit(2)
	}

	results, err := console.loop(commandArgs)
	if err != nil {
		panic(err)
	}

	if len(results) == 0 {
		return
	}

	commandPath, err := exec.LookPath(commandArgs[0])
	if err != nil {
		panic(err)
	}

	commandArgs = append(commandArgs, results...)

	err = syscall.Exec(commandPath, commandArgs, os.Environ())
	if err != nil {
		panic(err)
	}
}
