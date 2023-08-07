package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"runtime"
	"sync"
)

// contains routines to execute set of tasks containing into a file like described in `example.file`.
// the pattern consists of having each task infos presented as a json object at each line into the file.

// DispatchTasks processes a file containing a set of tasks
// and send them to the worker input queue for execution.
func DispatchTasks(filename string, tasksQueue chan<- *Task) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var taskStr string

	for scanner.Scan() {
		taskStr = scanner.Text()
		taskPtr, err := ParseTaskJson(taskStr)
		if err != nil {
			log.Printf("[failed ] %s", taskStr)
		} else {
			log.Printf("[loaded ] %s", taskStr)
		}

		tasksQueue <- taskPtr
	}

	err = scanner.Err()
	if err != nil {
		log.Fatalf("failed reading file content: %v", err)
	}
}

// ParseTaskJson converts a json string to a Task object.
func ParseTaskJson(data string) (*Task, error) {
	task := Task{}
	if err := json.Unmarshal([]byte(data), &task); err != nil {
		return &task, err
	}
	return &task, nil
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

func PreBootTaskWorkers(max int, wg *sync.WaitGroup, quit <-chan struct{}, tasksQueue <-chan *Task) {
	for i := 0; i < max; i++ {
		wg.Add(1)
		id := i
		go TaskWorker(id, wg, quit, tasksQueue)
	}
}

func ProcessTasksFile(filename string, quit <-chan struct{}) {
	wg := &sync.WaitGroup{}
	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)
	// assuming 2 threads per core and relaxing one for the main.
	maxWorkers := (cores * 2) - 1
	tasksQueue := make(chan *Task, maxWorkers)
	PreBootTaskWorkers(maxWorkers, wg, quit, tasksQueue)
	DispatchTasks(filename, tasksQueue)
	close(tasksQueue)
	wg.Wait()
}
