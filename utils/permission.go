package utils

import "syscall"

// isRoot checks if the process is run with root permission
func IsRoot() bool {
	return syscall.Getuid() == 0 && syscall.Geteuid() == 0
}
