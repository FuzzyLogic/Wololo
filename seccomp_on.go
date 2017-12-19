// +build seccomp

package wololo

import(
    "syscall"
    "github.com/seccomp/libseccomp-golang"
)

var syscalls []string

// Enable sandboxing for this application by registering only the required
// syscalls with seccomp.
func Sandbox() error {
    // Steps for creating a filter:
    // - Create new filter
    // - Get and add syscalls to filter
    // - Load filter

    // List of allowed system calls
    syscalls = []string {
        "openat",
        "read",
        "mmap",
        "close",
        "getpid",
        "write",
        "socket",
        "setsockopt",
        "bind",
        "listen",
        "epoll_ctl",
        "getsockname",
        "accept4",
        "epoll_wait",
        "futex",
        // Not sure about the following
        "connect",
        "epoll_create1",
        "sendto",
        "recvfrom",
        "getpeername",
    }

    // Create the filter - the default action for unmatched calls is to kill the process
    filter, err := seccomp.NewFilter(seccomp.ActErrno.SetReturnCode(int16(syscall.EPERM)))
    //filter, err := seccomp.NewFilter(seccomp.ActAllow)
    if err != nil {
        return err
    }
    defer filter.Release()

    // Add rules for the allowed syscalls
    for _, scName := range syscalls {
        call, err := seccomp.GetSyscallFromName(scName)
        if err != nil {
            return err
        }

        //err = filter.AddRule(call, seccomp.ActErrno.SetReturnCode(0x2))
        err = filter.AddRule(call, seccomp.ActAllow)
        if err != nil {
            return err
        }
    }

    // Try to load the syscall filter
    err = filter.Load()
    if err != nil {
        return err
    }

    return nil
}
