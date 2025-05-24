package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ANSI color codes
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
)

type Config struct {
	ApiKey string `json:"api_key"`
	Model  string `json:"model"`
}

type AnthropicRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func main() {
	configCmd := flag.NewFlagSet("config", flag.ExitOnError)
	apiKey := configCmd.String("api-key", "", "Anthropic API key")
	model := configCmd.String("model", "claude-3-7-sonnet-latest", "Anthropic model to use")

	commitCmd := flag.NewFlagSet("commit", flag.ExitOnError)
	viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
	modelsCmd := flag.NewFlagSet("models", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println(Bold + Magenta + "Claude Commit" + Reset)
		fmt.Println(Dim + Magenta + "Generate conventional commit messages with Anthropic's Claude" + Reset)
		fmt.Println(Dim + "Expected 'config', 'view', 'commit', or 'models' subcommands" + Reset)

		// Show usage examples
		fmt.Println("\n" + Bold + "Examples:" + Reset)
		fmt.Println("  claude_commit config -api-key \"your-api-key\" -model \"claude-3-haiku-20240307\"")
		fmt.Println("  claude_commit view")
		fmt.Println("  claude_commit models")
		fmt.Println("  claude_commit commit")

		// Show conventional commit info
		fmt.Println("\n" + Bold + "Commit Types:" + Reset)
		fmt.Println("  feat:     A new feature")
		fmt.Println("  fix:      A bug fix")
		fmt.Println("  docs:     Documentation changes")
		fmt.Println("  style:    Code style changes (formatting, etc.)")
		fmt.Println("  refactor: Code refactoring without changes to functionality")
		fmt.Println("  perf:     Performance improvements")
		fmt.Println("  test:     Adding or updating tests")
		fmt.Println("  chore:    Maintenance tasks, dependency updates, etc.")
		fmt.Println("  ci:       Continuous integration changes")
		fmt.Println("  build:    Changes that affect the build system or external dependencies")
		fmt.Println("  revert:   Reverts a previous commit")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "config":
		err := configCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(Red + fmt.Sprintf("Error parsing config arguments: %v", err) + Reset)
			os.Exit(1)
		}
		saveConfig(*apiKey, *model)
	case "view":
		err := viewCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(Red + fmt.Sprintf("Error parsing view arguments: %v", err) + Reset)
			os.Exit(1)
		}
		viewConfig()
	case "models":
		err := modelsCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(Red + fmt.Sprintf("Error parsing models arguments: %v", err) + Reset)
			os.Exit(1)
		}
		showModels()
	case "commit":
		err := commitCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(Red + fmt.Sprintf("Error parsing commit arguments: %v", err) + Reset)
			os.Exit(1)
		}
		generateCommitMessage()
	default:
		fmt.Println(Red + "Expected 'config', 'view' or 'commit' subcommands" + Reset)
		os.Exit(1)
	}
}

func saveConfig(apiKey, model string) {
	if apiKey == "" {
		fmt.Println(Red + "API key is required" + Reset)
		os.Exit(1)
	}

	config := Config{
		ApiKey: apiKey,
		Model:  model,
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(Red + fmt.Sprintf("Error getting home directory: %v", err) + Reset)
		os.Exit(1)
	}

	configDir := filepath.Join(homeDir, ".claude-commit")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		fmt.Println(Red + fmt.Sprintf("Error creating config directory: %v", err) + Reset)
		os.Exit(1)
	}

	configFile := filepath.Join(configDir, "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println(Red + fmt.Sprintf("Error marshaling config: %v", err) + Reset)
		os.Exit(1)
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		fmt.Println(Red + fmt.Sprintf("Error writing config file: %v", err) + Reset)
		os.Exit(1)
	}

	fmt.Println(Green + "Configuration saved successfully" + Reset)
	fmt.Println(Bold + "API Key: " + Reset + maskAPIKey(apiKey))
	fmt.Println(Bold + "Model: " + Reset + model)
}

func viewConfig() {
	config := loadConfig()

	fmt.Println(Bold + Cyan + "Current Configuration:" + Reset)
	fmt.Println(Bold + "API Key: " + Reset + maskAPIKey(config.ApiKey))
	fmt.Println(Bold + "Model: " + Reset + config.Model)
}

