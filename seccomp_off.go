// +build !seccomp

package wololo

// Dummy function if no sandboxing is used
func Sandbox() error {
    return nil
}
