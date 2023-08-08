package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

type Tasks struct {
	Tasks []Task `json:"tasks" yaml:"tasks" toml:"tasks"`
}

type Task struct {
	Task    string   `json:"task" yaml:"task" toml:"task"`
	Timeout int      `json:"timeout" yaml:"timeout" toml:"timeout"`
	Files   []string `json:"files" yaml:"files" toml:"files"`
	Save    bool     `json:"save" yaml:"save" toml:"save"`
	Console bool     `json:"console" yaml:"console" toml:"console"`
	Web     bool     `json:"web" yaml:"web" toml:"web"`
}

// ToExecCommand turns a task into a ready to run exec/command.
func (task *Task) ToExecCommand(out io.Writer) *exec.Cmd {
	var cmd *exec.Cmd
	// command syntax for windows platform.
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", task.Task)
	} else {
		// set default shell to use on linux.
		shell := "/bin/sh"
		// load shell name from env variable.
		if os.Getenv("SHELL") != "" {
			shell = os.Getenv("SHELL")
		}
		// syntax for linux-based platforms.
		cmd = exec.Command(shell, "-c", task.Task)
	}

	// set the command standard io.
	cmd.Stdout, cmd.Stderr = out, out
	return cmd
}

// IOWriter constructs the multiwriter to be used
// by the commnad as its outputs.
func (task *Task) IOWriter() (io.Writer, error) {
	outWriters := []io.Writer{}

	if task.Console {
		outWriters = append(outWriters, os.Stdout)
	}

	// open or create the day file to stream output of the task execution - default to nil.
	if task.Save {
		now := time.Now()
		dailyFile, err := os.OpenFile(fmt.Sprintf("outputs-%d%02d%02d.txt", now.Year(), now.Month(), now.Day()), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Printf("failed to create or open daily saving file for the task - errmsg : %v", err)
			return nil, err
		}
		outWriters = append(outWriters, dailyFile)
	}

	// open or create each file and add to the outputs stream list.
	if len(task.Files) > 0 {
		for _, filename := range task.Files {
			dstfile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Printf("failed to create or open destination file [%s] - errmsg : %v", err, filename)
				return nil, err
			}
			outWriters = append(outWriters, dstfile)
		}
	}

	return io.MultiWriter(outWriters...), nil
}

// Executor is the core function of this tool. same behavior can be achieved with builtin
// os/exec CommandContext function which is available from version >= 1.7.
func (task *Task) Execute(cmd *exec.Cmd, quit <-chan struct{}) {
	var err error
	// this start the task asynchronously.
	err = cmd.Start()
	if err != nil {
		// failed to start the task. no need to continue
		log.Printf("failed to start the task - errmsg : %v", err)
		return
	}
	pid := cmd.Process.Pid
	log.Printf("[started] [pid:%05d] task: %q\n", pid, task.Task)
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
	case <-time.After(time.Duration(task.Timeout) * time.Second):
		// timeout reached - so try to kill the job process.
		log.Printf("[timeout] [pid:%05d] task: %q timeout: %d\n", pid, task.Task, task.Timeout)
		// kill the process and exit from this function.
		if err = cmd.Process.Kill(); err != nil {
			log.Printf("[timeout] [pid:%05d] task: %q timeout: %d error: %s\n", pid, task.Task, task.Timeout, err.Error())
		}
		return
	case err = <-done:
		// task execution completed [cmd.wait func] - check if for error.
		if err != nil {
			log.Printf("[failed ] [pid:%05d] task: %q error: %s\n", pid, task.Task, err.Error())
		} else {
			log.Printf("[success] [pid:%05d] task: %q\n", pid, task.Task)
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

// Worker consumes available tasks from the queue and execute them.
func TaskWorker(id int, wg *sync.WaitGroup, quit <-chan struct{}, tasksQueue <-chan *Task) {
	defer wg.Done()
	for task := range tasksQueue {
		out, err := task.IOWriter()
		if err != nil {
			log.Printf("failed to build task output stream: %v\n", err)
			continue
		}
		cmd := task.ToExecCommand(out)
		task.Execute(cmd, quit)
	}
}

// PreBootTaskWorkers spins up all workers ready to consume tasks from queue.
func PreBootTaskWorkers(max int, wg *sync.WaitGroup, quit <-chan struct{}, tasksQueue <-chan *Task) {
	for i := 0; i < max; i++ {
		wg.Add(1)
		id := i
		go TaskWorker(id, wg, quit, tasksQueue)
	}
}
