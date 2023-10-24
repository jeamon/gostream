package main

import (
	"log"
	"os"
	"runtime"
	"sync"

	"gopkg.in/yaml.v3"
)

// contains routines to execute set of tasks defined into yaml-encoded files like described in `examples/example.yaml`.

// DispatchTasksYaml processes yaml files containing a set of tasks
// and send them to the worker input queue for execution.
func DispatchTasksYaml(filenames *[]string, tasksQueue chan<- *Task) {
	for _, filename := range *filenames {
		file, err := os.Open(filename)
		if err != nil {
			file.Close()
			log.Fatalf("cannot open file: %v", err)
		}

		var tasks Tasks

		decoder := yaml.NewDecoder(file)
		err = decoder.Decode(&tasks)
		if err != nil {
			file.Close()
			log.Fatalf("cannot parse yaml: %v", err)
		}

		for i := 0; i < len(tasks.Tasks); i++ {
			tasksQueue <- &(tasks.Tasks[i])
			log.Printf("[loaded ] %s", tasks.Tasks[i].Task)
		}
		file.Close()
	}
}

// ProcessTasksYaml orchestrates processing of yaml files
// containing all tasks with their details.
func ProcessTasksYaml(filenames *[]string, quit <-chan struct{}) {
	wg := &sync.WaitGroup{}
	cores := runtime.NumCPU()
	// assuming at least 2 threads per core and
	// relaxing one for the main.
	maxWorkers := (cores * 2) - 1
	tasksQueue := make(chan *Task, maxWorkers)
	PreBootTaskWorkers(maxWorkers, wg, quit, tasksQueue)
	DispatchTasksYaml(filenames, tasksQueue)
	close(tasksQueue)
	wg.Wait()
}
