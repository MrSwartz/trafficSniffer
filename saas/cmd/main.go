package main

import (
	"fmt"
	"log"
	"runtime"
	packets "saas/internal/core/capture"
	"saas/internal/optimizer"
)

var (
	DEVICES = 1
	ALLOCATE = 0
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	optimizer.Allocator(ALLOCATE)  //see description inside package

	fmt.Printf("Allocated %d GB\n", ALLOCATE)
	fmt.Printf("Compiller  %s %s\nGOARCH     %s\nGOOS       %s\nGOROOT     %s\nNumCPU     %d\n\n", runtime.Compiler, runtime.Version(), runtime.GOARCH, runtime.GOOS, runtime.GOROOT(), runtime.NumCPU())
}

func main() {
	init := make([]packets.Interface, 1)
	tmp := packets.Interface{}

	dev := packets.ParseDevs(1)

	for _, v := range dev {
		if v.Name != "" {
			tmp.Device = v.Name
			tmp.Snaplen = 1500
		}
		init = append(init, tmp)
	}

	for i, v := range init {
		fmt.Println(i, v)
	}

	packets.Capture(init[4], "tcp")
}