package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

// DispatchTasksJson processes a json file containing a set of tasks
// and send them to the worker input queue for execution.
func DispatchTasksJson(filename string, tasksQueue chan<- *Task) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read file: %v", err)
	}
	var tasks Tasks
	err = json.Unmarshal(data, &tasks)

	if err != nil {
		log.Fatalf("cannot parse json: %v", err)
	}

	for i := 0; i < len(tasks.Tasks); i++ {
		tasksQueue <- &(tasks.Tasks[i])
		log.Printf("[loaded ] %s", tasks.Tasks[i].Task)
	}
}

// ProcessTasksJson orchestrates processing of a json file
// containing all tasks with their details.
func ProcessTasksJson(filename string, quit <-chan struct{}) {
	wg := &sync.WaitGroup{}
	cores := runtime.NumCPU()
	// assuming at least 2 threads per core and
	// relaxing one for the main.
	maxWorkers := (cores * 2) - 1
	tasksQueue := make(chan *Task, maxWorkers)
	PreBootTaskWorkers(maxWorkers, wg, quit, tasksQueue)
	DispatchTasksJson(filename, tasksQueue)
	close(tasksQueue)
	wg.Wait()
}
