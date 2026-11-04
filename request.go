package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func download(plugin string) error {

	url := fmt.Sprintf("http://192.168.4.23:8080/plugins/%s", plugin)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Failed to create request %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Download rejected, status code: %s", resp.Status)
	}

	outFile, err := os.Create("./Downloads/" + plugin)
	if err != nil {
		return fmt.Errorf("Could not create file: %w", err)
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}
