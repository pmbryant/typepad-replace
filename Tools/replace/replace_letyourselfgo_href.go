package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	rootDir := "/Users/peterbryant/Documents/Codebase/typepad-replace/lyg-blog/pmbryant.typepad.com/letyourselfgo/2017/10" // Change as needed
	replaceHref := "https://pmbryant.com/x/letyourselfgo/index.html"                                                       // Configurable replacement value
	pattern := regexp.MustCompile(`href=["']([^"']*\\&quot;6Ldg1s4SAAAAAEvvZX2ILFkWp7KB-jjdL4v0JV2e\\&quot;\.html)["']`)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			input, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			content := string(input)
			newContent := pattern.ReplaceAllStringFunc(content, func(m string) string {
				return pattern.ReplaceAllString(m, "href=\""+replaceHref+"\"")
			})
			if newContent != content {
				err = os.WriteFile(path, []byte(newContent), info.Mode())
				if err != nil {
					return err
				}
				fmt.Printf("Updated: %s\n", path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
