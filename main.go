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
		fmt.Println("Expected 'config' or 'commit' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "config":
		err := configCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Printf("Error parsing config arguments: %v\n", err)
			os.Exit(1)
		}
		saveConfig(*apiKey, *model)
	case "commit":
		err := commitCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Printf("Error parsing commit arguments: %v\n", err)
			os.Exit(1)
		}
		generateCommitMessage()
	default:
		fmt.Println("Expected 'config' or 'commit' subcommands")
		os.Exit(1)
	}
}

func saveConfig(apiKey, model string) {
	if apiKey == "" {
		fmt.Println("API key is required")
		os.Exit(1)
	}

	config := Config{
		ApiKey: apiKey,
		Model:  model,
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	configDir := filepath.Join(homeDir, ".claude-commit")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		fmt.Printf("Error creating config directory: %v\n", err)
		os.Exit(1)
	}

	configFile := filepath.Join(configDir, "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling config: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Configuration saved successfully")
}

func loadConfig() Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	configFile := filepath.Join(homeDir, ".claude-commit", "config.json")
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %v\nPlease run 'config' first\n", err)
		os.Exit(1)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
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
		fmt.Printf("Error running git diff: %v\n", err)
		os.Exit(1)
	}
	return out.String()
}

func generateCommitMessage() {
	config := loadConfig()
	diff := getGitDiff()

	if strings.TrimSpace(diff) == "" {
		fmt.Println("No staged changes found. Use git add to stage changes.")
		os.Exit(1)
	}

	prompt := "Generate a concise and descriptive git commit message based on the following git diff. Focus on what was changed and why, and limit the message to a single line less than 72 characters:\n\n" + diff

	commitMsg := callAnthropicAPI(config, prompt)
	commitMsg = strings.TrimSpace(commitMsg)

	fmt.Printf("git commit -m \"%s\"\n", commitMsg)
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
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.ApiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making API call: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("API error (status %d): %s\n", resp.StatusCode, body)
		os.Exit(1)
	}

	var anthropicResp AnthropicResponse
	err = json.NewDecoder(resp.Body).Decode(&anthropicResp)
	if err != nil {
		fmt.Printf("Error parsing API response: %v\n", err)
		os.Exit(1)
	}

	if len(anthropicResp.Content) == 0 {
		fmt.Println("Empty response from API")
		os.Exit(1)
	}

	return anthropicResp.Content[0].Text
}
