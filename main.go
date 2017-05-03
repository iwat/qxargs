package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "SYNOPSIS")
		fmt.Fprintln(os.Stderr, "    ", os.Args[0], "[flags] command filters ...")
		fmt.Fprintln(os.Stderr, "    ", os.Args[0], "[flags] command commandargs ... -- filters ...")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "DESCRIPTION")
		fmt.Fprintln(os.Stderr, "     Execute command on the list of files that match the given filter.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "     command  The command to be executed with the matches files.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "     filters  There are 2 kinds of filters are supported, name filter and content filter.")
		fmt.Fprintln(os.Stderr, "              Simple string will be treated as file name filter.")
		fmt.Fprintln(os.Stderr, "              String with leading '?' will be treat as content filter.")
		fmt.Fprintln(os.Stderr, "              Multiple filters will be treated as AND.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "EXAMPLES")
		fmt.Fprintln(os.Stderr, "     To execute vim on any file that has go in their name. ")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "         $", os.Args[0], " vim go")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "     To execute vim -p on any file that has go in their name. ")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "         $", os.Args[0], " vim -p -- go")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "     To execute vim -p on any file that has go in their name and has newGrepper in their contents. ")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "         $", os.Args[0], " vim -p -- go ?newGrepper")
		fmt.Fprintln(os.Stderr)
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

	commandArgs := []string(nil)
	queryArgs := []string(nil)

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

	finder := newFinder(false)
	grepper := newGrepper()

	console := newConsole(finder, grepper)
	results := console.update(strings.Join(queryArgs, " "))

	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "no files matched")
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
