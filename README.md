# nub

A simple, fast CLI tool for crawling websites, summarizing them with AI, and displaying results.

## Features

- Crawl and cache websites
- AI-powered summarization (OpenAI-compatible APIs)
- Smart caching (24-hour cache validity)
- Daemon mode for scheduled crawling
- Minimalist HTML viewer (HackerNews-inspired design)
- Full markdown support
- Logs viewer with pager support
- Fast, simple, zero bloat

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
  "summary_prompt": "Summarize the key news topics and main stories from this website. Focus on the most important headlines and provide a concise overview in markdown format."
}
```

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

# View summaries in browser
nub --show

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
5. **Store**: Saves summary as markdown in `~/.local/nub/summaries/`
6. **Display**: View all summaries as formatted HTML with enhanced markdown rendering

### Default Behavior

Running `nub` without arguments shows the help menu. Use `nub --run` to crawl and summarize sites, or `nub --show` to view existing summaries.

## Markdown Rendering

The HTML viewer features a comprehensive markdown renderer with minimalist HackerNews-inspired design:

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
- Orange header bar (#ff6600) for visual anchor
- Mobile responsive with optimized spacing
- Fast loading, zero bloat

## Data Storage

- **Config**: `~/.config/nub/config.json` (preserved by `--clear-data`)
- **Cache**: `~/.local/nub/cache/` (HTML content)
- **Summaries**: `~/.local/nub/summaries/` (Markdown files)
- **Logs**: `~/.local/nub/nub.log` (Daemon logs)
- **PID File**: `~/.local/nub/nub.pid` (Daemon process ID)
- **HTML View**: `~/.local/nub/view.html`

### Clearing Data

- `--clear-cache`: Removes only cached website content
- `--clear-data`: Removes **all** data in `~/.local/nub/` (cache, summaries, logs, PID)
- Note: Config file in `~/.config/nub/` is **not** affected by `--clear-data`

## Supported LLM Providers

Any OpenAI-compatible API:
- OpenAI
- Mistral AI
- Anthropic (via compatible endpoint)
- Local LLMs (Ollama, LM Studio, etc.)

## License

MIT
