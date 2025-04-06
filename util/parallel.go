package util

import (
	"runtime"
	"sync"
)

func ParallelFilter[T any](items []T, matchFn func(T) bool) []T {
	cores := runtime.NumCPU()
	size := (len(items) + cores - 1) / cores
	var wg sync.WaitGroup
	outCh := make(chan []T, cores)

	for i := 0; i < len(items); i += size {
		end := i + size
		if end > len(items) {
			end = len(items)
		}
		chunk := items[i:end]
		wg.Add(1)
		go func(chunk []T) {
			defer wg.Done()
			var filtered []T
			for _, item := range chunk {
				if matchFn(item) {
					filtered = append(filtered, item)
				}
			}
			outCh <- filtered
		}(chunk)
	}

	wg.Wait()
	close(outCh)

	var final []T
	for chunk := range outCh {
		final = append(final, chunk...)
	}
	return final
}
