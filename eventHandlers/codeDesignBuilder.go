package eventhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"geminiProject/gemini"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

type Content struct {
	Parts []string `json:Parts`
	Role  string   `json:Role`
}
type Candidates struct {
	Content *Content `json:Content`
}

type ContentResponse struct {
	Candidates *[]Candidates `json:Candidates`
}

func CodeDesigner(input string, stackToUse string, approachToUse string, databaseToUse string) (string, error) {
	clientGemini := gemini.GetGeminiCLient()
	// Construct the prompt
	prompt := fmt.Sprintf(
		`write just the **project structure** as a JSON for "%s".
		Approach to use: %s
		Stack to use: %s
		Database to use: %s

		keep it Modular, Clean Code, Easy to Read, Easy to Customize, Basic Security`,
		input, approachToUse, stackToUse, databaseToUse)

	ctx := context.Background()

	// Generate content using the Gemini client
	resp, err := clientGemini.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("error generating content: %v", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates found")
	}
	marshalResponse, _ := json.MarshalIndent(resp, "", "  ")

	var generateResponse ContentResponse
	if err := json.Unmarshal(marshalResponse, &generateResponse); err != nil {
		return "", err
	}

	// Extract the project structure from the response
	var projectStructure string
	for _, candidate := range *generateResponse.Candidates {
		if candidate.Content != nil && len(candidate.Content.Parts) > 0 {
			part := candidate.Content.Parts[0]
			projectStructure = part
			break
		}
	}

	if projectStructure == "" {
		return "", fmt.Errorf("no suitable project structure found")
	}
	projectStructure = strings.TrimPrefix(projectStructure, "```json\n")
	projectStructure = strings.TrimSuffix(projectStructure, "\n```")

	return projectStructure, nil
}
