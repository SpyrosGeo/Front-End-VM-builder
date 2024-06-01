package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func BuildProject(remoteIp, branchName, buildType string) error {
	// Get the path to the config file
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// Navigate to project directory
	config, err := readConfigFile(configFilePath)
	if err != nil {
		return err
	}
	projectDir := config["project_dir"]
	os.Chdir(projectDir)
	if err != nil {
		return err
	}

	// Checkout to the specified branch
	fmt.Printf("Checking out to branch: %s\n", branchName)
	gitCheckoutCmd := exec.Command("git", "checkout", branchName)
	gitCheckoutCmd.Stdout = os.Stdout
	gitCheckoutCmd.Stderr = os.Stderr
	if err := gitCheckoutCmd.Run(); err != nil {
		return err
	}

	// Pull the latest changes
	fmt.Println("Pulling the latest changes")
	gitPullCmd := exec.Command("git", "pull")
	gitPullCmd.Stdout = os.Stdout
	gitPullCmd.Stderr = os.Stderr
	if err := gitPullCmd.Run(); err != nil {
		return err
	}

	// Updating dependencies
	fmt.Println("Running npm install to update dependencies")
	npmInstallCmd := exec.Command("npm", "i", "--legacy-peer-deps")
	npmInstallCmd.Stdout = os.Stdout
	npmInstallCmd.Stderr = os.Stderr
	if err := npmInstallCmd.Run(); err != nil {
		return err
	}
	var npmCmd string
	var remoteFolderName string
	switch buildType {
	case "gr":
		npmCmd = "build:preprod"
		remoteFolderName = "public.gr"
	case "cy":
		npmCmd = "build:preprod_cyprus"
		remoteFolderName = "public.com.cy"
	case "b2b":
		npmCmd = "build:b2b"
		remoteFolderName = "publicbusiness.gr"
	default:
		return fmt.Errorf("Invalid build type: %s", buildType)
	}

	// Run npm build command
	fmt.Printf("Running npm build command: %s\n", npmCmd)
	npmBuildCmd := exec.Command("npm", "run", npmCmd)
	npmBuildCmd.Stdout = os.Stdout
	npmBuildCmd.Stderr = os.Stderr
	if err := npmBuildCmd.Run(); err != nil {
		return err
	}
	generatedDirectory := filepath.Join(projectDir, "dist/pbc/browser")
	CopyFilesViaSCP(generatedDirectory, remoteIp, remoteFolderName)

	fmt.Println("Build process completed successfully.")
	return nil
}
