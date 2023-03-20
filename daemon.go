package main

import (
	"fmt"
	"time"

	"github.com/alitto/pond"
)

func StartDaemon(scan_dirs []string, repo_location string) {
	interval_minutes := 2 // minutes
	pool := pond.New(100, 1000, pond.MinWorkers(2))

	fmt.Println("Starting backup daemon. Refresh time", interval_minutes, "minutes")
	for {
		time.Sleep(time.Minute * time.Duration(interval_minutes))
		// Scan for directories to backup
		var dirs []string
		for i := 0; i < len(scan_dirs); i++ {
			ldirs := Scan(scan_dirs[i])
			dirs = append(dirs, LookForFiles(ldirs)...)
		}
		// Backup
		for i := 0; i < len(dirs); i++ {
			path := dirs[i]
			pool.Submit(func() {
				fmt.Println("Start backing up ", path)
				err := Backup(repo_location, path, []string{})
				if err == nil {
					fmt.Println("Finished backup of ", path)
					fmt.Println()
				} else {
					fmt.Println("Error backup of ", path, ":\n", err)
				}
			})
		}
		pool.Group().Wait()
	}
}
