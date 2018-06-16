package main

import (
	"os"
	"runtime"
)

func userHome() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}
