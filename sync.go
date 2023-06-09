package main

import (
	"fmt"
	"io/ioutil"
	"log"
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

func GetRepositoryPassword() {
	if os.Getenv("RESTIC_PASSWORD") == "" {
		password := PasswordPrompt("Repository password:")
		os.Setenv("RESTIC_PASSWORD", password)
	}
}

func Backup(repository_path string, backup_location string, excluded []string) error {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Backup of ", backup_location, "failed:\n", err)
		}
	}()
	// Check for empty dir
	files, err := ioutil.ReadDir(backup_location)
	if err != nil || len(files) == 0 {
		fmt.Println("Skip empty directory ", backup_location)
		return nil
	}

	// Find and exclude child directories
	dirs := Scan(backup_location)
	dirs = LookForFiles(dirs)
	for index := range dirs {
		excluded = append(excluded, dirs[index+1:]...)
	}

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
	return c.Wait()
}