func showModels() {
	config := loadConfig()
	models := [6]string{"claude-opus-4-0", "claude-sonnet-4-0", "claude-3-7-sonnet-latest", "claude-3-5-sonnet-latest", "claude-3-5-haiku-latest", "claude-3-opus-latest"}

	fmt.Println(Bold + Cyan + "Available Models:" + Reset)
	for _, model := range models {
		switch model {
		case config.Model:
			fmt.Println(Bold + Green + model + " [CURRENT]" + Reset)
		case "claude-3-7-sonnet-latest":
			fmt.Println(Bold + model + " [DEFAULT]" + Reset)
		default:
			fmt.Println(Bold + model + Reset)
		}
	}
}

// maskAPIKey masks most of the API key for display purposes
func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "********"
	}
	return apiKey[:4] + "****" + apiKey[len(apiKey)-4:]
}

func loadConfig() Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(Red + fmt.Sprintf("Error getting home directory: %v", err) + Reset)
		os.Exit(1)
	}

	configFile := filepath.Join(homeDir, ".claude-commit", "config.json")
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println(Red + fmt.Sprintf("Error reading config file: %v\nPlease run 'config' first", err) + Reset)
		os.Exit(1)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(Red + fmt.Sprintf("Error parsing config file: %v", err) + Reset)
		os.Exit(1)
	}

	return config
}

func getGitDiff() string {
	cmd := exec.Command("git", "diff", "--staged")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(Red + fmt.Sprintf("Error running git diff: %v", err) + Reset)
		os.Exit(1)
	}
	return out.String()
}

func getFileNames() string {
	cmd := exec.Command("git", "diff", "--staged", "--name-only")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(Red + fmt.Sprintf("Error getting changed files: %v", err) + Reset)
		os.Exit(1)
	}
	return out.String()
}

func generateCommitMessage() {
	config := loadConfig()
	diff := getGitDiff()
	files := getFileNames()

	if strings.TrimSpace(diff) == "" {
		fmt.Println(Yellow + "No staged changes found. Use git add to stage changes." + Reset)
		os.Exit(1)
	}

	// Show a nice "Thinking..." message
	fmt.Println(Dim + "⚙️  Analyzing git diff with Claude AI..." + Reset)

	prompt := `Generate a conventional commit message based on the following git diff.

The message should follow this format: <type>: <description>

Types include:
- feat: A new feature
- fix: A bug fix
- docs: Documentation changes
- style: Code style changes (formatting, etc.)
- refactor: Code refactoring without changes to functionality
- perf: Performance improvements
- test: Adding or updating tests
- chore: Maintenance tasks, dependency updates, etc.
- ci: Continuous integration changes
- build: Changes that affect the build system or external dependencies
- revert: Reverts a previous commit

Guidelines:
1. Use the imperative mood ("add feature" not "Added feature")
2. All lowercase characters
3. No period at the end
4. Be concise but descriptive (what was changed and why)
5. Maximum 50 characters

Here are the files changed:
` + files + `

Here is the git diff:
` + diff

	commitMsg := callAnthropicAPI(config, prompt)
	commitMsg = strings.TrimSpace(commitMsg)

	// Format the final command nicely
	gitCommand := fmt.Sprintf("git commit -m \"%s\"", commitMsg)

	fmt.Println(Green + "✓ Commit message generated" + Reset)
	fmt.Println()
	fmt.Println(Bold + gitCommand + Reset)
}

func callAnthropicAPI(config Config, prompt string) string {
	requestBody := AnthropicRequest{
		Model: config.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens: 100,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println(Red + fmt.Sprintf("Error creating request: %v", err) + Reset)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(Red + fmt.Sprintf("Error creating request: %v", err) + Reset)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.ApiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(Red + fmt.Sprintf("Error making API call: %v", err) + Reset)
		os.Exit(1)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println(Red + fmt.Sprintf("Error closing response body: %v", err) + Reset)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(Red + fmt.Sprintf("API error (status %d): %s", resp.StatusCode, body) + Reset)
		os.Exit(1)
	}

	var anthropicResp AnthropicResponse
	err = json.NewDecoder(resp.Body).Decode(&anthropicResp)
	if err != nil {
		fmt.Println(Red + fmt.Sprintf("Error parsing API response: %v", err) + Reset)
		os.Exit(1)
	}

	if len(anthropicResp.Content) == 0 {
		fmt.Println(Red + "Empty response from API" + Reset)
		os.Exit(1)
	}

	return anthropicResp.Content[0].Text
}
