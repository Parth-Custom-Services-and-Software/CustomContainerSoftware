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
	case "child":
		child()
	default:
		panic("Unrecognized command, aborting.")
	}
}

func run() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	
	// Once the shell is ran, it should reinvoke the same process but inside the new namespace (which is the child function)
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	// Mapping the shell's stdin/stdout/stderr to the commands stdin/stdout/stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Namespace handling (sending in the sysprocattributes of the syscall and setting it equal to the command)
	cmd.SysProcAttr = &syscall.SysProcAttr {
		// UTS = unix time sharing system namespace (doesn't share the hostname)
		// PID = process id namespace (doesn't share pids)
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID, // unix time sharing system namespace (only shares the hostname)
	}

	// Necessary to see the errors appearing in the code and have it be outputted to stdout
	must(cmd.Run())
}

func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	// but we do want to set a new hostname
	syscall.Sethostname([]byte("container"))

	fsPath := os.Getenv("CHROOT")
	if fsPath == "" {
		panic("CHROOT env variable cannot be found :(")
	}
	// Note that we can't use "~" here since the command is not being interpreted
	// Its getting directly run, which means that ~ is not translated to /home/USERNAME
	// So we set the root to a cloned ubuntu filesystem (so that we can have a separate /proc folder to isolate the output of the ps)
	// Then we set the current working directory of the function to the root (which is the cloned ubuntu filesystem)
	syscall.Chroot(fsPath)
	syscall.Chdir("/")
	
	// Mount is necessary for the kernel to recognize the new /proc folder as a proxy
	syscall.Mount("proc", "proc", "proc", 0, "")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)

	// Mapping the shell's stdin/stdout/stderr to the commands stdin/stdout/stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	// Don't need to create another namespace here

	// Necessary to see the errors appearing in the code and have it be outputted to stdout
	must(cmd.Run())

	// Unmount the proc folder to cleanup
	syscall.Unmount("/proc", 0)
}
func must(err error){
	if err != nil {
		panic(err)
	}
}
