# nub

A simple, fast CLI tool for crawling websites, summarizing them with AI, and displaying results in both plain text and HTML formats.

## Features

- **Smart Crawling**: Crawl and cache websites with 24-hour cache validity
- **AI-Powered Summarization**: Uses OpenAI-compatible APIs for intelligent content summarization
- **Focus Topics**: Filter content to show only what matters to you
- **Dual Display Modes**: 
  - Clean ASCII text in pager (optimized for terminal reading)
  - Rich HTML in browser (HackerNews-inspired design)
- **Daemon Mode**: Scheduled background crawling with configurable intervals
- **Full Markdown Support**: Complete markdown rendering with syntax highlighting
- **Logs Viewer**: Built-in pager support for daemon logs
- **Fast & Simple**: Zero bloat, pure Go implementation

## Installation

### Prerequisites

- Go 1.21 or later
- Git

### macOS

```bash
# Clone the repository
git clone https://github.com/yourusername/nub.git
cd nub

# Build and install
go build -o nub
sudo mv nub /usr/local/bin/

# Verify installation
nub --help
```

### Linux

```bash
# Clone the repository
git clone https://github.com/yourusername/nub.git
cd nub

# Build and install
go build -o nub
sudo mv nub /usr/local/bin/

# Verify installation
nub --help
```

**Alternative: Install to user directory (no sudo required)**

```bash
# Build
go build -o nub

# Install to ~/.local/bin (ensure this is in your PATH)
mkdir -p ~/.local/bin
mv nub ~/.local/bin/

# Add to PATH if needed (add to ~/.bashrc or ~/.zshrc)
export PATH="$HOME/.local/bin:$PATH"
```

## Configuration

Config file location: `~/.config/nub/config.json`

Example configuration:
```json
{
  "sources": [
    "https://example.com",
    "https://news.ycombinator.com"
  ],
  "llm_api_key": "your-api-key-here",
  "llm_api_url": "https://api.mistral.ai/v1/chat/completions",
  "llm_api_model": "mistral-small-latest",
  "schedule_minutes": 15,
  "focus_topics": "go,javascript,rust",
  "summary_prompt": "Summarize the key news topics and main stories from this website. Focus on the most important headlines and provide a concise overview in markdown format."
}
```

### Configuration Fields

- **sources**: Array of URLs to crawl and summarize
- **llm_api_key**: API key for your LLM provider
- **llm_api_url**: API endpoint URL (OpenAI-compatible)
- **llm_api_model**: Model name to use for summarization
- **schedule_minutes**: Interval for daemon mode (default: 15)
- **focus_topics**: Comma-separated topics to filter content (optional)
- **summary_prompt**: Custom prompt for AI summarization (optional)

## Usage

### Initial Setup

```bash
# Set LLM configuration
nub --set-llm-api-key sk-your-api-key
nub --set-llm-api-url https://api.mistral.ai/v1/chat/completions
nub --set-llm-api-model mistral-small-latest

# Add sources to crawl
nub --add-source https://news.ycombinator.com
nub --add-source https://example.com

# Customize the summarization prompt (optional)
nub --set-prompt "Summarize the key news topics in 5 bullet points"

# Set focus topics to filter content (optional)
nub --set-focus "go,javascript,rust"

# Set schedule time (for daemon mode)
nub --set-schedule-time 15
```

### Focus Topics

Set focus topics to filter and highlight only the content you care about:

```bash
# Set focus topics (comma-separated)
nub --set-focus "go,javascript,ai,kubernetes"

# Clear focus (show all content)
nub --set-focus ""
```

When focus is set, nub will:
- Use AI to extract only content related to your topics
- Display focused content at the top of the page in a highlighted section
- Keep full summaries available below

### Managing Sources

```bash
# List all sources with IDs
nub --list

# Add a new source
nub --add-source https://github.com/trending

# Remove source by ID
nub --rem-source 2

# Remove source by URL
nub --rem-source https://example.com
```

### Running

```bash
# Show help (default when no arguments)
nub

# Run crawl and summarization once
nub --run

# View summaries in terminal (plain text in pager)
nub --show

# View summaries in browser (rich HTML)
nub --show-html

# Run in daemon mode (detaches and runs in background)
nub -d

# Stop running daemon
nub --stop

# View logs from daemon
nub --logs

# Clear cached websites (forces re-crawl next time)
nub --clear-cache

# Clear all stored data (cache, summaries, logs, PID)
nub --clear-data
```

### Daemon Mode

When running in daemon mode with `nub -d`:
- The process detaches from terminal and runs in background
- Only one daemon instance can run at a time
- Output is written to `~/.local/nub/nub.log`
- Process ID is saved to `~/.local/nub/nub.pid`
- To view logs: `nub --logs` (opens in your pager)
- To stop: `nub --stop`

## How It Works

### Workflow

