package handler

import (
	"fmt"
	"net/mail"
	"os"
	"os/exec"
	"strings"
)

type RealGitService struct{}

func (r *RealGitService) RunGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *RealGitService) getChangedFiles() ([]string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting changed files:", err)
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var files []string
	for _, line := range lines {
		if len(line) > 3 {
			fmt.Println("Line:", line)
			files = append(files, strings.TrimSpace(line[3:]))
		}
	}
	return files, nil
}

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

var DefaultGitService GitService = &RealGitService{}
