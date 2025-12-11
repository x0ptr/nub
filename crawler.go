package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

func CrawlWebsite(url string) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	reader, err := charset.NewReader(resp.Body, contentType)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getCacheFilePath(url string) (string, error) {
	dataDir, err := GetDataDir()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(dataDir, "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	hash := md5.Sum([]byte(url))
	filename := hex.EncodeToString(hash[:]) + ".html"
	return filepath.Join(cacheDir, filename), nil
}

func IsCached(url string) (bool, error) {
	cachePath, err := getCacheFilePath(url)
	if err != nil {
		return false, err
	}

	info, err := os.Stat(cachePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	age := time.Since(info.ModTime())
	if age > 24*time.Hour {
		return false, nil
	}

	return true, nil
}

func GetCachedContent(url string) (string, error) {
	cachePath, err := getCacheFilePath(url)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func CacheContent(url, content string) error {
	cachePath, err := getCacheFilePath(url)
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, []byte(content), 0644)
}

func ClearCache() error {
	dataDir, err := GetDataDir()
	if err != nil {
		return err
	}

	cacheDir := filepath.Join(dataDir, "cache")
	
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return nil
	}

	err = os.RemoveAll(cacheDir)
	if err != nil {
		return err
	}

	return os.MkdirAll(cacheDir, 0755)
}

func extractTextFromHTML(html string) string {
	text := html
	text = strings.ReplaceAll(text, "<script", "\n<script")
	text = strings.ReplaceAll(text, "</script>", "</script>\n")
	text = strings.ReplaceAll(text, "<style", "\n<style")
	text = strings.ReplaceAll(text, "</style>", "</style>\n")
	
	inScript := false
	inStyle := false
	result := strings.Builder{}
	
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "<script") {
			inScript = true
			continue
		}
		if strings.Contains(line, "</script>") {
			inScript = false
			continue
		}
		if strings.Contains(line, "<style") {
			inStyle = true
			continue
		}
		if strings.Contains(line, "</style>") {
			inStyle = false
			continue
		}
		
		if !inScript && !inStyle {
			result.WriteString(line)
			result.WriteString("\n")
		}
	}
	
	text = result.String()
	
	for strings.Contains(text, "<") && strings.Contains(text, ">") {
		start := strings.Index(text, "<")
		end := strings.Index(text[start:], ">")
		if end == -1 {
			break
		}
		text = text[:start] + " " + text[start+end+1:]
	}
	
	lines := strings.Split(text, "\n")
	cleanLines := make([]string, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && len(line) > 3 {
			cleanLines = append(cleanLines, line)
		}
	}
	
	return strings.Join(cleanLines, "\n")
}