1. **Load Config**: Reads configuration from `~/.config/nub/config.json`
2. **Check Cache**: Checks if website is cached (24-hour validity, or until cleared)
3. **Crawl**: If not cached, fetches website content
4. **Summarize**: Uses OpenAI-compatible API to generate summary
5. **Focus (Optional)**: Extracts only content matching your focus topics
6. **Store**: Saves summaries as markdown in `~/.local/nub/summaries/`
7. **Display**: View as plain text (`--show`) or HTML (`--show-html`)

### Display Modes

**Plain Text Mode (`--show`)**
- Strips all markdown formatting for clean ASCII reading
- Optimized text wrapping at 78 characters
- Visual separators between summaries
- Opens in your configured pager (less, more, etc.)
- Perfect for terminal-only workflows

**HTML Mode (`--show-html`)**
- Rich markdown rendering with full formatting
- HackerNews-inspired minimalist design
- Syntax highlighting for code blocks
- Opens in your default browser
- Mobile responsive

### Default Behavior

Running `nub` without arguments shows the help menu. Use `nub --run` to crawl and summarize sites, `nub --show` for terminal view, or `nub --show-html` for browser view.

## Markdown Rendering

### HTML Mode Features

The HTML viewer features comprehensive markdown rendering with a minimalist HackerNews-inspired design:

**Supported Markdown:**
- **Formatting**: Bold (`**text**`), italic (`*text*`), inline code (`` `code` ``)
- **Headers**: H1-H6 (`#` to `######`)
- **Lists**: Unordered (`-`, `*`, `+`) and ordered (`1.`, `2.`, etc.)
- **Links**: `[text](url)` - opens in new tab
- **Code Blocks**: Triple backticks with language syntax highlighting
- **Blockquotes**: `> quote text`
- **Horizontal Rules**: `---`, `***`, or `___`

**Design:**
- Clean, minimal aesthetic inspired by HackerNews
- Verdana font (10pt) for maximum readability
- Beige background (#f6f6ef) easy on the eyes
- Pink/rose header bar (#dc94ba) for visual anchor
- Mobile responsive with optimized spacing
- Fast loading, zero JavaScript

### Plain Text Mode Features

The plain text viewer provides clean, terminal-friendly output:

- **Markdown Stripping**: Removes all markdown syntax for pure ASCII reading
- **Smart Text Wrapping**: Wraps at 78 characters for comfortable terminal viewing
- **Visual Separators**: Uses Unicode box drawing characters for clear section breaks
- **Focus Highlighting**: Distinguished headers for focused content sections
- **Preserved Structure**: Maintains paragraph breaks and logical content flow
- **Pager Integration**: Works seamlessly with less, more, or your configured $PAGER

## Data Storage

- **Config**: `~/.config/nub/config.json` (preserved by `--clear-data`)
- **Cache**: `~/.local/nub/cache/` (HTML content from websites)
- **Summaries**: `~/.local/nub/summaries/` (AI-generated markdown summaries)
- **Focus**: `~/.local/nub/focus/` (Filtered content based on topics)
- **Logs**: `~/.local/nub/nub.log` (Daemon operation logs)
- **PID File**: `~/.local/nub/nub.pid` (Daemon process tracking)
- **View Files**: 
  - `~/.local/nub/view.md` (Plain text view for pager)
  - `~/.local/nub/view.html` (HTML view for browser)

### Clearing Data

- `--clear-cache`: Removes only cached website content (forces fresh crawls)
- `--clear-data`: Removes **all** data in `~/.local/nub/` (cache, summaries, focus, logs, PID)
- Note: Config file in `~/.config/nub/` is **not** affected by `--clear-data`

## Supported LLM Providers

Any OpenAI-compatible API:
- OpenAI
- Mistral AI
- Anthropic (via compatible endpoint)
- Local LLMs (Ollama, LM Studio, etc.)
- Any service implementing the OpenAI chat completions API

## Quick Reference

```bash
# Setup
nub --set-llm-api-key <key>          # Set API key
nub --set-llm-api-url <url>          # Set API endpoint
nub --set-llm-api-model <model>      # Set model name
nub --add-source <url>               # Add website to track

# Running
nub --run                            # Crawl and summarize once
nub -d                               # Start daemon (background)
nub --stop                           # Stop daemon

# Viewing
nub --show                           # View in terminal (plain text)
nub --show-html                      # View in browser (HTML)
nub --logs                           # View daemon logs

# Managing
nub --list                           # List all sources
nub --rem-source <id>                # Remove source
nub --clear-cache                    # Clear cached websites
nub --clear-data                     # Clear all data

# Optional
nub --set-focus <topics>             # Filter by topics
nub --set-prompt <text>              # Custom prompt
nub --set-schedule-time <mins>       # Set daemon interval
```

## Tips

- Use `--show` for quick terminal checks, `--show-html` for detailed browsing
- Set focus topics to reduce noise and see only what matters
- Daemon mode is perfect for morning news digests
- Cache is valid for 24 hours - use `--clear-cache` to force fresh content
- Customize the summary prompt to match your reading style
- View files are temporary and regenerated on each `--show` call

## License

MIT
