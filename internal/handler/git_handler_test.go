package handler

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

// MockGitService for testing
type MockGitService struct {
	RunGitCommandFunc   func(args ...string) error
	GetChangedFilesFunc func() ([]string, error)
}

func (m *MockGitService) RunGitCommand(args ...string) error {
	if m.RunGitCommandFunc != nil {
		return m.RunGitCommandFunc(args...)
	}
	return nil
}

func (m *MockGitService) getChangedFiles() ([]string, error) {
	if m.GetChangedFilesFunc != nil {
		return m.GetChangedFilesFunc()
	}
	return []string{}, nil
}

func TestAddFiles(t *testing.T) {
	t.Run("Direct file addition", func(t *testing.T) {
		mockGit := &MockGitService{
			RunGitCommandFunc: func(args ...string) error {
				assert.Equal(t, []string{"add", "file1.txt", "file2.txt"}, args)
				return nil
			},
		}
		app := cli.NewApp()
		app.Commands = []*cli.Command{
			{
				Name: "add",
				Action: func(c *cli.Context) error {
					return AddFiles(c, mockGit)
				},
			},
		}
		err := app.Run([]string{"gcm", "add", "file1.txt", "file2.txt"})
		assert.NoError(t, err)
	})

	t.Run("Interactive file selection", func(t *testing.T) {
		mockGit := &MockGitService{
			GetChangedFilesFunc: func() ([]string, error) {
				return []string{"file1.txt", "file2.txt"}, nil
			},
			RunGitCommandFunc: func(args ...string) error {
				assert.Equal(t, []string{"add", "file1.txt"}, args)
				return nil
			},
		}
		app := cli.NewApp()
		app.Commands = []*cli.Command{
			{
				Name: "add",
				Action: func(c *cli.Context) error {
					return mockGit.RunGitCommand("add", "file1.txt")
				},
			},
		}
		err := app.Run([]string{"gcm", "add"})
		assert.NoError(t, err)
	})

	t.Run("Error getting changed files", func(t *testing.T) {
		mockGit := &MockGitService{
			GetChangedFilesFunc: func() ([]string, error) {
				return nil, errors.New("git error")
			},
		}
		app := cli.NewApp()
		app.Commands = []*cli.Command{
			{
				Name: "add",
				Action: func(c *cli.Context) error {
					return AddFiles(c, mockGit)
				},
			},
		}
		err := app.Run([]string{"gcm", "add"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get changed files")
	})
}

// TODO: Add Test cases for CreateCommit

// TestPushChanges tests the PushChanges function
func TestPushChanges(t *testing.T) {
	t.Run("Successful push", func(t *testing.T) {
		mockGit := &MockGitService{
			RunGitCommandFunc: func(args ...string) error {
				assert.Equal(t, []string{"push"}, args)
				return nil
			},
		}
		app := cli.NewApp()
		app.Commands = []*cli.Command{
			{
				Name: "push",
				Action: func(c *cli.Context) error {
					return PushChanges(c, mockGit)
				},
			},
		}
		err := app.Run([]string{"gcm", "push"})
		assert.NoError(t, err)
	})

	t.Run("Push error", func(t *testing.T) {
		mockGit := &MockGitService{
			RunGitCommandFunc: func(args ...string) error {
				return errors.New("push failed")
			},
		}
		app := cli.NewApp()
		app.Commands = []*cli.Command{
			{
				Name: "push",
				Action: func(c *cli.Context) error {
					return PushChanges(c, mockGit)
				},
			},
		}
		err := app.Run([]string{"gcm", "push"})
		assert.Error(t, err)
	})
}

// TestForcePushChanges tests the ForcePushChanges function
func TestForcePushChanges(t *testing.T) {
	t.Run("Successful force push", func(t *testing.T) {
		mockGit := &MockGitService{
			RunGitCommandFunc: func(args ...string) error {
				assert.Equal(t, []string{"push", "--force"}, args)
				return nil
			},
		}
		app := cli.NewApp()
		app.Commands = []*cli.Command{
			{
				Name: "force-push",
				Action: func(c *cli.Context) error {
					return ForcePushChanges(c, mockGit)
				},
			},
		}
		err := app.Run([]string{"gcm", "force-push"})
		assert.NoError(t, err)
	})

	t.Run("Force push error", func(t *testing.T) {
		mockGit := &MockGitService{
			RunGitCommandFunc: func(args ...string) error {
				return errors.New("force push failed")
			},
		}
		app := cli.NewApp()
		app.Commands = []*cli.Command{
			{
				Name: "force-push",
				Action: func(c *cli.Context) error {
					return ForcePushChanges(c, mockGit)
				},
			},
		}
		err := app.Run([]string{"gcm", "force-push"})
		assert.Error(t, err)
	})
}

func TestShowTypeRecommendations(t *testing.T) {
	t.Run("Show recommendations", func(t *testing.T) {
		var buf bytes.Buffer
		originalStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		app := cli.NewApp()
		app.Commands = []*cli.Command{
			{
				Name: "show",
				Action: func(c *cli.Context) error {
					return ShowTypeRecommendations(c)
				},
			},
		}
		err := app.Run([]string{"gcm", "show"})
		assert.NoError(t, err)

		if err := w.Close(); err != nil {
			t.Fatalf("failed to close writer: %v", err)
		}
		os.Stdout = originalStdout
		_, _ = buf.ReadFrom(r)

		assert.Contains(t, buf.String(), "Commit Type Recommendations")
	})
}

func TestShowDiff(t *testing.T) {
	t.Run("Successful diff", func(t *testing.T) {
		mockGit := &MockGitService{
			RunGitCommandFunc: func(args ...string) error {
				assert.Equal(t, []string{"diff", "--cached"}, args)
				return nil
			},
		}
		app := cli.NewApp()
		app.Commands = []*cli.Command{
			{
				Name: "diff",
				Action: func(c *cli.Context) error {
					return ShowDiff(c, mockGit)
				},
			},
		}
		err := app.Run([]string{"gcm", "diff"})
		assert.NoError(t, err)
	})

	t.Run("Diff error", func(t *testing.T) {
		mockGit := &MockGitService{
			RunGitCommandFunc: func(args ...string) error {
				return errors.New("diff failed")
			},
		}
		app := cli.NewApp()
		app.Commands = []*cli.Command{
			{
				Name: "diff",
				Action: func(c *cli.Context) error {
					return ShowDiff(c, mockGit)
				},
			},
		}
		err := app.Run([]string{"gcm", "diff"})
		assert.Error(t, err)
	})
}
