# Claude Commit

A simple CLI tool that uses the Claude API to generate Git commit messages from staged changes, following conventional commit best practices.

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
build  # or: go build
```

## Quick Start

```bash
# Get help
claude_commit
# or
claude_commit --help

# Check version
claude_commit --version

# Configure
claude_commit config -api-key "your-api-key"

# Generate commit message
git add .
claude_commit commit
```

## Commands

### Help and Version

```bash
claude_commit              # Show help
claude_commit --help       # Show help
claude_commit help         # Show help
claude_commit --version    # Show version info
```

### Configuration

```bash
# Configure with your API key (uses claude-3-7-sonnet-latest by default)
claude_commit config -api-key "your-api-key"

# Configure with specific model
claude_commit config -api-key "your-api-key" -model "claude-3-5-sonnet-latest"

# View current configuration
claude_commit view

# List available models
claude_commit models
```

### Generate Commit Messages

```bash
git add .                # Stage your changes
claude_commit commit     # Generate a commit message
```

## Available Models

- `claude-opus-4-0` - Most capable, slower and more expensive
- `claude-sonnet-4-0` - Balanced performance and speed  
- `claude-3-7-sonnet-latest` - **Default** - Fast and efficient
- `claude-3-5-sonnet-latest` - Previous generation, reliable
- `claude-3-5-haiku-latest` - Fastest and most cost-effective
- `claude-3-opus-latest` - Previous generation, most capable

## Example Usage

### Configuration

```bash
$ claude_commit config -api-key "sk-ant-api03-..." -model "claude-3-7-sonnet-latest"
Configuration saved successfully
API Key: sk-a****...
Model: claude-3-7-sonnet-latest

$ claude_commit view
Current Configuration:
API Key: sk-a****...
Model: claude-3-7-sonnet-latest

$ claude_commit models
Available Models:
claude-opus-4-0
claude-sonnet-4-0
claude-3-7-sonnet-latest [CURRENT]
claude-3-5-sonnet-latest
claude-3-5-haiku-latest
claude-3-opus-latest
```

### Generating Commits

```bash
$ git add .
$ claude_commit commit
⚙️  Analyzing git diff with Claude AI...
✓ Commit message generated

git commit -m "feat: add user authentication and password reset functionality"
```

### Version Information

```bash
$ claude_commit --version
Claude Commit v1.2.3
Build Date: 2024-01-15T10:30:00Z
Commit: abc1234
Generate conventional commit messages with Anthropic's Claude
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

## Configuration Storage

Your configuration is stored in a JSON file at `~/.claude-commit/config.json`. The API key is stored in plaintext, so ensure appropriate file permissions are set.

## Features

- Zero dependencies
- Follows conventional commit best practices
- Uses conventional commit format
- Configuration stored in `~/.claude-commit/config.json`
- API key masking for display security
- Colorized terminal output
- Version information with build details
- Comprehensive help system
- Multiple model support with easy switching
- Clean error handling and user feedback

## Development

### Building from Source

```bash
git clone https://github.com/natrimmer/claude_commit.git
cd claude_commit

# The devenv environment provides all necessary tools
# Install dependencies are handled automatically by devenv

# Run tests
test

# Build with version info
build

# Build release version
build-release

# Run all quality checks
ci
```

### Available Commands

When you enter the devenv shell, you'll have access to these commands:

- `build` - Build with version info
- `build-release` - Build optimized release binary
- `test` - Run tests
- `test-coverage` - Run tests with coverage
- `test-race` - Run tests with race detection
- `bench` - Run benchmark tests
- `lint` - Run linter
- `fmt` - Format code
- `vet` - Run go vet
- `version` - Show version information
- `clean` - Clean build artifacts
- `ci` - Run all CI checks

### Version Management

This project uses **Semantic Versioning (SemVer)**. Versions are managed through git tags:

```bash
# Create a new version tag
git tag v1.2.3
git push origin v1.2.3

# Build will automatically use the tag
build
./claude_commit --version  # Shows: Claude Commit v1.2.3
```

### Development Workflow

```bash
# Enter the development environment
cd claude_commit  # devenv activates automatically with direnv

# Make changes, then test
fmt      # Format code
lint     # Check for issues
test     # Run tests
ci       # Run full CI suite

# Build and test
build
./claude_commit --version
```
