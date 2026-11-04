package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func download(plugin string) string {

	url := fmt.Sprintf("http://192.168.4.23:8080/plugins/%s.zip", plugin)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Sprintf("Failed to create request %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Sprintf("request failed %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Download rejected, status code: %s", resp.Status)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Sprintf("failed to get home dir: %s", err)
	}

	destPath := filepath.Join(homeDir, "Downloads", plugin+".zip")

	outFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Sprintf("Could not create file: %s", err)
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	return "Thank you for downloading! \n The plugin will be in your downloads folder. \n Hit backspace to return to the main menu"
}
