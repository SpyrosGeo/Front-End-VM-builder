package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const configFileName = ".fe-build-settings.conf"

func getConfigFilePath() (string, error) {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Construct the path to the config file
	configFilePath := filepath.Join(homeDir, configFileName)
	return configFilePath, nil
}

func readConfigFile(filePath string) (map[string]string, error) {
	config := make(map[string]string)

	// Open the config file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Split each line by ": " to get key-value pairs
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

func initConfig() {
	// Get the path to the config file
	configFilePath, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("Error getting config file path: %v\n", err)
		return
	}

	// Delete the previous config file if it exists
	if _, err := os.Stat(configFilePath); err == nil {
		err := os.Remove(configFilePath)
		if err != nil {
			fmt.Printf("Error deleting previous config file: %v\n", err)
			return
		}
	}

	// Prompt the user to set SSH password
	fmt.Println("Enter SSH password:")
	reader := bufio.NewReader(os.Stdin)
	sshPassword, _ := reader.ReadString('\n')

	// Prompt the user to set project directory
	fmt.Println("Enter project directory:")
	projectDir, _ := reader.ReadString('\n')

	// Remove newline characters
	sshPassword = strings.TrimSpace(sshPassword)
	projectDir = strings.TrimSpace(projectDir)

	// Create the config file
	err = createConfigFile(configFilePath, sshPassword, projectDir)
	if err != nil {
		fmt.Printf("Error creating config file: %v\n", err)
		return
	}

	fmt.Println("Configuration file initialized successfully.")
}

func createConfigFile(filePath, sshPassword, projectDir string) error {
	// Create the config file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write SSH password and project directory to the file
	_, err = fmt.Fprintf(file, "ssh_password: %s\n", sshPassword)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, "project_dir: %s\n", projectDir)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Check if there are no arguments
	if len(os.Args) == 1 {
		// Display help message
		fmt.Println("Help:")
		fmt.Println("1) init: initialize the settings for the build")
		fmt.Println("2) To run the script, please use './fe-build <VM_Name_or_IP> <Branch_Name> <Type_of_Build>'")
		return
	}

	// Handle other cases based on the provided arguments
	args := os.Args
	if len(args) == 2 && args[1] == "init" {
		// Initialization logic
		initConfig()
	} else if len(args) == 4 {
		// Check if the first argument is an IP address
		vmName := args[1]
		branchName := args[2]
		buildType := args[3]
		ip := net.ParseIP(vmName)
		if ip == nil {
			// The first argument is not an IP address, proceed with ping
			remoteAddress := fmt.Sprintf("%s.public.gr", vmName)
			ipAddr, err := getPingIPAddress(remoteAddress)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("IP address of %s: %s\n", vmName, ipAddr)
		} else {
			// The first argument is an IP address
			ipAddr := ip.String()
			fmt.Printf("Using provided IP address: %s\n", ipAddr)
		}
		BuildProject(ip.String(), branchName, buildType)

	} else {
		fmt.Println("Invalid arguments. Please see the help message.")
	}
}

func getPingIPAddress(address string) (string, error) {
	out, err := exec.Command("ping", "-c", "1", address).Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("unable to determine IP address")
	}

	parts := strings.Split(lines[1], " ")
	if len(parts) < 2 {
		return "", fmt.Errorf("unable to determine IP address")
	}

	return strings.Trim(parts[1], "()"), nil
}
