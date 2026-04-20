/*
	Namespaces - controlling what the container can see and what they can't see. Basically restricting the container's
	view of the host process. Created with syscalls
*/

package main

import (
	"os"
	"fmt"
)

// This is the command we're trying to run (in docker) and this is the command we're going to parse:
// docker              run image <cmd> <params>
// go run container.go run       <cmd> <params>
// the image is not going to get added 

func main() {
	switch os.Args[1] {
	case "run":
		run()
	
	default:
		panic("Unrecognized command, aborting.")
	}
}

func run() {
	fmt.Printf("Running")
}

func must(err error){
	if err != nil {
		panic(err)
	}
}
