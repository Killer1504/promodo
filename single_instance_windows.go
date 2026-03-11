//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

const (
	// Unique mutex name — must match across all instances.
	mutexName = "PomodoroFocusTimer-{F3A2B1C0-9D8E-4F67-A312-14DemoWails}"

	errAlreadyExists uintptr = 183 // ERROR_ALREADY_EXISTS

	swRestore = 9 // ShowWindow: restore from minimized
	swShow    = 5 // ShowWindow: show normal
)

var (
	kernel32        = syscall.NewLazyDLL("kernel32.dll")
	user32          = syscall.NewLazyDLL("user32.dll")
	procCreateMutex = kernel32.NewProc("CreateMutexW")
	procFindWindow  = user32.NewProc("FindWindowW")
	procShowWindow  = user32.NewProc("ShowWindow")
	procSetForeground = user32.NewProc("SetForegroundWindow")
	procIsIconic    = user32.NewProc("IsIconic")
)

// ensureSingleInstance returns true if another instance is already running.
// If so, it brings that instance's window to the foreground and the caller should exit.
// If this is the first instance, it holds the named mutex until the process exits.
func ensureSingleInstance() bool {
	namePtr, _ := syscall.UTF16PtrFromString(mutexName)

	handle, _, err := procCreateMutex.Call(
		0,                          // lpMutexAttributes (nil)
		0,                          // bInitialOwner = false
		uintptr(unsafe.Pointer(namePtr)),
	)

	if handle == 0 {
		// CreateMutex failed entirely — allow running.
		return false
	}

	if err == syscall.Errno(errAlreadyExists) {
		// Another instance owns the mutex → bring it to front.
		bringExistingToFront()
		return true
	}

	// We now own the mutex; it is released automatically when the process exits.
	// Keep the handle in a global so the GC doesn't close it prematurely.
	_ = handle
	return false
}

// bringExistingToFront finds the existing app window by its title and raises it.
func bringExistingToFront() {
	titlePtr, _ := syscall.UTF16PtrFromString("Pomodoro Focus Timer")

	hwnd, _, _ := procFindWindow.Call(
		0, // lpClassName — any class
		uintptr(unsafe.Pointer(titlePtr)),
	)
	if hwnd == 0 {
		return
	}

	// If minimized, restore it; otherwise just show.
	iconic, _, _ := procIsIconic.Call(hwnd)
	if iconic != 0 {
		procShowWindow.Call(hwnd, swRestore)
	} else {
		procShowWindow.Call(hwnd, swShow)
	}

	procSetForeground.Call(hwnd)
}
