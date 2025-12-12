package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	stripmd "github.com/writeas/go-strip-markdown/v2"
)

func getSummaryFilePath(url string) (string, error) {
	dataDir, err := GetDataDir()
	if err != nil {
		return "", err
	}

	summariesDir := filepath.Join(dataDir, "summaries")
	if err := os.MkdirAll(summariesDir, 0755); err != nil {
		return "", err
	}

	hash := md5.Sum([]byte(url))
	filename := hex.EncodeToString(hash[:]) + ".md"
	return filepath.Join(summariesDir, filename), nil
}

func StoreSummarization(url, summary string) error {
	summaryPath, err := getSummaryFilePath(url)
	if err != nil {
		return err
	}

	timestamp := time.Now().Format(time.RFC3339)
	content := fmt.Sprintf("# Summary for: %s\n\nGenerated: %s\n\n---\n\n%s\n", url, timestamp, summary)

	return os.WriteFile(summaryPath, []byte(content), 0644)
}

func getFocusFilePath(url string) (string, error) {
	dataDir, err := GetDataDir()
	if err != nil {
		return "", err
	}

	focusDir := filepath.Join(dataDir, "focus")
	if err := os.MkdirAll(focusDir, 0755); err != nil {
		return "", err
	}

	hash := md5.Sum([]byte(url))
	filename := hex.EncodeToString(hash[:]) + ".md"
	return filepath.Join(focusDir, filename), nil
}

func StoreFocusedContent(url, focused string) error {
	focusPath, err := getFocusFilePath(url)
	if err != nil {
		return err
	}
	return os.WriteFile(focusPath, []byte(focused), 0644)
}

func StoreCombinedFocusedContent(focused string) error {
	dataDir, err := GetDataDir()
	if err != nil {
		return err
	}

	focusDir := filepath.Join(dataDir, "focus")
	if err := os.MkdirAll(focusDir, 0755); err != nil {
		return err
	}

	focusPath := filepath.Join(focusDir, "combined.md")
	timestamp := time.Now().Format(time.RFC3339)
	content := fmt.Sprintf("# Combined Focus Summary\n\nGenerated: %s\n\n---\n\n%s\n", timestamp, focused)
	return os.WriteFile(focusPath, []byte(content), 0644)
}

func GetFocusedContent(url string) (string, error) {
	focusPath, err := getFocusFilePath(url)
	if err != nil {
		return "", err
	}
	
	data, err := os.ReadFile(focusPath)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ShowSummarizations() error {
	dataDir, err := GetDataDir()
	if err != nil {
		return err
	}

	summariesDir := filepath.Join(dataDir, "summaries")
	if _, err := os.Stat(summariesDir); os.IsNotExist(err) {
		fmt.Println("No summarizations found")
		return nil
	}

	files, err := os.ReadDir(summariesDir)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		fmt.Println("No summarizations found")
		return nil
	}

	tempFile := filepath.Join(dataDir, "view.md")
	var content string

	config, err := LoadConfig()
	if err == nil && config.FocusTopics != "" {
		content += fmt.Sprintf("═══════════════════════════════════════════════════════════════════\n")
		content += fmt.Sprintf("  FOCUS: %s\n", config.FocusTopics)
		content += fmt.Sprintf("═══════════════════════════════════════════════════════════════════\n\n")

		focusDir := filepath.Join(dataDir, "focus")
		if focusFiles, err := os.ReadDir(focusDir); err == nil {
			for _, focusFile := range focusFiles {
				if filepath.Ext(focusFile.Name()) != ".md" {
					continue
				}
				focusPath := filepath.Join(focusDir, focusFile.Name())
				focusContent, err := os.ReadFile(focusPath)
				if err == nil {
					content += markdownToPlainText(string(focusContent)) + "\n\n"
				}
			}
		}
		content += "───────────────────────────────────────────────────────────────────\n\n"
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".md" {
			continue
		}

		filePath := filepath.Join(summariesDir, file.Name())
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		content += markdownToPlainText(string(fileContent)) + "\n\n"
		content += "───────────────────────────────────────────────────────────────────\n\n"
	}

	if err := os.WriteFile(tempFile, []byte(content), 0644); err != nil {
		return err
	}

	pager := "less"
	if os.Getenv("PAGER") != "" {
		pager = os.Getenv("PAGER")
	}

	exec := &execCmd{name: pager, args: []string{tempFile}}
	return exec.runWait()
}

