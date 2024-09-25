package config

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"reflect"
	"runtime"
)

func msg_hook(ze zapcore.Entry) error {
	v := reflect.ValueOf(ze)
	level := v.FieldByName("Level").Interface().(zapcore.Level)
	switch level {
	case zapcore.DebugLevel:
		// Print stack trace
		stackBuf := make([]byte, 1024)
		stackSize := runtime.Stack(stackBuf, true)
		fmt.Printf("Stack trace:\n%s\n", stackBuf[:stackSize])

		// Print heap information
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		fmt.Printf("Heap Alloc = %v MiB", bToMb(memStats.HeapAlloc))
		fmt.Printf("\tHeap Sys = %v MiB", bToMb(memStats.HeapSys))
		fmt.Printf("\tHeap Idle = %v MiB", bToMb(memStats.HeapIdle))
		fmt.Printf("\tHeap Inuse = %v MiB", bToMb(memStats.HeapInuse))
		fmt.Printf("\tHeap Released = %v MiB", bToMb(memStats.HeapReleased))
		fmt.Printf("\tHeap Objects = %v\n", memStats.HeapObjects)
	default:
	}
	return nil
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
