package handler

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

var commitTypes = []string{
	"feat", "fix", "docs", "style",
	"refactor", "test", "chore", "perf",
	"ci", "build", "revert",
}

func AddFiles(c *cli.Context) error {
	files := c.Args().Slice()
	if len(files) == 0 {
		changedFiles, err := getChangedFiles()
		if err != nil {
			return err
		}

		var selectedFiles []string
		prompt := &survey.MultiSelect{
			Message: "Select files to stage:",
			Options: changedFiles,
		}
		survey.AskOne(prompt, &selectedFiles)

		if len(selectedFiles) == 0 {
			fmt.Println("No files selected")
			return nil
		}
		return RunGitCommand(append([]string{"add"}, selectedFiles...)...)
	}
	return RunGitCommand(append([]string{"add"}, files...)...)
}


// Create Commit
func CreateCommit(c *cli.Context) error {
	var commitType string
	promptType := &survey.Select{
		Message: "Select commit type:",
		Options: commitTypes,
	}
	survey.AskOne(promptType, &commitType)

	var scope string
	promptScope := &survey.Input{
		Message: "Enter scope (optional, e.g., 'ci', 'database'):",
	}
	survey.AskOne(promptScope, &scope)

	var message string
	promptMessage := &survey.Input{
		Message: "Enter commit message:",
	}
	survey.AskOne(promptMessage, &message)

	commitMsg := commitType
	if scope != "" {
		commitMsg += fmt.Sprintf("(%s)", scope)
	}
	commitMsg += ": " + message

	return RunGitCommand("commit", "-m", commitMsg)
}

// PushChanges handles git push
func PushChanges(c *cli.Context) error {
	return RunGitCommand("push")
}

// ForcePushChanges handles git push --force
func ForcePushChanges(c *cli.Context) error {
	return RunGitCommand("push", "--force")
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
