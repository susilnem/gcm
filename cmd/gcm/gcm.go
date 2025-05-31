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
			{
				Name: "diff",
				Aliases: []string{"d"},
				Usage: "Show the diff of staged changes",
				Action: handler.ShowDiff,
			},
			// Profile management Commands
			{
				Name:    "profile",
				Usage:   "Manage Git user profiles",
				Subcommands: []*cli.Command{
					{
						Name:    "add",
						Usage: "Add a new profile(name, username, email)",
						ArgsUsage: "[profile-name] [user-name] [email]",
						Action:  handler.AddProfile,
					},
					{
						Name:    "list",
						Usage:   "List all profiles",
						Action:  handler.ListProfiles,
					},
					{
						Name:    "use",
						Usage:   "Switch to a profile",
						ArgsUsage: "[profile_name]",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "global",
								Usage: "Apply profile globally",
								DefaultText: "false",
							},
						},
						Action: handler.UseProfile,
					},
					{
						Name: "remove",
						Usage: "Remove a profile",
						ArgsUsage: "[profile_name]",
						Action: handler.RemoveProfile,
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
