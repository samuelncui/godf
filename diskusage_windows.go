//go:build windows
// +build windows

package godf

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

type DiskUsage struct {
	freeBytes  int64
	totalBytes int64
	availBytes int64
}

var (
	proc    *syscall.Proc
	procErr error
	once    sync.Once
)

func initProc() {
	h, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		procErr = fmt.Errorf("load dll kernel32.dll fail, %w", err)
		return
	}

	p, err := h.FindProc("GetDiskFreeSpaceExW")
	if err != nil {
		procErr = fmt.Errorf("load proc GetDiskFreeSpaceExW from kernel32.dll fail, %w", err)
		return
	}

	proc = p
}

// NewDiskUsages returns an object holding the disk usage of volumePath
// or nil in case of error (invalid path, etc)
func NewDiskUsage(volumePath string) (*DiskUsage, error) {
	once.Do(initProc)
	if procErr != nil {
		return nil, procErr
	}

	pathPtr, err := syscall.UTF16PtrFromString(volumePath)
	if err != nil {
		return nil, fmt.Errorf("convert string to utf16 failed, %w", err)
	}

	du := new(DiskUsage)
	if _, _, err := proc.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(&du.freeBytes)),
		uintptr(unsafe.Pointer(&du.totalBytes)),
		uintptr(unsafe.Pointer(&du.availBytes)),
	); err != nil {
		return nil, fmt.Errorf("call dll method fail, %w", err)
	}

	return du, nil
}

// Free returns total free bytes on file system
func (du *DiskUsage) Free() uint64 {
	return uint64(du.freeBytes)
}

// Available returns total available bytes on file system to an unprivileged user
func (du *DiskUsage) Available() uint64 {
	return uint64(du.availBytes)
}

// Size returns total size of the file system
func (du *DiskUsage) Size() uint64 {
	return uint64(du.totalBytes)
}

// Used returns total bytes used in file system
func (du *DiskUsage) Used() uint64 {
	return du.Size() - du.Free()
}

// Usage returns percentage of use on the file system
func (du *DiskUsage) Usage() float32 {
	return float32(du.Used()) / float32(du.Size())
}
