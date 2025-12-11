package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	helpFlag := flag.Bool("help", false, "Show help")
	
	runMode := flag.Bool("run", false, "Run crawl and summarization once")
	daemonMode := flag.Bool("d", false, "Run in daemon mode")
	stopDaemon := flag.Bool("stop", false, "Stop running daemon")
	showMode := flag.Bool("show", false, "Show summarizations in HTML")
	
	listSources := flag.Bool("list", false, "List all sources")
	addSource := flag.String("add-source", "", "Add a source URL")
	remSource := flag.String("rem-source", "", "Remove a source by ID or URL")
	
	setLLMAPIKey := flag.String("set-llm-api-key", "", "Set LLM API key")
	setLLMAPIURL := flag.String("set-llm-api-url", "", "Set LLM API URL")
	setLLMAPIModel := flag.String("set-llm-api-model", "", "Set LLM API model")
	setScheduleTime := flag.Int("set-schedule-time", 0, "Set schedule time in minutes")
	setPrompt := flag.String("set-prompt", "", "Set custom summarization prompt")
	setFocus := flag.String("set-focus", "", "Set focus topics (comma-separated)")
	
	logsMode := flag.Bool("logs", false, "View logs in pager")
	clearCache := flag.Bool("clear-cache", false, "Clear cached websites")
	clearData := flag.Bool("clear-data", false, "Clear all stored data")
	
	daemonChild := flag.Bool("daemon-child", false, "Internal: daemon child process")

	flag.Parse()

	if *helpFlag {
		showHelp()
		return
	}

	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if *daemonChild {
		if config.ScheduleMinutes == 0 {
			config.ScheduleMinutes = 15
		}
		RunDaemonChild(config)
		return
	}

	if *setLLMAPIKey != "" {
		config.LLMAPIKey = *setLLMAPIKey
		if err := SaveConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("LLM API key set successfully")
		return
	}

	if *setLLMAPIURL != "" {
		config.LLMAPIURL = *setLLMAPIURL
		if err := SaveConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("LLM API URL set successfully")
		return
	}

	if *setLLMAPIModel != "" {
		config.LLMAPIModel = *setLLMAPIModel
		if err := SaveConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("LLM API model set successfully")
		return
	}

	if *setScheduleTime > 0 {
		config.ScheduleMinutes = *setScheduleTime
		if err := SaveConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Schedule time set to %d minutes\n", *setScheduleTime)
		return
	}

	if *setPrompt != "" {
		config.SummaryPrompt = *setPrompt
		if err := SaveConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Summary prompt set successfully")
		return
	}

	if *setFocus != "" {
		config.FocusTopics = *setFocus
		if err := SaveConfig(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Focus topics set to: %s\n", *setFocus)
		return
	}

	if *addSource != "" {
		if err := AddSource(config, *addSource); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding source: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Source added: %s\n", *addSource)
		return
	}

	if *remSource != "" {
		if err := RemoveSource(config, *remSource); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing source: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Source removed: %s\n", *remSource)
		return
	}

	if *listSources {
		ListSources(config)
		return
	}

	if *clearCache {
		if err := ClearCache(); err != nil {
			fmt.Fprintf(os.Stderr, "Error clearing cache: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Cache cleared successfully")
		return
	}

	if *clearData {
		if err := ClearAllData(); err != nil {
			fmt.Fprintf(os.Stderr, "Error clearing data: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("All data cleared successfully")
		return
	}

	if *logsMode {
		if err := ShowLogs(); err != nil {
			fmt.Fprintf(os.Stderr, "Error showing logs: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *daemonMode {
		if config.ScheduleMinutes == 0 {
			config.ScheduleMinutes = 15
		}
		if err := StartDaemon(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting daemon: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *stopDaemon {
		if err := StopDaemon(); err != nil {
			fmt.Fprintf(os.Stderr, "Error stopping daemon: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *runMode {
		if err := runOnce(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *showMode {
		if err := ShowSummarizations(); err != nil {
			fmt.Fprintf(os.Stderr, "Error showing summarizations: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if flag.NFlag() == 0 {
		showHelp()
		return
	}
}

func showHelp() {
	fmt.Println("nub - A website crawler and summarizer")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  nub --run                        Run crawl and summarization once")
	fmt.Println("  nub -d                           Run in daemon mode")
	fmt.Println("  nub --stop                       Stop running daemon")
	fmt.Println("  nub --show                       Show summarizations in HTML")
	fmt.Println()
	fmt.Println("Source Management:")
	fmt.Println("  nub --list                       List all sources")
	fmt.Println("  nub --add-source <url>           Add a source URL")
	fmt.Println("  nub --rem-source <id or url>     Remove a source by ID or URL")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Println("  nub --set-llm-api-key <key>      Set LLM API key")
	fmt.Println("  nub --set-llm-api-url <url>      Set LLM API URL")
	fmt.Println("  nub --set-llm-api-model <model>  Set LLM API model")
	fmt.Println("  nub --set-schedule-time <mins>   Set schedule time in minutes")
	fmt.Println("  nub --set-prompt <text>          Set custom summarization prompt")
	fmt.Println("  nub --set-focus <topics>         Set focus topics (comma-separated)")
	fmt.Println()
	fmt.Println("Utilities:")
	fmt.Println("  nub --logs                       View logs in pager")
	fmt.Println("  nub --clear-cache                Clear cached websites")
	fmt.Println("  nub --clear-data                 Clear all stored data")
	fmt.Println("  nub --help                       Show this help")
	fmt.Println()
	fmt.Println("Files:")
	fmt.Println("  Config file: ~/.config/nub/config.json")
	fmt.Println("  Data directory: ~/.local/nub/")
	fmt.Println()
	fmt.Println("Example config.json:")
	fmt.Println(`  {
    "sources": ["https://example.com", "https://news.ycombinator.com"],
    "llm_api_key": "your-api-key",
    "llm_api_url": "https://api.mistral.ai/v1/chat/completions",
    "llm_api_model": "mistral-small-latest",
    "schedule_minutes": 15,
    "focus_topics": "go,javascript,rust"
  }`)
}

func runOnce(config *Config) error {
	if err := validateConfig(config); err != nil {
		return err
	}

	fmt.Println("Starting crawl and summarization...")
	var allSummaries []string
	for _, source := range config.Sources {
		summary, err := processSource(config, source)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", source, err)
			continue
		}
		if config.FocusTopics != "" && summary != "" {
			allSummaries = append(allSummaries, summary)
		}
	}

	if config.FocusTopics != "" && len(allSummaries) > 0 {
		fmt.Println("Extracting focused content from all summaries...")
		combinedSummaries := strings.Join(allSummaries, "\n\n---\n\n")
		focused, err := ExtractFocusedContent(config, combinedSummaries)
		if err != nil {
			fmt.Printf("Warning: failed to extract focused content: %v\n", err)
		} else if focused != "" && focused != "No relevant content found." {
			if err := StoreCombinedFocusedContent(focused); err != nil {
				fmt.Printf("Warning: failed to store focused content: %v\n", err)
			}
		}
	}

	fmt.Println("Done!")
	return nil
}

func runDaemon(config *Config) {
	ticker := time.NewTicker(time.Duration(config.ScheduleMinutes) * time.Minute)
	defer ticker.Stop()

	if err := runOnce(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error in initial run: %v\n", err)
	}

	for range ticker.C {
		fmt.Printf("\n[%s] Starting scheduled crawl...\n", time.Now().Format(time.RFC3339))
		if err := runOnce(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
}

func validateConfig(config *Config) error {
	if len(config.Sources) == 0 {
		return fmt.Errorf("no sources configured")
	}
	if config.LLMAPIKey == "" {
		return fmt.Errorf("LLM API key not set, use --set-llm-api-key")
	}
	if config.LLMAPIURL == "" {
		return fmt.Errorf("LLM API URL not set, use --set-llm-api-url")
	}
	if config.LLMAPIModel == "" {
		return fmt.Errorf("LLM API model not set, use --set-llm-api-model")
	}
	return nil
}

func processSource(config *Config, source string) (string, error) {
	fmt.Printf("Processing: %s\n", source)

	cached, err := IsCached(source)
	if err != nil {
		return "", err
	}

	var content string
	if cached {
		fmt.Printf("  Using cached content for %s\n", source)
		content, err = GetCachedContent(source)
		if err != nil {
			return "", err
		}
	} else {
		fmt.Printf("  Crawling %s\n", source)
		content, err = CrawlWebsite(source)
		if err != nil {
			return "", err
		}
		if err := CacheContent(source, content); err != nil {
			return "", err
		}
	}

	fmt.Printf("  Summarizing %s\n", source)
	summary, err := SummarizeWithAI(config, content, source)
	if err != nil {
		return "", err
	}

	if err := StoreSummarization(source, summary); err != nil {
		return "", err
	}

	fmt.Printf("  ✓ Completed %s\n", source)
	return summary, nil
}
