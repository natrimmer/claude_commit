# Claude Commit

A simple CLI tool that uses the Claude API to generate Git commit messages based on your staged changes.

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

### 1. Configure the tool with your Anthropic API key

```bash
claude_commit config -api-key "your-api-key" -model "claude-3-haiku-20240307"
```

Available models:

- `claude-3-7-sonnet-20250219`
- `claude-3-5-sonnet-20241022`
- `claude-3-5-haiku-20241022`
- `claude-3-opus-20240229`
- `claude-3-sonnet-20240229`
- `claude-3-haiku-20240307`

### 2. Generate a commit message based on staged changes

First, stage your changes using `git add`:

```bash
git add .
```

Then generate a commit message:

```bash
claude_commit commit
```

This will output a suggested git commit command with a generated message:

```
git commit -m "Add user authentication and password reset functionality"
```

You can then execute this command directly.

## How It Works

1. The tool reads your Anthropic API key and model preference from the config
2. When generating a commit message, it runs `git diff --staged` to get your staged changes
3. It sends this diff to the Anthropic API with a prompt to generate a concise commit message
4. The tool outputs the commit command with the generated message
