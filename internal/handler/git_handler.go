package handler

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

type GitService interface {
	RunGitCommand(args ...string) error
	getChangedFiles() ([]string, error)
}

var commitTypes = []string{
	"feat", "fix", "docs", "style",
	"refactor", "test", "chore", "perf",
	"ci", "build", "revert",
}

func AddFiles(c *cli.Context, git GitService) error {
	files := c.Args().Slice()
	if len(files) == 0 {
		changedFiles, err := git.getChangedFiles()
		if err != nil {
			return fmt.Errorf("failed to get changed files: %w", err)
		}

		var selectedFiles []string
		prompt := &survey.MultiSelect{
			Message: "Select files to stage:",
			Options: changedFiles,
		}
		if err := survey.AskOne(prompt, &selectedFiles); err != nil {
			return fmt.Errorf("failed to select files: %w", err)
		}

		if len(selectedFiles) == 0 {
			fmt.Println("No files selected")
			return nil
		}
		return git.RunGitCommand(append([]string{"add"}, selectedFiles...)...)
	}
	return git.RunGitCommand(append([]string{"add"}, files...)...)
}

// Create Commit
func CreateCommit(c *cli.Context, git GitService) error {
	var commitType string
	promptType := &survey.Select{
		Message: "Select commit type:",
		Options: commitTypes,
	}
	if err := survey.AskOne(promptType, &commitType); err != nil {
		return err
	}

	var scope string
	promptScope := &survey.Input{
		Message: "Enter scope (optional, e.g., 'ci', 'database'):",
	}
	if err := survey.AskOne(promptScope, &scope); err != nil {
		return err
	}

	var message string
	promptMessage := &survey.Input{
		Message: "Enter commit message:",
	}
	if err := survey.AskOne(promptMessage, &message); err != nil {
		return err
	}

	commitMsg := commitType
	if scope != "" {
		commitMsg += fmt.Sprintf("(%s)", scope)
	}
	commitMsg += ": " + message

	return git.RunGitCommand("commit", "-m", commitMsg)
}

// PushChanges handles git push
func PushChanges(c *cli.Context, git GitService) error {
	return git.RunGitCommand("push")
}

// ForcePushChanges handles git push --force
func ForcePushChanges(c *cli.Context, git GitService) error {
	return git.RunGitCommand("push", "--force")
}

// Show commit type recommendations
var typeDescriptions = map[string]string{
	"feat":     "A new feature",
	"fix":      "A bug fix",
	"docs":     "Documentation only changes",
	"style":    "Changes that do not affect the meaning of the code",
	"refactor": "A code change that neither fixes a bug nor adds a feature",
	"test":     "Adding missing tests or correcting existing tests",
	"chore":    "Changes to the build process or auxiliary tools",
	"perf":     "A code change that improves performance",
	"ci":       "Changes to CI configuration files and scripts",
	"build":    "Changes that affect the build system or external dependencies",
	"revert":   "Reverts a previous commit",
}

func ShowTypeRecommendations(c *cli.Context) error {
	fmt.Println("Commit Type Recommendations:")
	for commitType, description := range typeDescriptions {
		fmt.Printf("- %s: %s\n", commitType, description)
	}
	return nil
}

// ShowDiff displays the diff of staged Changes
func ShowDiff(c *cli.Context, git GitService) error {
	return git.RunGitCommand("diff", "--cached")
}
