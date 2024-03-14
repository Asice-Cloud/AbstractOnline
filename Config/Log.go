package Config

import (
	"log"
	"sync"
)

//create thread pool for recording log

type LogTask struct {
	message string
	done    chan bool
}

type LogPool struct {
	tasks chan LogTask
	wg    sync.WaitGroup
}

// Implement the methods of LogPool
func NewLogPool(workCount int) *LogPool {
	pool := &LogPool{
		tasks: make(chan LogTask),
	}

	pool.Start(workCount)

	return pool
}

func (pool *LogPool) Start(workCount int) {
	for i := 0; i < workCount; i++ {
		pool.wg.Add(1)

		go func() {
			defer pool.wg.Done()
			for task := range pool.tasks {
				log.Println(task.message)
				task.done <- true
			}
		}()
	}
}

func (pool *LogPool) Log(message string) {
	done := make(chan bool)

	pool.tasks <- LogTask{
		message: message,
		done:    done,
	}

	<-done
}

func (pool *LogPool) Stop() {
	close(pool.tasks)
	pool.wg.Wait()
}
