package main

import (
	"fmt"
	"os"
	"os/exec"
)

func CopyFilesViaSCP(generatedFolder, remoteIP, directoryName string) error {
	// Construct the destination path on the remote server
	remotePath := fmt.Sprintf("/var/www/html/%s", directoryName)
	fmt.Println("Starting scp process on %v:%v", remoteIP, remotePath)
	// Construct the SCP command
	scpCmd := exec.Command("scp", "-r", generatedFolder, fmt.Sprintf("thatguy@%s:%s", remoteIP, remotePath))
	scpCmd.Stdout = os.Stdout
	scpCmd.Stderr = os.Stderr

	// Run the SCP command
	fmt.Printf("Copying files via SCP to remote server: %s", remoteIP)
	if err := scpCmd.Run(); err != nil {
		return fmt.Errorf("error copying files via SCP: %v", err)
	}

	fmt.Println("Files copied successfully.")
	return nil
}
