package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

func main() {
	workspaceDir := "/Users/peterbryant/Documents/Codebase/typepad-replace"
	//baseDownloadDir := "/Users/peterbryant/Documents/Codebase/typepad-replace/downloads" // Change as needed

	hrefSet := make(map[string]struct{})
	anchorPattern := regexp.MustCompile(`<a\s[^>]*href=["']([^"']+)['"]`)

	_ = filepath.Walk(workspaceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			matches := anchorPattern.FindAllStringSubmatch(string(data), -1)
			for _, match := range matches {
				if len(match) > 1 && isTopLevelSubdir(match[1]) {
					hrefSet[match[1]] = struct{}{}
				}
			}
		}
		return nil
	})

	hrefs := make([]string, 0, len(hrefSet))
	for href := range hrefSet {
		hrefs = append(hrefs, href)
	}

	sort.Strings(hrefs)
	for _, href := range hrefs {
		fmt.Println(href)
	}

	//downloadFiles(hrefs, baseDownloadDir)
}

func isTopLevelSubdir(href string) bool {
	// Match: https://pmbryant.typepad.com/.a/..., .shared/..., or photos/...
	pattern := regexp.MustCompile(`^https://pmbryant\.typepad\.com/(\.a|\.shared|photos)/`)
	return pattern.MatchString(href)
}

/*
func downloadFiles(urls []string, baseDir string) {
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36"
	for _, url := range urls {
		// Space out downloads by 0.5 seconds
		time.Sleep(500 * time.Millisecond)

		// Get subpath after domain
		parts := strings.SplitN(url, "pmbryant.typepad.com/", 2)
		if len(parts) != 2 {
			fmt.Printf("Invalid URL: %s\n", url)
			continue
		}
		subPath := parts[1]
		localPath := filepath.Join(baseDir, subPath)

		// Ensure directory exists
		dir := filepath.Dir(localPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Failed to create directory %s: %v\n", dir, err)
			continue
		}

		// Prepare HTTP request with User-Agent
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Failed to create request for %s: %v\n", url, err)
			continue
		}
		req.Header.Set("User-Agent", userAgent)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Failed to download %s: %v\n", url, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Non-OK HTTP status for %s: %s\n", url, resp.Status)
			continue
		}

		out, err := os.Create(localPath)
		if err != nil {
			fmt.Printf("Failed to create file %s: %v\n", localPath, err)
			continue
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			fmt.Printf("Failed to write file %s: %v\n", localPath, err)
			continue
		}

		fmt.Printf("Downloaded %s to %s\n", url, localPath)

	}
}

*/
