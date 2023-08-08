package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// this cross-platform tool allows to execute full command from shell with possibility
// to specify the execution timeout with others options such as streaming output to multiple
// destinations like multiple files. This timeout approach uses only time.After & goroutines
// so this tool works and compiles without any problem on golang version < 1.7.

// handlesignal is a function that process SIGTERM from kill command or CTRL-C or more.
func handlesignal(exit chan<- struct{}) {
	// one signal to be handled.
	sigch := make(chan os.Signal, 1)
	// setup supported exit signals.
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGQUIT,
		syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)

	// block until something comes in.
	<-sigch
	signal.Stop(sigch)
	// then notify executor to stop.
	exit <- struct{}{}
}

func main() {
	// enforce usage of all cores means number of
	// active os threads simultaneously.
	runtime.GOMAXPROCS(runtime.NumCPU())
	// will be triggered to display usage instructions.
	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }

	taskPtr := flag.String("task", "", "full command with its arguments to be executed")
	timeoutPtr := flag.Int("timeout", 3600, "command execution timetout value in seconds")
	// displayPtr := flag.Int("display", 0, "interval between each output line display")
	// list of files names to stream command outputs.
	filesPtr := flag.String("files", "", "filenames to stream execution output")
	// declare the boolean flag save. if mentioned save stream output to daily file.
	savePtr := flag.Bool("save", false, "specify if wanted to stream as well output to daily file")
	// declare the boolean flag save. if mentioned save stream output to daily file.
	consolePtr := flag.Bool("console", false, "specify if wanted to stream as well output to console")

	tasksFilePtr := flag.String("tasksFile", "", "define a file where to load all commands details")
	tasksJsonPtr := flag.String("tasksJson", "", "define a json file where to load all commands details")

	// check for any valid subcommands : version or help
	if len(os.Args) == 2 {
		if os.Args[1] == "version" || os.Args[1] == "--version" || os.Args[1] == "-v" {
			fmt.Fprintf(os.Stderr, "\n%s\n", version)
			os.Exit(0)
		} else {
			fmt.Fprintf(os.Stderr, "\n%s\n", usage)
			os.Exit(0)
		}
	}

	// move on for flag processing.
	flag.Parse()

	if len(*taskPtr) != 0 && len(*tasksFilePtr) != 0 && len(*tasksJsonPtr) != 0 {
		// must use one option.
		flag.Usage()
		return
	}

	// run a goroutine to handle exit signals.
	quit := make(chan struct{}, 1)
	go handlesignal(quit)

	// ption: single task with its details.
	if *taskPtr != "" {
		ProcessSingleTask(quit, taskPtr, filesPtr, timeoutPtr, savePtr, consolePtr)
		return
	}
	// option: load tasks from file.
	if *tasksFilePtr != "" {
		ProcessTasksFile(*tasksFilePtr, quit)
		return
	}

	// option: load tasks from json-encoded file.
	if *tasksJsonPtr != "" {
		ProcessTasksJson(*tasksJsonPtr, quit)
		return
	}

	// no valid option.
	flag.Usage()
}
