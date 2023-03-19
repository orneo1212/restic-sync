package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func handle_info(cCtx *cli.Context) {
	dir_path := cCtx.Args().First()
	dirs := Scan(dir_path)
	dirs = LookForFiles(dirs)
	fmt.Println("Found", len(dirs), "directories to backup")
	for _, dirpath := range dirs {
		fmt.Println(dirpath)
	}
	os.Exit(0)
}

func handle_backup(cCtx *cli.Context) {
	repo_location := cCtx.Args().First()
	backup_path := cCtx.Args().Get(1)
	if repo_location == "" || backup_path == "" {
		println("Specify repository path and backup_location")
		return
	}

	dirs := Scan(backup_path)
	dirs = LookForFiles(dirs)
	fmt.Println("Found", len(dirs), "directories to backup")
	if len(dirs) > 0 && os.Getenv("RESTIC_PASSWORD") == "" {
		password := PasswordPrompt("Repository password:")
		os.Setenv("RESTIC_PASSWORD", password)
	}
	for index, dirpath := range dirs {
		exclude := dirs[index+1:]
		fmt.Println("Starting backup ", dirpath)
		Backup(repo_location, dirpath, exclude)
		fmt.Println()
	}
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:     "sync",
				Aliases:  []string{"s"},
				Usage:    "Make backup of all directories found",
				Category: "backup",
				Action: func(cCtx *cli.Context) error {
					handle_backup(cCtx)
					return nil
				},
			},
			{
				Name:     "info",
				Aliases:  []string{"i"},
				Usage:    "Show informations about directories",
				Category: "informations",
				Action: func(cCtx *cli.Context) error {
					handle_info(cCtx)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
