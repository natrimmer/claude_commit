# Claude Commit

A simple CLI tool that uses the Claude API to generate Git commit messages from staged changes, following best practices.

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
# Standard commit format
claude_commit config -api-key "your-api-key" -model "claude-3-haiku-20240307"

# Use conventional commit format
claude_commit config -api-key "your-api-key" -model "claude-3-haiku-20240307" -conventional
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
Conventional Commits: true
```

### 3. Generate a commit message

```bash
git add .                # Stage your changes
claude_commit commit     # Generate a commit message
```

Output (standard format):
```
⚙️  Analyzing git diff with Claude AI...
✓ Commit message generated
git commit -m "Add user authentication and password reset functionality"
```

Output (conventional format):
```
⚙️  Analyzing git diff with Claude AI...
✓ Commit message generated
git commit -m "feat: add user authentication and password reset functionality"
```

## Commit Message Best Practices

The tool enforces two message formats:

### Standard Format
- Capitalized first word
- Imperative mood ("Add feature" not "Added feature")
- No period at end
- Descriptive and concise
- Less than 50 characters when possible

### Conventional Commit Format
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

## How It Works

1. Reads your Anthropic API key from config
2. Gets staged changes with `git diff --staged`
3. Sends the diff and detailed prompt to Claude API
4. Returns a formatted git commit command

## Features

- Zero dependencies - pure Go standard library
- Follows commit message best practices
- Optional conventional commit format
- API key stored securely with masking for display
- Colorized terminal output
