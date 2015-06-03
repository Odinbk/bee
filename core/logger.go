package core

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type report struct {
	errorCount       int
	responseDuration time.Duration
	responseCount    int

	rps                 float64
	errorRate           float64
	fastestResponseTime float64
	slowestResponseTime float64
	averageResponseTime float64
}

func collectResults(results chan *result, allDone chan bool, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	r := new(report)
	start := time.Now()
	tick := time.Tick(time.Second)

	for {
		select {
		case result := <-results:
			if result.err != nil {
				r.errorCount++
			}
			r.responseDuration += result.duration
			r.responseCount++
		case <-tick:
			fmt.Printf(strings.Repeat(">", 1))
		case <-allDone:
			duration := int(time.Now().Sub(start) / time.Second)
			r.rps = float64(r.responseCount / duration)
			r.averageResponseTime = float64(int(r.responseDuration/time.Millisecond) / r.responseCount)
			r.errorRate = float64(r.errorCount / r.responseCount)

			fmt.Println()
			fmt.Println(fmt.Sprintf("RPS %.2f", r.rps))
			fmt.Println(fmt.Sprintf("Average response time %.2f ms", r.averageResponseTime))
			fmt.Println(fmt.Sprintf("ErrorRate %.2f%%", r.errorRate*100))

			return
		}
	}
}
