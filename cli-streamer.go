package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// this cross-platform tool allows to execute full command from shell with possibility
// to specify the execution timeout with others options such as streaming output to multiple
// destinations like multiple files. This timeout approach uses only time.After & goroutines
// so this tool works and compiles without any problem on golang version < 1.7.

// executor is the core function of this tool. same behavior can be achieved with builtin
// os/exec CommandContext function which is available from version >= 1.7.
func executor(cmd *exec.Cmd, timeout time.Duration, quit <-chan struct{}) {
	var err error
	// this start the task asynchronously.
	err = cmd.Start()
	if err != nil {
		// failed to start the task. no need to continue
		log.Printf("failed to start the task - errmsg : %v", err)
		return
	}
	log.Printf("task started under process id [%d]\n", cmd.Process.Pid)
	// goroutine to handle the blocking behavior of wait func - channel used to notify.
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	// watch on both channels and handle the case which hits/triggers first.
	select {
	case <-quit:
		// best effort to kill process and leave.
		cmd.Process.Kill()
		return
	// start the timer and keep watching until expired.
	case <-time.After(timeout):
		// timeout reached - so try to kill the job process.
		log.Printf("task execution timeout reached - killing the process id [%d]\n", cmd.Process.Pid)
		// kill the process and exit from this function.
		if err = cmd.Process.Kill(); err != nil {
			log.Printf("task execution timeout reached - failed to kill process id [%d] - errmsg: %v\n", cmd.Process.Pid, err)
		} else {
			log.Printf("task execution timeout reached - succeeded to kill process id [%d]\n", cmd.Process.Pid)
		}

		return
	case err = <-done:
		// task execution completed [cmd.wait func] - check if for error.
		if err != nil {
			fmt.Printf("task completed with failure - errmsg : %v", err)
		}
		return
		// if needed to dump the buffer content to console
		// cmdstdout, _ := cmd.StdoutPipe()
		// data, _ := ioutil.ReadAll(bufio.NewReader(cmdstdout))
		// fmt.Println()
		// fmt.Printf(string(data))
		// cmdstdout.Close()
		// or inside main func - acheive the same with fmt.Println(result.String())
	}
}

func formatSyntax(task string) *exec.Cmd {
	var cmd *exec.Cmd
	// command syntax for windows platform.
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", task)
	} else {
		// set default shell to use on linux.
		shell := "/bin/sh"
		// load shell name from env variable.
		if os.Getenv("SHELL") != "" {
			shell = os.Getenv("SHELL")
		}
		// syntax for linux-based platforms.
		cmd = exec.Command(shell, "-c", task)
	}

	return cmd
}

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

	if *taskPtr == "" {
		// no command provided - abort.
		flag.Usage()
		return
	}

	// command to execute provided so format syntax based on platform.
	// we can use as demo the task "netstat -n 2"
	cmd := formatSyntax(*taskPtr)

	// 0 or negative value provided - reset to 1hr = 3600 secs.
	if *timeoutPtr <= 0 {
		*timeoutPtr = 3600
	}

	// build the multi destinations io writer.
	// by default we stream in memory buffer the combined output.
	// resultbuffer := &bytes.Buffer{}
	// outWriters := []io.Writer{resultbuffer}
	outWriters := []io.Writer{}

	if *consolePtr {
		// user wants output being displayed at terminal.
		outWriters = append(outWriters, os.Stdout)
	}

	// open or create the day file to stream output of the task execution - default to nil.
	if *savePtr {
		// user wants to save output to daily file so try open or create it.
		starttime := time.Now()
		dailyFile, err := os.OpenFile(fmt.Sprintf("outputs-%d%02d%02d.txt", starttime.Year(), starttime.Month(), starttime.Day()), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Printf("failed to create or open saving file for the task - errmsg : %v", err)
			return
		}
		// add to daily file to io writers list.
		outWriters = append(outWriters, dailyFile)
	}

	// retrieve list of filenames based on space.
	filenames := strings.Fields(*filesPtr)

	if len(filenames) > 0 {
		// other destinations files were mentionned.
		for _, filename := range filenames {
			// open or create each file and add to the destinations list.
			dstfile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Printf("failed to create or open destination file [%s] - errmsg : %v", err, filename)
				continue
			}
			outWriters = append(outWriters, dstfile)
		}
	}

	// use standard console output if no output specified.
	if len(outWriters) == 0 {
		outWriters = append(outWriters, os.Stdout)
	}

	// to display on console add also os.Stdout.
	// create the final multi-destination writer.
	outWr := io.MultiWriter(outWriters...)
	// set the command to use it as its standard io.
	cmd.Stdout, cmd.Stderr = outWr, outWr

	// build the timeout duration.
	timeout := time.Duration(*timeoutPtr) * time.Second
	// wanted to display output each 5 secs one line at a time.

	// run a goroutine to handle exit signals.
	quit := make(chan struct{}, 1)
	go handlesignal(quit)
	// below function is blocking with select statement.
	executor(cmd, timeout, quit)

	// lets reset the memory buffer and use built-in feature.
	// resultbuffer.Reset()
}

const version = "This tool is <cli-streamer> • version 1.0 By Jerome AMON"

const usage = `Usage:
    
    cli-streamer [-task <quoted-command-to-execute>] [-timeout <execution-deadline-in-seconds>] [-files <filenames-to-stream-output>] [-save] [-console]

Subcommands:
    version    Display the current version of this tool.
    help       Display the help - how to use this tool.


Options:
    -task      Specify the command with its arguments to run in quoted format.
    -timeout   Specify the number of seconds to allow the task to be running.
    -files     Specify all filenames to stream the output of the execution.
    -save      If present then execution output must also be saved in daily file.
    -console   If present then execution output will also be displayed on terminal.
    

Arguments:
    quoted-command-to-execute      complete command to execute from the shell.
    execution-deadline-in-seconds  number of seconds to have the task running.
    filenames-to-stream-output     space separed filenames where to stream output.

You have to provide at least one mandatory argument value [-task]. The list of files
where to have the output duplicated (in real-time) must be in one word and space
separed. To have the output be saved into a daily file named outputs-<year><month><day>,
just add -save flag when launching the program. Upcoming version will add capabilities
to mention multiple tasks and each task with its own destinations output filenames.
In the meantime please see below examples for current version 1.0 :


Examples:
	$ cli-streamer version
	$ cli-streamer --help
    $ cli-streamer -task "netstat -n 2 | findstr ESTAB" -timeout 180 -files "a.txt b.txt" -save
    $ cli-streamer -task "ping 127.0.0.1 -t" -timeout 3600 -files "ping.txt" --console
    $ cli-streamer -task "journalctl -f | grep <xx>" -timeout 120 -files "proclog.txt" --save
    $ cli-streamer -task "tail -f /var/log/syslog" -timeout 3600 -files "syslog.txt"`
