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
