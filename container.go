/*
	Namespaces - controlling what the container can see and what they can't see. Basically restricting the container's
	view of the host process. Created with syscalls
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
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
	fmt.Printf("Running %v\n", os.Args[2:])

	

	cmd := exec.Command(os.Args[2], os.Args[3:]...)

	// Mapping the shell's stdin/stdout/stderr to the commands stdin/stdout/stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Namespace handling (sending in the sysprocattributes of the syscall and setting it equal to the command)
	cmd.SysProcAttr = &syscall.SysProcAttr {
		Cloneflags: syscall.CLONE_NEWUTS, // unix time sharing system namespace (only shares the hostname)
	}

	// Necessary to see the errors appearing in the code and have it be outputted to stdout
	must(cmd.Run())
}

func must(err error){
	if err != nil {
		panic(err)
	}
}
