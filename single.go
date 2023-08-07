package main

import (
	"log"
	"os"
	"strings"
)

// contains routines to execute single task defined through the command line.

func ProcessSingleTask(quit <-chan struct{}, taskPtr *string, filesPtr *string, timeoutPtr *int, savePtr, consolePtr *bool) {
	task := Task{}
	task.Task = *taskPtr

	// 0 or negative value provided - reset to 1hr = 3600 secs.
	if *timeoutPtr <= 0 {
		*timeoutPtr = 3600
	}
	task.Timeout = *timeoutPtr

	if *consolePtr {
		// user wants output being displayed at terminal.
		task.Console = true
	}

	if *savePtr {
		// user wants to save output to daily file.
		task.Save = true
	}

	// retrieve list of filenames based on space.
	filenames := strings.Fields(*filesPtr)
	if len(filenames) > 0 {
		// user mentionned other destinations file(s).
		task.Files = append(task.Files, filenames...)
	}

	// use standard console output if no output specified.
	if !task.Console && len(task.Files) == 0 && !task.Save {
		task.Console = true
	}

	// build output stream and command to execute.
	out, err := task.IOWriter()
	if err != nil {
		log.Printf("failed to build task output stream: %v\n", err)
		os.Exit(1)
	}
	cmd := task.ToExecCommand(out)
	task.Execute(cmd, quit)
}
