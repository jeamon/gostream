package main

import (
	"io"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/pelletier/go-toml"
)

// contains routines to execute set of tasks defined into toml-encoded files like described in `example.toml`.

// DispatchTasksToml processes toml files containing a set of tasks
// and send them to the worker input queue for execution.
func DispatchTasksToml(filenames *[]string, tasksQueue chan<- *Task) {
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
		err = toml.Unmarshal(data, &tasks)

		if err != nil {
			file.Close()
			log.Fatalf("cannot parse toml: %v", err)
		}

		for i := 0; i < len(tasks.Tasks); i++ {
			tasksQueue <- &(tasks.Tasks[i])
			log.Printf("[loaded ] %s", tasks.Tasks[i].Task)
		}
		file.Close()
	}
}

// ProcessTasksToml orchestrates processing of toml files
// containing all tasks with their details.
func ProcessTasksToml(filenames *[]string, quit <-chan struct{}) {
	wg := &sync.WaitGroup{}
	cores := runtime.NumCPU()
	// assuming at least 2 threads per core and
	// relaxing one for the main.
	maxWorkers := (cores * 2) - 1
	tasksQueue := make(chan *Task, maxWorkers)
	PreBootTaskWorkers(maxWorkers, wg, quit, tasksQueue)
	DispatchTasksToml(filenames, tasksQueue)
	close(tasksQueue)
	wg.Wait()
}
