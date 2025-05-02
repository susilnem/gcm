package gcm

import (
	"fmt"
	"os"

	handler "github.com/susilnem/gcm/internal/handler"
	"github.com/urfave/cli/v2"
)

func Execute() {
	app := &cli.App{
		Version: Version,
		Name:  "gcm",
		Usage: "Git Conventional Commit Manager",
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Stage files for commit",
				Action:  handler.AddFiles,
			},
			{
				Name:    "commit",
				Aliases: []string{"c"},
				Usage:   "Create a conventional commit",
				Action:  handler.CreateCommit,
			},
			{
				Name:    "push",
				Aliases: []string{"p"},
				Usage:   "Push changes to remote",
				Action:  handler.PushChanges,
			},
			{
				Name:    "force-push",
				Aliases: []string{"fp"},
				Usage:   "Force push changes to remote",
				Action:  handler.ForcePushChanges,
			},
			{
				Name:   "show",
				Aliases: []string{"s"},
				Usage:  "Show the last commit message",
				Action: handler.ShowTypeRecommendations,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
