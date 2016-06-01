package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/procfs"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Getting proc PID and create process object.
var pid = os.Getpid()
var proc, _ = procfs.NewProc(pid)

var iterationCount = flag.Int("iterationCount", 100, "Count of iterations in goroutine")
var goroutineCount = flag.Int("goroutineCount", 1000, "Count of goroutines to be spawned")
var maxThreadCount = 0
var lock sync.Mutex

func logStats() {
	stats, _ := proc.NewStat()
	log.Println("Num threads:" + strconv.Itoa(stats.NumThreads))
	lock.Lock()
	defer lock.Unlock()
	if maxThreadCount < stats.NumThreads {
		maxThreadCount = stats.NumThreads
	}
}

func worker(wg *sync.WaitGroup) {

	for i := 0; i < *iterationCount; i++ {
		logStats()
		time.Sleep(time.Second * 1)
	}
	wg.Done()
}

func RunTest() {
	var wg sync.WaitGroup
	wg.Add(*goroutineCount)

	for i := 0; i < *goroutineCount; i++ {
		go worker(&wg)
	}

	wg.Wait()
}

func main() {
	log.Println("Starting benchark")
	log.Println("Numbers of CPU " + strconv.Itoa(runtime.NumCPU()))
	// In Go version 1.5 or higher this string is redundant. Added for backward compatibility
	// and clarity.
	log.Println("GOMAXPROC:" + strconv.Itoa(runtime.GOMAXPROCS(runtime.NumCPU())))
	log.Println("Process PID: " + strconv.Itoa(pid))
	time.Sleep(time.Second * 10)
	fmt.Scan()
	logStats()
	RunTest()
	log.Println("Test finished, max thread count " + strconv.Itoa(maxThreadCount))
}
