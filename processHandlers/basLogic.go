package processhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"geminiProject/gemini"

	"github.com/google/generative-ai-go/genai"
)

func BaseLogic(input string, stackToUse string, approachToUse string, databaseToUse string, projectStructure string) (string, error) {
	prompt := fmt.Sprintf("Given the approach: %s,\n", approachToUse)
	prompt += fmt.Sprintf("Stack to use:\n```json\n%s```,\n", stackToUse)
	prompt += fmt.Sprintf("Database to use: %s\n", databaseToUse)
	prompt += fmt.Sprintf(", please generate the *well described logic* for each file in the project structure :\n%s", projectStructure)
	fmt.Println("Ultimate prompt:", prompt)
	ctx := context.Background()
	clientGemini := gemini.GetGeminiCLient()
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
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}
	// Extract the project structure from the response
	var baseLogic string
	for _, candidate := range *generateResponse.Candidates {
		if candidate.Content != nil && len(candidate.Content.Parts) > 0 {
			part := candidate.Content.Parts[0]
			baseLogic = part
			break
		}
	}
	return baseLogic, nil
}
