package main

import (
	"fmt"
	"os"

	"github.com/horacehylee/go-youtube-dl/pkg/youtube"
)

func main() {
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

	fmt.Printf("Downloaded video at %v\n", f.Name())
}

func checkError(append string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, append+": %v\n", err)
		os.Exit(1)
	}
}
