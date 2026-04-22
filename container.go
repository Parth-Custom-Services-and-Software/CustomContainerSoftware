/*
	Namespaces - controlling what the container can see and what they can't see. Basically restricting the container's
	view of the host process. Created with syscalls
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

	// Running the control security groups
	cg()
	
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
		// NS = namespace namespace (first namespace so thats why its called that. Its primary purpose is for mounts)
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		// Currently the root directory on the host recursively shares a property across all namespaces and mounts, 
		// so you need to deliberately turn that off:
		// Now that we have CLONE_NEWNS in both the unshareflags and the cloneflags,
		// running mount | grep /proc
		// on the host machien doesn't show the mounted proc through our container :D
		Unshareflags: syscall.CLONE_NEWNS,
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
	defer syscall.Unmount("proc", 0) // defer so that it'll run whenver this function returns or panics

	// Mount is also necessary for the sys folder
	syscall.Mount("sysfs", "sys", "sysfs", 0, "")
	defer syscall.Unmount("sys", 0)

	// Same thing with cgroup2
	syscall.Mount("cgroup2", "/sys/fs/cgroup", "cgroup2", 0, "")
	defer syscall.Unmount("/sys/fs/cgroup", 0)
	
	cmd := exec.Command(os.Args[2], os.Args[3:]...)

	// Mapping the shell's stdin/stdout/stderr to the commands stdin/stdout/stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	// Don't need to create another namespace here

	// Necessary to see the errors appearing in the code and have it be outputted to stdout
	must(cmd.Run())
}

func cg() {
	MAX_PROCESES := "20"

	// These are pretty self explanatory but
	// All its doing is creating a path for the control groups:
	// /sys/fs/cgroup/container
	// and then creating a directory with that path and giving it the 
	// access values: 0755
	// which is just drwxr-xr-x (rwx for owner, rx for group, and rx for other)
	cgroups := "/sys/fs/cgroup/container"
	err := os.Mkdir(cgroups, 0755)

	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	// Inside the control group, there can only be 20 processes (or whatever I define with MAX_PROCESSES)
	must(os.WriteFile(filepath.Join(cgroups, "pids.max"), []byte(MAX_PROCESES), 0700))
	// Adds the current process (the bash shell) to the control group procs file, making it subject to the same
	// specifications given by the control group
	must(os.WriteFile(filepath.Join(cgroups, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must(err error){
	if err != nil {
		panic(err)
	}
}
