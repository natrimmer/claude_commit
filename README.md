# Claude Commit

A simple CLI tool that uses the Claude API to generate Git commit messages from staged changes.

## Installation

```bash
go install github.com/natrimmer/claude_commit@latest
```

Or build from source:

```bash
git clone https://github.com/natrimmer/claude_commit.git
cd claude_commit
go build
```

## Usage

### 1. Configure with your API key

```bash
claude_commit config -api-key "your-api-key" -model "claude-3-haiku-20240307"
```

Available models:

- `claude-3-7-sonnet-20250219`
- `claude-3-5-sonnet-20241022`
- `claude-3-opus-20240229`
- `claude-3-haiku-20240307`

### 2. View current configuration

```bash
claude_commit view
```

Output:

```
Current Configuration:
┌─────────────────────────────────┐
│ API Key: abcd****wxyz           │
│ Model: claude-3-haiku-20240307  │
└─────────────────────────────────┘
```

### 3. Generate a commit message

```bash
git add .                # Stage your changes
claude_commit commit     # Generate a commit message
```

Output:

```
✓ Commit message generated
┌──────────────────────────────────────────────────────────────────────────┐
│ git commit -m "Add user authentication and password reset functionality" │
└──────────────────────────────────────────────────────────────────────────┘
```

## How It Works

1. Reads your Anthropic API key from config
2. Gets staged changes with `git diff --staged`
3. Sends the diff to Claude API
4. Returns a formatted git commit command

## Features

- Zero dependencies - pure Go standard library
- Simple configuration with view option
- API key stored securely with masking for display
- Colorized terminal output
