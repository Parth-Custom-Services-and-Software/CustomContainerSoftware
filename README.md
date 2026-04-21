# CustomContainerSoftware
Docker-like containerization software coded in Go. 
# Initial commit

# Notes:
In order to run container.go, you need to have root access. So if you do not have that, please use the command `sudo` before
the rest of the command. You should be running something like this:
`sudo go run container.go run /bin/bash`

# Child
We would like to change the hostname from the getgo, to differentiate it from our regular root. Because currently, even if you
do change the hostname, it'll still say root@ubuntu. 
Note that whenever writing out the command, you'll want to include the CHROOT environment variable like so:
`sudo CHROOT="/home/`[YOUR_NAME]`/`[NAME_OF_CLONED_UBUNTU_SYSTEM]` go run container.go run /bin/bash`

## Cloning an Ubuntu File System
Navigate to your desired directory (or create one)
`mkdir ubuntu-fs` (to make a new directory)
`cd ubuntu-fs`
Then use docker to create an optimized ubuntu fs
`docker extract $(docker create ubuntu) ubuntu.tar`
Extract that tar file:
`tar -xf ubuntu.tar`
Remove the tar file:
`rm ubuntu.tar`
Now your cloned ubuntu system should be listed in ubuntu-fs


Note that /proc is a pseudo-filesystem, which means that a space for the kernel and the userspace to share information
Given that the /proc in the ubuntu filesystem copy doesn't have anything in it, it needs to be mounted as a proxy so that
the kernel knows it needs to populate that directory with all the information about the running processes. 
Trying to run ps before mounting it results in this error: `Error, do this: mount -t proc proc /proc`
