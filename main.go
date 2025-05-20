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

	"github.com/charmbracelet/lipgloss"
)

// Lipgloss styles
var (
	// Colors
	purple      = lipgloss.Color("#7D56F4")
	lightPurple = lipgloss.Color("#9E83F5")
	gray        = lipgloss.Color("#888888")
	green       = lipgloss.Color("#2ECC71")
	yellow      = lipgloss.Color("#F1C40F")
	red         = lipgloss.Color("#E74C3C")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(purple).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lightPurple).
			MarginBottom(1)

	successStyle = lipgloss.NewStyle().
			Foreground(green)

	errorStyle = lipgloss.NewStyle().
			Foreground(red)

	warningStyle = lipgloss.NewStyle().
			Foreground(yellow)

	infoStyle = lipgloss.NewStyle().
			Foreground(gray)

	commandStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple).
			Padding(1, 2).
			MarginTop(1)
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
	model := configCmd.String("model", "claude-3-haiku-20240307", "Anthropic model to use")

	commitCmd := flag.NewFlagSet("commit", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println(titleStyle.Render("Claude Commit"))
		fmt.Println(subtitleStyle.Render("Generate commit messages with Claude AI"))
		fmt.Println(infoStyle.Render("Expected 'config' or 'commit' subcommands"))

		// Show usage examples in a nice box
		usageText := "Examples:\n" +
			"  claude_commit config -api-key \"your-api-key\" -model \"claude-3-haiku-20240307\"\n" +
			"  claude_commit commit"
		fmt.Println(boxStyle.Render(usageText))
		os.Exit(1)
	}

	switch os.Args[1] {
	case "config":
		err := configCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("Error parsing config arguments: %v", err)))
			os.Exit(1)
		}
		saveConfig(*apiKey, *model)
	case "commit":
		err := commitCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("Error parsing commit arguments: %v", err)))
			os.Exit(1)
		}
		generateCommitMessage()
	default:
		fmt.Println(errorStyle.Render("Expected 'config' or 'commit' subcommands"))
		os.Exit(1)
	}
}

func saveConfig(apiKey, model string) {
	if apiKey == "" {
		fmt.Println(errorStyle.Render("API key is required"))
		os.Exit(1)
	}

	config := Config{
		ApiKey: apiKey,
		Model:  model,
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error getting home directory: %v", err)))
		os.Exit(1)
	}

	configDir := filepath.Join(homeDir, ".claude-commit")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error creating config directory: %v", err)))
		os.Exit(1)
	}

	configFile := filepath.Join(configDir, "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error marshaling config: %v", err)))
		os.Exit(1)
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error writing config file: %v", err)))
		os.Exit(1)
	}

	// Create a nice success message with config details
	configDetails := fmt.Sprintf("API Key: %s\nModel: %s", maskAPIKey(apiKey), model)
	configBox := boxStyle.Render(configDetails)

	fmt.Println(successStyle.Render("Configuration saved successfully"))
	fmt.Println(configBox)
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
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error getting home directory: %v", err)))
		os.Exit(1)
	}

	configFile := filepath.Join(homeDir, ".claude-commit", "config.json")
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error reading config file: %v\nPlease run 'config' first", err)))
		os.Exit(1)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error parsing config file: %v", err)))
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
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error running git diff: %v", err)))
		os.Exit(1)
	}
	return out.String()
}

func generateCommitMessage() {
	config := loadConfig()
	diff := getGitDiff()

	if strings.TrimSpace(diff) == "" {
		fmt.Println(warningStyle.Render("No staged changes found. Use git add to stage changes."))
		os.Exit(1)
	}

	// Show a nice "Thinking..." message
	fmt.Println(infoStyle.Render("⚙️  Analyzing git diff with Claude AI..."))

	prompt := "Generate a concise and descriptive git commit message based on the following git diff. Focus on what was changed and why, and limit the message to a single line less than 72 characters:\n\n" + diff

	commitMsg := callAnthropicAPI(config, prompt)
	commitMsg = strings.TrimSpace(commitMsg)

	// Format the final command nicely
	gitCommand := fmt.Sprintf("git commit -m \"%s\"", commitMsg)

	fmt.Println(successStyle.Render("✓ Commit message generated"))
	fmt.Println(boxStyle.Render(commandStyle.Render(gitCommand)))
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
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error creating request: %v", err)))
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error creating request: %v", err)))
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.ApiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error making API call: %v", err)))
		os.Exit(1)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("Error closing response body: %v", err)))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(errorStyle.Render(fmt.Sprintf("API error (status %d): %s", resp.StatusCode, body)))
		os.Exit(1)
	}

	var anthropicResp AnthropicResponse
	err = json.NewDecoder(resp.Body).Decode(&anthropicResp)
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error parsing API response: %v", err)))
		os.Exit(1)
	}

	if len(anthropicResp.Content) == 0 {
		fmt.Println(errorStyle.Render("Empty response from API"))
		os.Exit(1)
	}

	return anthropicResp.Content[0].Text
}
