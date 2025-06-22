// ProfileHandler handles user profile related commands
package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

type Profile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ProfileStore struct {
	Profiles map[string]Profile `json:"profiles"`
}

const profileFile = ".gcm_profiles.json"

// getProfilePath returns the path to the profile configuration file
func getProfilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, profileFile), nil
}

// LoadProfiles reads the profile configuration from the file
func LoadProfiles() (ProfileStore, error) {
	profilePath, err := getProfilePath()
	if err != nil {
		return ProfileStore{}, err
	}

	// Check if the profile file exists
	store := ProfileStore{Profiles: make(map[string]Profile)}
	data, err := os.ReadFile(profilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file does not exist, return an empty store
			return store, nil
		}
		return store, fmt.Errorf("failed to read profile file: %w", err)
	}
	// Parse the JSON data into the ProfileStore
	if err := json.Unmarshal(data, &store); err != nil {
		return ProfileStore{}, fmt.Errorf("failed to parse profiles: %w", err)
	}
	return store, nil
}

// SaveProfiles writes the profile configuration to the file
func SaveProfiles(store ProfileStore) error {
	profilePath, err := getProfilePath()
	if err != nil {
		return err
	}

	// Ensure the directory exists
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profiles: %w", err)
	}

	if err := os.WriteFile(profilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write profile file: %w", err)
	}
	return nil
}

// AddProfile adds a new profile to the profile store
func AddProfile(c *cli.Context) error {
	profileName := c.Args().Get(0)
	userName := c.Args().Get(1)
	email := c.Args().Get(2)

	if profileName == "" || userName == "" || email == "" {
		return fmt.Errorf("profile name, username, and email are required")
	}

	if !ValidEmail(email) {
		return fmt.Errorf("invalid email address: %s", email)
	}

	store, err := LoadProfiles()
	if err != nil {
		return fmt.Errorf("failed to load profiles: %w", err)
	}

	// Check if the profile already exists
	if _, exists := store.Profiles[profileName]; exists {
		return fmt.Errorf("profile '%s' already exists", profileName)
	}

	store.Profiles[profileName] = Profile{
		Name:  userName,
		Email: email,
	}

	if err := SaveProfiles(store); err != nil {
		return fmt.Errorf("failed to save profiles: %w", err)
	}

	fmt.Printf("Profile '%s' added successfully\n", profileName)
	return nil
}

// ListProfiles lists all profiles in the profile store
func ListProfiles(c *cli.Context) error {
	store, err := LoadProfiles()
	if err != nil {
		return fmt.Errorf("failed to load profiles: %w", err)
	}

	if len(store.Profiles) == 0 {
		fmt.Println("No profiles found")
		return nil
	}

	fmt.Println("Available profiles:")
	for name, profile := range store.Profiles {
		fmt.Printf("- %s: %s <%s>\n", name, profile.Name, profile.Email)
	}
	return nil
}

// UseProfile switches to a specified profile
func UseProfile(c *cli.Context, git GitService) error {
	profileName := c.Args().Get(0)
	isGlobal := c.Bool("global")

	if profileName == "" {
		return fmt.Errorf("profile name is required")
	}

	store, err := LoadProfiles()
	if err != nil {
		return fmt.Errorf("failed to load profiles: %w", err)
	}

	profile, exists := store.Profiles[profileName]
	if !exists {
		return fmt.Errorf("profile '%s' does not exist", profileName)
	}

	args := []string{"config", "--local"}
	if isGlobal {
		args = []string{"config", "--global"}
	}

	if err := git.RunGitCommand(append(args, "user.name", profile.Name)...); err != nil {
		return fmt.Errorf("failed to set user.name: %w", err)
	}
	if err := git.RunGitCommand(append(args, "user.email", profile.Email)...); err != nil {
		return fmt.Errorf("failed to set user.email: %w", err)
	}

	fmt.Printf(
		"Switched to profile '%s' (%s <%s>) %s\n",
		profileName,
		profile.Name,
		profile.Email,
		map[bool]string{true: " globally", false: " locally"}[isGlobal],
	)

	return nil
}

// RemoveProfile deletes profile from the profile store
func RemoveProfile(c *cli.Context) error {
	profileName := c.Args().Get(0)

	if profileName == "" {
		return fmt.Errorf("profile name is required")
	}

	store, err := LoadProfiles()
	if err != nil {
		return fmt.Errorf("failed to load profiles: %w", err)
	}

	if _, exists := store.Profiles[profileName]; !exists {
		return fmt.Errorf("profile '%s' does not exist", profileName)
	}

	delete(store.Profiles, profileName)

	if err := SaveProfiles(store); err != nil {
		return fmt.Errorf("failed to save profiles: %w", err)
	}

	fmt.Printf("Profile '%s' removed successfully\n", profileName)
	return nil
}
