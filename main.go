package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/KimMachineGun/automemlimit/memlimit"
)

// get the memory limit
func getMemoryLimit() uint64 {
	// assume we are running in a container with cgroup V2
	limit, err := memlimit.FromCgroupV2()

    if err != nil {
		limit = 0
        fmt.Printf("error reading memory limit %s", err.Error())
    } 
	return limit
}

// returns true if system memory is still available
func checkMemAvailable() (ok bool){
	ok = true
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	limit:= getMemoryLimit()

	// if we are using ~90% of memory, return false
	if float64(m.Alloc / limit) >= 0.9 ||  float64(m.Sys / limit) >= 0.9 {
		ok = false
		fmt.Println("Near memory limit")
	}
	fmt.Println("-------------------------------")
	fmt.Printf("Limit = %v MB\n",bToMb(limit))
	fmt.Printf("TotalAlloc = %v MB\n", bToMb(m.TotalAlloc))
	fmt.Printf("Alloc = %v MB", bToMb(m.Alloc))
	fmt.Printf("\tSys = %v MB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
	fmt.Println("-------------------------------")

	return ok
}

func main() {
	// disable GC
	debug.SetGCPercent(-1)
	fmt.Printf("Starting go-memtest.  Detected system memory = %v MB\n", bToMb(getMemoryLimit()))
	
	var data []int
	// loop forever to fill data
	for i := 0; ; i++ {
		if checkMemAvailable(){
			data = append(data, i)
		} else {
			for !checkMemAvailable(){
				fmt.Println("FEED ME MORE MEMORY")
				time.Sleep(1 * time.Second)
			}
		}
		if i%100000 == 0 {
			time.Sleep(1 * time.Second) // Pause for 1 second
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Println("-------------------------------")
			fmt.Printf("TotalAlloc = %v MB\n", bToMb(m.TotalAlloc))
			fmt.Printf("Alloc = %v MB", bToMb(m.Alloc))
			fmt.Printf("\tSys = %v MB", bToMb(m.Sys))
			fmt.Printf("\tNumGC = %v\n", m.NumGC)
			fmt.Println("-------------------------------")
		}
	}
}


func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
