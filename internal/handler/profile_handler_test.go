package handler

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/urfave/cli/v2"
)

// Setup a temporary profile file for testing
func setupTestProfileFile(t *testing.T, content string) string {
	homeDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", homeDir)

	profilePath := filepath.Join(homeDir, profileFile)
	if content != "" {
		if err := os.WriteFile(profilePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test profile file: %v", err)
		}
	}

	t.Cleanup(func() {
		// Clean up the temporary profile file after the test
		if err := os.Remove(profilePath); err != nil && !os.IsNotExist(err) {
			t.Fatalf("Failed to remove test profile file: %v", err)
		}
		os.Setenv("HOME", originalHome)
	})

	return profilePath
}

func TestLoadProfiles(t *testing.T) {
	// Test with an empty profile file
	setupTestProfileFile(t, "")
	store, err := LoadProfiles()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(store.Profiles) != 0 {
		t.Errorf("Expected empty profiles, got %d profiles", len(store.Profiles))
	}

	// Test with a valid profile file
	content := `{"profiles":{"work":{"name":"John Doe","email":"john@work.com"}}}`
	_ = setupTestProfileFile(t, content)

	store, err = LoadProfiles()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(store.Profiles) != 1 {
		t.Errorf("Expected 1 profile, got %d profiles", len(store.Profiles))
	}

	profile, exists := store.Profiles["work"]
	if !exists {
		t.Errorf("Expected profile 'work' to exist")
	}
	if profile.Name != "John Doe" || profile.Email != "john@work.com" {
		t.Errorf("Expected profile 'work' to have name 'John Doe' and email 'john@work.com', got %v", profile)
	}
}

// Test AddProfile function
func TestAddProfile(t *testing.T) {
	// Setup a temporary profile file
	profilePath := setupTestProfileFile(t, "")
	// Setup the CLI context with appropriate flags
	set := flag.NewFlagSet("test", 0)
	set.String("profile", "work", "")
	set.String("name", "John Doe", "")
	set.String("email", "john@work.com", "")
	set.Parse([]string{"work", "John Doe", "john@work.com"})

	app := cli.NewApp()
	ctx := cli.NewContext(app, set, nil)

	err := AddProfile(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	store, err := LoadProfiles()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(store.Profiles) != 1 {
		t.Errorf("Expected 1 profile, got %d profiles", len(store.Profiles))
	}
	profile, exists := store.Profiles["work"]
	if !exists || profile.Name != "John Doe" || profile.Email != "john@work.com" {
		t.Errorf("Expected profile work with correct details, got %v", profile)
	}

	// Verify the profile file was updated
	profileFileContent, err := os.ReadFile(profilePath)
	if err != nil {
		t.Fatalf("Failed to read profile file: %v", err)
	}
	if string(profileFileContent) == "" {
		t.Errorf("Expected profile file to be updated, but it is empty")
	}
}

// Test Duplicate Profiles
func TestAddDuplicateProfile(t *testing.T) {
	content := `{"profiles":{"work":{"name":"John Doe","email":"john@work.com"}}}`
	setupTestProfileFile(t, content)

	set := flag.NewFlagSet("test", 0)
	set.String("profile", "work", "")
	set.String("name", "John Doe", "")
	set.String("email", "john@work.com", "")
	set.Parse([]string{"work", "John Doe", "john@work.com"})

	app := cli.NewApp()
	ctx := cli.NewContext(app, set, nil)

	err := AddProfile(ctx)
	if err == nil || err.Error() != "profile 'work' already exists" {
		t.Fatalf("Expected error for duplicate profile, got %v", err)
	}
}

// Test RemoveProfile function
func TestRemoveProfile(t *testing.T) {
	content := `{"profiles":{"work":{"name":"John Doe","email":"john@work.com"}}}`
	setupTestProfileFile(t, content)

	set := flag.NewFlagSet("test", 0)
	set.String("profile", "work", "")
	set.String("name", "John Doe", "")
	set.String("email", "john@work.com", "")
	set.Parse([]string{"work", "John Doe", "john@work.com"})

	app := cli.NewApp()
	ctx := cli.NewContext(app, set, nil)

	err := RemoveProfile(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	store, err := LoadProfiles()
	if err != nil {
		t.Fatalf("Expected no error loading profile, got %v", err)
	}

	if len(store.Profiles) != 0 {
		t.Errorf("Expected 0 profiles after removal, got %d profiles", len(store.Profiles))
	}
}
