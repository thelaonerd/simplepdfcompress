package worker

import (
	"sync"

	"simplepdfcompress/internal/compression"
)

// Job represents a single compression task
type Job struct {
	InputPath  string
	OutputPath string
	Options    compression.CompressionOptions
}

// Result represents the outcome of a compression job
type Result struct {
	Job          Job
	OriginalSize int64
	FinalSize    int64
	Error        error
}

// RunPool processes a list of jobs using a specified number of concurrent workers
// It returns a channel that streams results as they complete.
func RunPool(jobs []Job, numWorkers int) <-chan Result {
	jobChan := make(chan Job, len(jobs))
	resultChan := make(chan Result, len(jobs))
	var wg sync.WaitGroup

	// Fill job channel
	for _, job := range jobs {
		jobChan <- job
	}
	close(jobChan)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobChan {
				initial, final, err := compression.CompressPDF(job.InputPath, job.OutputPath, job.Options)
				resultChan <- Result{
					Job:          job,
					OriginalSize: initial,
					FinalSize:    final,
					Error:        err,
				}
			}
		}()
	}

	// Wait and close results in a separate goroutine
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}
