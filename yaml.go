package main

import (
	"log"
	"os"
	"runtime"
	"sync"

	"gopkg.in/yaml.v3"
)

// contains routines to execute set of tasks defined into a yaml-encoded file like described in `example.yaml`.

// DispatchTasksYaml processes a yaml file containing a set of tasks
// and send them to the worker input queue for execution.
func DispatchTasksYaml(filename string, tasksQueue chan<- *Task) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open file: %v", err)
	}
	defer file.Close()

	var tasks Tasks

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		log.Fatalf("cannot parse yaml: %v", err)
	}

	for i := 0; i < len(tasks.Tasks); i++ {
		tasksQueue <- &(tasks.Tasks[i])
		log.Printf("[loaded ] %s", tasks.Tasks[i].Task)
	}
}

// ProcessTasksYaml orchestrates processing of a yaml file
// containing all tasks with their details.
func ProcessTasksYaml(filename string, quit <-chan struct{}) {
	wg := &sync.WaitGroup{}
	cores := runtime.NumCPU()
	// assuming at least 2 threads per core and
	// relaxing one for the main.
	maxWorkers := (cores * 2) - 1
	tasksQueue := make(chan *Task, maxWorkers)
	PreBootTaskWorkers(maxWorkers, wg, quit, tasksQueue)
	DispatchTasksYaml(filename, tasksQueue)
	close(tasksQueue)
	wg.Wait()
}
