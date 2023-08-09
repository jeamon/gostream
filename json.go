package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

// contains routines to execute set of tasks defined into json-encoded files like described in `example.json`.

// DispatchTasksJson processes json files containing a set of tasks
// and send them to the worker input queue for execution.
func DispatchTasksJson(filenames *[]string, tasksQueue chan<- *Task) {
	for _, filename := range *filenames {
		file, err := os.Open(filename)
		if err != nil {
			file.Close()
			log.Fatalf("cannot open file: %v", err)
		}

		data, err := io.ReadAll(file)
		if err != nil {
			file.Close()
			log.Fatalf("cannot read file: %v", err)
		}
		var tasks Tasks
		err = json.Unmarshal(data, &tasks)

		if err != nil {
			file.Close()
			log.Fatalf("cannot parse json: %v", err)
		}

		for i := 0; i < len(tasks.Tasks); i++ {
			tasksQueue <- &(tasks.Tasks[i])
			log.Printf("[loaded ] %s", tasks.Tasks[i].Task)
		}
		file.Close()
	}
}

// ProcessTasksJson orchestrates processing of json files
// containing all tasks with their details.
func ProcessTasksJson(filenames *[]string, quit <-chan struct{}) {
	wg := &sync.WaitGroup{}
	cores := runtime.NumCPU()
	// assuming at least 2 threads per core and
	// relaxing one for the main.
	maxWorkers := (cores * 2) - 1
	tasksQueue := make(chan *Task, maxWorkers)
	PreBootTaskWorkers(maxWorkers, wg, quit, tasksQueue)
	DispatchTasksJson(filenames, tasksQueue)
	close(tasksQueue)
	wg.Wait()
}
