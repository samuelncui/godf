//go:build !windows
// +build !windows

package godf

import "syscall"

// DiskUsage contains usage data and provides user-friendly access methods
type DiskUsage struct {
	stat *syscall.Statfs_t
}

// NewDiskUsages returns an object holding the disk usage of volumePath
// or nil in case of error (invalid path, etc)
func NewDiskUsage(volumePath string) (*DiskUsage, error) {
	stat := new(syscall.Statfs_t)
	if err := syscall.Statfs(volumePath, stat); err != nil {
		return nil, err
	}

	return &DiskUsage{stat}, nil
}

// Free returns total free bytes on file system
func (du *DiskUsage) Free() int64 {
	return int64(du.stat.Bfree) * int64(du.stat.Bsize)
}

// Available return total available bytes on file system to an unprivileged user
func (du *DiskUsage) Available() int64 {
	return int64(du.stat.Bavail) * int64(du.stat.Bsize)
}

// Size returns total size of the file system
func (du *DiskUsage) Size() int64 {
	return int64(du.stat.Blocks) * int64(du.stat.Bsize)
}

// Used returns total bytes used in file system
func (du *DiskUsage) Used() int64 {
	return du.Size() - du.Free()
}

// Usage returns percentage of use on the file system
func (du *DiskUsage) Usage() float64 {
	return float64(du.Used()) / float64(du.Size())
}