func ShowSummarizationsHTML() error {
	dataDir, err := GetDataDir()
	if err != nil {
		return err
	}

	summariesDir := filepath.Join(dataDir, "summaries")
	if _, err := os.Stat(summariesDir); os.IsNotExist(err) {
		fmt.Println("No summarizations found")
		return nil
	}

	files, err := os.ReadDir(summariesDir)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		fmt.Println("No summarizations found")
		return nil
	}

	htmlFile := filepath.Join(dataDir, "view.html")
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>nub</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: Verdana, Geneva, sans-serif;
            font-size: 10pt;
            color: #000;
            background: #f6f6ef;
            padding: 8px;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: #f6f6ef;
        }
        header {
            background: #dc94ba;
            padding: 2px 4px;
            margin-bottom: 10px;
        }
        h1 {
            font-size: 11pt;
            font-weight: bold;
            color: #000;
            display: inline;
        }
        .summary {
            background: #fff;
            padding: 8px;
            margin-bottom: 8px;
            border: 1px solid #e0e0e0;
        }
        h2 {
            font-size: 11pt;
            font-weight: bold;
            margin: 8px 0 4px 0;
            color: #000;
        }
        h3 {
            font-size: 10pt;
            font-weight: bold;
            margin: 6px 0 3px 0;
            color: #000;
        }
        p {
            margin: 6px 0;
            line-height: 1.4;
        }
        a {
            color: #000;
            text-decoration: underline;
        }
        a:visited { color: #828282; }
        code {
            font-family: monospace;
            font-size: 9pt;
            background: #f0f0f0;
            padding: 1px 3px;
        }
        pre {
            background: #f5f5f5;
            border: 1px solid #ddd;
            padding: 8px;
            overflow-x: auto;
            margin: 8px 0;
            font-size: 9pt;
            line-height: 1.3;
        }
        pre code {
            background: none;
            padding: 0;
        }
        ul, ol {
            margin: 6px 0 6px 20px;
        }
        li {
            margin: 2px 0;
            line-height: 1.4;
        }
        blockquote {
            border-left: 2px solid #ccc;
            padding-left: 10px;
            margin: 6px 0;
            color: #555;
        }
        hr {
            border: none;
            border-top: 1px solid #ccc;
            margin: 10px 0;
        }
        strong { font-weight: bold; }
        em { font-style: italic; }
        .meta {
            font-size: 8pt;
            color: #828282;
            margin-bottom: 4px;
        }
        .focus-section {
            background: #fce4f0;
            padding: 10px;
            margin-bottom: 10px;
            border: 1px solid #dc94ba;
        }
        .focus-section h2 {
            color: #c2608a;
            margin-top: 0;
        }
        @media (max-width: 700px) {
            body { padding: 4px; font-size: 9pt; }
            .summary { padding: 6px; }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>nub</h1>
        </header>
`

	config, err := LoadConfig()
	if err == nil && config.FocusTopics != "" {
		html += fmt.Sprintf(`        <div class="focus-section">
            <h2>Focus: %s</h2>
`, config.FocusTopics)

		focusDir := filepath.Join(dataDir, "focus")
		if focusFiles, err := os.ReadDir(focusDir); err == nil {
			for _, focusFile := range focusFiles {
				if filepath.Ext(focusFile.Name()) != ".md" {
					continue
				}
				focusPath := filepath.Join(focusDir, focusFile.Name())
				focusContent, err := os.ReadFile(focusPath)
				if err == nil {
					html += markdownToHTML(string(focusContent))
				}
			}
		}

		html += `        </div>
`
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".md" {
			continue
		}

		filePath := filepath.Join(summariesDir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		html += `<div class="summary">
` + markdownToHTML(string(content)) + `
</div>
`
	}

	html += `    </div>
</body>
</html>`

	if err := os.WriteFile(htmlFile, []byte(html), 0644); err != nil {
		return err
	}

	fmt.Printf("Opening summaries in browser...\n")
	fmt.Printf("File: %s\n", htmlFile)

	return openInBrowser(htmlFile)
}

func markdownToHTML(md string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

func markdownToPlainText(md string) string {
	text := stripmd.Strip(md)
	
	lines := strings.Split(text, "\n")
	var formatted []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			formatted = append(formatted, "")
			continue
		}
		
		if strings.HasPrefix(line, "Summary for:") || strings.HasPrefix(line, "Generated:") {
			formatted = append(formatted, line)
		} else if len(line) > 0 {
			wrapped := wrapText(line, 78)
			formatted = append(formatted, wrapped...)
		}
	}
	
	return strings.Join(formatted, "\n")
}

func wrapText(text string, width int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}
	
	var lines []string
	var currentLine string
	
	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	
	return lines
}

func openInBrowser(path string) error {
	var cmd string
	switch {
	case fileExists("/usr/bin/open"):
		cmd = "open"
	case fileExists("/usr/bin/xdg-open"):
		cmd = "xdg-open"
	case fileExists("/usr/bin/wslview"):
		cmd = "wslview"
	default:
		return fmt.Errorf("no browser opener found, please open manually: %s", path)
	}

	exec := &execCmd{name: cmd, args: []string{path}}
	return exec.run()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ClearAllData() error {
	dataDir, err := GetDataDir()
	if err != nil {
		return err
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		return nil
	}

	err = os.RemoveAll(dataDir)
	if err != nil {
		return fmt.Errorf("failed to remove data directory: %v", err)
	}

	fmt.Printf("Removed: %s\n", dataDir)
	return nil
}
