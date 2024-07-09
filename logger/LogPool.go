package logger

import (
	"log"
	"sync"
	"time"
)

// create thread pool for recording log
type LogTask struct {
	message string
	done    chan bool
}
type LogPool struct {
	tasks chan LogEntry
	wg    sync.WaitGroup
}

func NewLogPool(workCount int) *LogPool {
	pool := &LogPool{
		tasks: make(chan LogEntry),
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
				log.Println(task.String())
				// Assuming a done channel is no longer needed as logging is synchronous
			}
		}()
	}
}

func (pool *LogPool) Log(level LogLevel, message string) {
	pool.tasks <- LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}
}

func (pool *LogPool) Stop() {
	close(pool.tasks)
	pool.wg.Wait()
}
