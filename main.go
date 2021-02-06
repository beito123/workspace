package main

/*
	workspace

	Copyright (c) 2021 bei2

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// WorkspaceFormat is a format for workspace name
const WorkspaceFormat = "2006-01-02"

// Deadline is a deadline(day) for workspace
// You can remove expired workspaces by command
// If It's zero, remove all workspaces
// If It's -1, don't remove
const Deadline = 7

func main() {
	targetPath := flag.String("path", "./", "Set a path to make or remove workspace")
	removeMode := flag.Bool("remove", false, "Remove expired workspaces")
	flag.Parse()

	path, err := filepath.Abs(*targetPath)
	if err != nil {
		panic(err)
	}

	if *removeMode {
		removed, err := remove(path)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Removed %d folders", removed)
		return
	}

	name, err := make(path)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Made %s folder", name)
}

func make(path string) (string, error) {
	name := time.Now().Format(WorkspaceFormat)

	err := os.MkdirAll(filepath.Join(path, name), os.ModePerm)
	if err != nil {
		return "", err
	}

	return name, nil
}

func remove(path string) (int, error) {
	if Deadline == -1 {
		return 0, nil
	}

	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return 0, err
	}

	nowTime := time.Now()
	var count int

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		t, err := time.Parse(WorkspaceFormat, dir.Name())
		if err != nil {
			continue
		}

		if Deadline > 0 {
			if nowTime.Sub(t) < Deadline*(time.Hour*24) {
				continue
			}
		}

		dirPath := filepath.Join(path, dir.Name())
		err = os.RemoveAll(dirPath)
		if err != nil {
			return 0, err
		}

		count++
	}

	return count, nil
}
