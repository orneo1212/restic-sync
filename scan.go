package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var excluded = []string{".git", ".hg", "node_modules", "cache", "src", "lib"}
var look_for_filenames = []string{".resticsync"}

func Scan(directory_path string) []string {
	directory_path, _ = filepath.Abs(directory_path)
	var dirs = make([]string, 0, 1)

	files, err := ioutil.ReadDir(directory_path)
	if err != nil {
		return []string{}
	}

	dirs = append(dirs, directory_path)

	for _, info := range files {
		path := filepath.Join(directory_path, info.Name())
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		}
		if !info.IsDir() {
			continue
		}
		//Exclude dirs
		exc := false
		for _, excluded := range excluded {
			if info.Name() == excluded {
				exc = true
				break
			}
		}
		if exc {
			continue
		}
		dirs = append(dirs, Scan(path)...)
	}
	return dirs
}

func LookForFiles(dirs []string) []string {
	result := make([]string, 0)
	for _, directory_path := range dirs {
		files, err := ioutil.ReadDir(directory_path)
		if err != nil {
			continue
		}
		for _, info := range files {
			if info.IsDir() {
				continue
			}
			for _, lookname := range look_for_filenames {
				if info.Name() == lookname {
					result = append(result, directory_path)
				}
			}
		}
	}
	return result
}
