package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"runtime"
	"sync"
)

// contains routines to execute set of tasks containing into files like described in `example.file`.
// the pattern consists of having each task presented as a json object at each line into each file.

// DispatchTasksFile processes files containing a set of tasks
// and send them to the worker input queue for execution.
func DispatchTasksFile(filenames *[]string, tasksQueue chan<- *Task) {
	for _, filename := range *filenames {
		file, err := os.Open(filename)
		if err != nil {
			file.Close()
			log.Fatalf("cannot open file: %v", err)
		}

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
		file.Close()
		if err != nil {
			log.Fatalf("failed reading file content: %v", err)
		}
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

// ProcessTasksFile orchestrates processing of input files
// containing at each line a json formatted task details.
func ProcessTasksFile(filenames *[]string, quit <-chan struct{}) {
	wg := &sync.WaitGroup{}
	cores := runtime.NumCPU()
	// assuming at least 2 threads per core and
	// relaxing one for the main.
	maxWorkers := (cores * 2) - 1
	tasksQueue := make(chan *Task, maxWorkers)
	PreBootTaskWorkers(maxWorkers, wg, quit, tasksQueue)
	DispatchTasksFile(filenames, tasksQueue)
	close(tasksQueue)
	wg.Wait()
}
