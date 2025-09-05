package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	//rootDir := "/Users/peterbryant/Documents/Codebase/typepad-replace/lyg-blog/pmbryant.typepad.com/letyourselfgo"
	rootDir := "/Users/peterbryant/Documents/Codebase/typepad-replace"

	// Define match and replace pairs
	matchReplace := []struct {
		matchStr   string
		replaceStr string
	}{
		{`href="\&quot;6Ldg1s4SAAAAAEvvZX2ILFkWp7KB-jjdL4v0JV2e\&quot;.html`, `href="../../index.html`}, // These are files in 2021/01 subdir
		{`href="../01/\&quot;6Ldg1s4SAAAAAEvvZX2ILFkWp7KB-jjdL4v0JV2e\&quot;.html`, `href="../../index.html`},
		{`href="2021/01/\&quot;6Ldg1s4SAAAAAEvvZX2ILFkWp7KB-jjdL4v0JV2e\&quot;.html`, `href="index.html`},
		{`href="../2021/01/\&quot;6Ldg1s4SAAAAAEvvZX2ILFkWp7KB-jjdL4v0JV2e\&quot;.html`, `href="../index.html`},
		{`href="../../2021/01/\&quot;6Ldg1s4SAAAAAEvvZX2ILFkWp7KB-jjdL4v0JV2e\&quot;.html`, `href="../../index.html`},
		{`href="../../../2021/01/\&quot;6Ldg1s4SAAAAAEvvZX2ILFkWp7KB-jjdL4v0JV2e\&quot;.html`, `href="../../../index.html`},
		{`href="../../../../2021/01/\&quot;6Ldg1s4SAAAAAEvvZX2ILFkWp7KB-jjdL4v0JV2e\&quot;.html`, `href="../../../../index.html`},
		{`href="atom.xml"`, ``},
		{`href="../atom.xml"`, ``},
		{`href="../../atom.xml"`, ``},
		{`href="../../../atom.xml"`, ``},
		{`href="../../../../atom.xml"`, ``},
		{`href="../../../../../atom.xml"`, ``},
		{`href="https://pmbryant.typepad.com/letyourselfgo/atom.xml"`, ``}, // These only show up in atom.xml files so no replacement
		{`href="https://pmbryant.typepad.com/b_and_b/atom.xml"`, ``},       // These only show up in atom.xml files so no replacement
		{`<a href="https://www.typepad.com/" title="Blog">Blog</a> powered by <a href="https://www.typepad.com/" title="Typepad">Typepad</a>`, `Blog powered by Typepad`},
	}

	var files []string
	// Walk all html files under rootDir
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			files = append(files, path)
		}
		return nil
	})

	for _, path := range files {
		input, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", path, err)
			continue
		}
		content := string(input)
		newContent := content
		for _, mr := range matchReplace {
			newContent = strings.ReplaceAll(newContent, mr.matchStr, mr.replaceStr)
		}
		if newContent != content {
			err = os.WriteFile(path, []byte(newContent), 0644)
			if err != nil {
				fmt.Printf("Error writing %s: %v\n", path, err)
				continue
			}
			fmt.Printf("Updated: %s\n", path)
		}
	}
}
