package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/gosimple/slug"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/term"
)

type RepoConfig struct {
	Id       string
	Name     string
	Category string
}

func create_excluded_params(excluded []string) []string {
	var list2 []string
	for _, x := range excluded {
		list2 = append(list2, "--exclude", x)
	}
	return list2
}

func read_config(directory_path string) RepoConfig {
	var cfg RepoConfig
	var config_path string = filepath.Join(directory_path, ".resticsync")
	var modified bool = false

	dat, err := os.ReadFile(config_path)
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal([]byte(dat), &cfg)
	if err != nil {
		panic(err)
	}
	if cfg.Name == "" {
		modified = true
		cfg.Name = "Unnamed"
	}
	if cfg.Id == "" {
		modified = true
		cfg.Id = generate(20)
	}

	if modified {
		dat, err := toml.Marshal(&cfg)
		if err != nil {
			panic(err)
		}
		os.WriteFile(config_path, dat, 0644)
	}
	return cfg
}

func generate(n int) string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
	str := make([]rune, n)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}

func PasswordPrompt(label string) string {
	var s string
	for {
		fmt.Fprint(os.Stderr, label+" ")
		b, _ := term.ReadPassword(int(syscall.Stdin))
		s = string(b)
		if s != "" {
			break
		}
	}
	fmt.Println()
	return s
}

func Backup(repository_path string, backup_location string, excluded []string) {
	config := read_config(backup_location)
	// Create tag
	tags := []string{"--tag=" + config.Id, "--tag=" + slug.Make(config.Name)}
	if config.Category != "" {
		tags = append(tags, "--tag="+slug.Make(config.Category))
	}

	// Create restric arguments list
	cmd := []string{"-r", repository_path, "--one-file-system", "--host=resticsync"}
	cmd = append(cmd, create_excluded_params(excluded)...)
	cmd = append(cmd, tags...)
	cmd = append(cmd, "backup", ".")

	c := exec.Command("restic", cmd...)
	c.Dir = backup_location
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Start()
	c.Wait()
}

func main() {
	args := os.Args[1:]
	if len(args) >= 2 {
		dirs := Scan(args[1])
		dirs = LookForFiles(dirs)
		fmt.Println("Found", len(dirs), "directories to backup")
		if len(dirs) > 0 && os.Getenv("RESTIC_PASSWORD") == "" {
			password := PasswordPrompt("Repository password:")
			os.Setenv("RESTIC_PASSWORD", password)
		}
		for index, dirpath := range dirs {
			exclude := dirs[index+1:]
			fmt.Println("Starting backup ", dirpath)
			Backup(args[0], dirpath, exclude)
			fmt.Println()
		}
	} else {
		fmt.Println("Specify repository path and backup directory")
	}
}
