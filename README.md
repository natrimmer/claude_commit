# Claude Commit

A simple CLI tool that uses the Claude API to generate Git commit messages from staged changes, following best practices.

## Installation

### Option 1: Download binary

Download the pre-built binary for your platform from the [GitHub Releases page](https://github.com/natrimmer/claude_commit/releases/latest):

```bash
# Example for Linux (amd64)
curl -L https://github.com/natrimmer/claude_commit/releases/latest/download/claude_commit_linux_amd64 -o claude_commit
chmod +x claude_commit
sudo mv claude_commit /usr/local/bin/

# Example for macOS (intel)
curl -L https://github.com/natrimmer/claude_commit/releases/latest/download/claude_commit_darwin_amd64 -o claude_commit
chmod +x claude_commit
sudo mv claude_commit /usr/local/bin/

# Example for macOS (Apple Silicon)
curl -L https://github.com/natrimmer/claude_commit/releases/latest/download/claude_commit_darwin_arm64 -o claude_commit
chmod +x claude_commit
sudo mv claude_commit /usr/local/bin/
```

### Option 2: Using Go

```bash
go install github.com/natrimmer/claude_commit@latest
```

### Option 3: Build from source

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
API Key: abcd****wxyz
Model: claude-3-haiku-20240307
```

### Configuration Storage

Your configuration is stored in a JSON file at `~/.claude-commit/config.json`. The API key is stored in plaintext, so ensure appropriate file permissions are set.

### 3. Generate a commit message

```bash
git add .                # Stage your changes
claude_commit commit     # Generate a commit message
```

Output:

```
⚙️  Analyzing git diff with Claude AI...
✓ Commit message generated
git commit -m "feat: add user authentication and password reset functionality"
```

## Commit Message Format

- Type prefix (feat, fix, docs, etc.)
- Lowercase throughout
- Imperative mood
- No period at end
- Format: `<type>: <description>`

## Conventional Commit Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring without functionality changes
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks, dependency updates, etc.
- `ci`: Continuous integration changes
- `build`: Changes that affect the build system or external dependencies
- `revert`: Reverts a previous commit

## How It Works

1. Reads your Anthropic API key from config (stored in `~/.claude-commit/config.json`)
2. Gets staged changes with `git diff --staged`
3. Sends the diff and detailed prompt to Claude API
4. Returns a formatted git commit command

## Features

- Zero dependencies
- Follows commit message best practices
- Uses conventional commit format
- Configuration stored in `~/.claude-commit/config.json`
- API key stored with masking for display
- Colorized terminal output
