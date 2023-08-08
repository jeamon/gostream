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

// DispatchTasksFile processes a file containing a set of tasks
// and send them to the worker input queue for execution.
func DispatchTasksFile(filename string, tasksQueue chan<- *Task) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		taskStr := scanner.Text()
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

// ProcessTasksFile orchestrates processing of input file
// containing at each line a json formatted task details.
func ProcessTasksFile(filename string, quit <-chan struct{}) {
	wg := &sync.WaitGroup{}
	cores := runtime.NumCPU()
	// assuming at least 2 threads per core and
	// relaxing one for the main.
	maxWorkers := (cores * 2) - 1
	tasksQueue := make(chan *Task, maxWorkers)
	PreBootTaskWorkers(maxWorkers, wg, quit, tasksQueue)
	DispatchTasksFile(filename, tasksQueue)
	close(tasksQueue)
	wg.Wait()
}
