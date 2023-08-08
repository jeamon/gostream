package main

import (
	"io"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/pelletier/go-toml"
)

// contains routines to execute set of tasks defined into a toml-encoded file like described in `example.toml`.

// DispatchTasksToml processes a toml file containing a set of tasks
// and send them to the worker input queue for execution.
func DispatchTasksToml(filename string, tasksQueue chan<- *Task) {
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
	err = toml.Unmarshal(data, &tasks)

	if err != nil {
		log.Fatalf("cannot parse toml: %v", err)
	}

	for i := 0; i < len(tasks.Tasks); i++ {
		tasksQueue <- &(tasks.Tasks[i])
		log.Printf("[loaded ] %s", tasks.Tasks[i].Task)
	}
}

// ProcessTasksToml orchestrates processing of a toml file
// containing all tasks with their details.
func ProcessTasksToml(filename string, quit <-chan struct{}) {
	wg := &sync.WaitGroup{}
	cores := runtime.NumCPU()
	// assuming at least 2 threads per core and
	// relaxing one for the main.
	maxWorkers := (cores * 2) - 1
	tasksQueue := make(chan *Task, maxWorkers)
	PreBootTaskWorkers(maxWorkers, wg, quit, tasksQueue)
	DispatchTasksToml(filename, tasksQueue)
	close(tasksQueue)
	wg.Wait()
}
