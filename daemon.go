package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func GetLogPath() (string, error) {
	dataDir, err := GetDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "nub.log"), nil
}

func GetPidPath() (string, error) {
	dataDir, err := GetDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "nub.pid"), nil
}

func StartDaemon(config *Config) error {
	logPath, err := GetLogPath()
	if err != nil {
		return err
	}

	pidPath, err := GetPidPath()
	if err != nil {
		return err
	}

	if isRunning, pid := IsDaemonRunning(); isRunning {
		return fmt.Errorf("daemon already running (PID: %d)", pid)
	}

	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	attr := &os.ProcAttr{
		Files: []*os.File{nil, logFile, logFile},
		Sys: &syscall.SysProcAttr{
			Setsid: true,
		},
	}

	process, err := os.StartProcess(execPath, []string{execPath, "--daemon-child"}, attr)
	if err != nil {
		logFile.Close()
		return fmt.Errorf("failed to start daemon: %v", err)
	}

	if err := os.WriteFile(pidPath, []byte(fmt.Sprintf("%d", process.Pid)), 0644); err != nil {
		logFile.Close()
		return fmt.Errorf("failed to write PID file: %v", err)
	}

	logFile.Close()
	process.Release()

	fmt.Printf("Daemon started successfully (PID: %d)\n", process.Pid)
	fmt.Printf("Logs: %s\n", logPath)
	fmt.Printf("Schedule: every %d minutes\n", config.ScheduleMinutes)

	return nil
}

func IsDaemonRunning() (bool, int) {
	pidPath, err := GetPidPath()
	if err != nil {
		return false, 0
	}

	data, err := os.ReadFile(pidPath)
	if err != nil {
		return false, 0
	}

	var pid int
	fmt.Sscanf(string(data), "%d", &pid)
	if pid <= 0 {
		return false, 0
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false, 0
	}

	err = process.Signal(syscall.Signal(0))
	if err != nil {
		os.Remove(pidPath)
		return false, 0
	}

	return true, pid
}

func StopDaemon() error {
	isRunning, pid := IsDaemonRunning()
	if !isRunning {
		return fmt.Errorf("daemon is not running")
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %v", err)
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to stop daemon: %v", err)
	}

	pidPath, _ := GetPidPath()
	os.Remove(pidPath)

	fmt.Printf("Daemon stopped (PID: %d)\n", pid)
	return nil
}

func RunDaemonChild(config *Config) {
	logPath, err := GetLogPath()
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags)

	log.Printf("Daemon started, crawling every %d minutes\n", config.ScheduleMinutes)

	runDaemonLoop(config)
}

func runDaemonLoop(config *Config) {
	ticker := time.NewTicker(time.Duration(config.ScheduleMinutes) * time.Minute)
	defer ticker.Stop()

	if err := runOnceLogged(config); err != nil {
		log.Printf("Error in initial run: %v\n", err)
	}

	for range ticker.C {
		log.Printf("Starting scheduled crawl...\n")
		if err := runOnceLogged(config); err != nil {
			log.Printf("Error: %v\n", err)
		}
	}
}

func runOnceLogged(config *Config) error {
	if err := validateConfig(config); err != nil {
		return err
	}

	log.Println("Starting crawl and summarization...")
	for _, source := range config.Sources {
		if err := processSourceLogged(config, source); err != nil {
			log.Printf("Error processing %s: %v\n", source, err)
			continue
		}
	}
	log.Println("Done!")
	return nil
}

func processSourceLogged(config *Config, source string) error {
	log.Printf("Processing: %s\n", source)

	cached, err := IsCached(source)
	if err != nil {
		return err
	}

	var content string
	if cached {
		log.Printf("  Using cached content for %s\n", source)
		content, err = GetCachedContent(source)
		if err != nil {
			return err
		}
	} else {
		log.Printf("  Crawling %s\n", source)
		content, err = CrawlWebsite(source)
		if err != nil {
			return err
		}
		if err := CacheContent(source, content); err != nil {
			return err
		}
	}

	log.Printf("  Summarizing %s\n", source)
	summary, err := SummarizeWithAI(config, content, source)
	if err != nil {
		return err
	}

	if err := StoreSummarization(source, summary); err != nil {
		return err
	}

	if config.FocusTopics != "" {
		focused, err := ExtractFocusedContent(config, summary)
		if err != nil {
			log.Printf("  Warning: failed to extract focused content: %v\n", err)
		} else if focused != "" && focused != "No relevant content found." {
			if err := StoreFocusedContent(source, focused); err != nil {
				log.Printf("  Warning: failed to store focused content: %v\n", err)
			}
		}
	}

	log.Printf("  ✓ Completed %s\n", source)
	return nil
}

func ShowLogs() error {
	logPath, err := GetLogPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		fmt.Println("No logs found")
		return nil
	}

	pager := "less"
	if os.Getenv("PAGER") != "" {
		pager = os.Getenv("PAGER")
	}

	exec := &execCmd{name: pager, args: []string{"+G", logPath}}
	return exec.runWait()
}
