package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/horacehylee/go-youtube-dl/pkg/youtube"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	// f, err := os.Create("cpu.pprof")
	// log.Printf("cpuprofile: %v\n", f.Name())
	// if err != nil {
	// 	log.Fatal("could not create CPU profile: ", err)
	// }
	// defer f.Close() // error handling omitted for example
	// if err := pprof.StartCPUProfile(f); err != nil {
	// 	log.Fatal("could not start CPU profile: ", err)
	// }
	// defer pprof.StopCPUProfile()

	if len(os.Args) < 2 {
		checkError("Argument is required", fmt.Errorf("no argument is passed"))
	}

	vid := os.Args[1]

	var c youtube.Client

	f, err := os.Create(fmt.Sprintf("./%v.m4a", vid))
	checkError("Failed to create file", err)
	defer f.Close()

	err = c.Download(f, vid)
	checkError("Failed to download video", err)

	fmt.Fprintf(os.Stdout, "Downloaded video at %v\n", f.Name())
}

func checkError(append string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, append+": %v\n", err)
		os.Exit(1)
	}
}
