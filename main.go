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

	// When RESTIC_REPOSITORY is isset ignore first argument
	if backup_path == "" && os.Getenv("RESTIC_REPOSITORY") != "" {
		repo_location = os.Getenv("RESTIC_REPOSITORY")
		backup_path = cCtx.Args().Get(0)
	}

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
	for _, dirpath := range dirs {
		Backup(repo_location, dirpath, []string{})
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
