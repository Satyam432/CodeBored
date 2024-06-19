package processhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"geminiProject/gemini"

	"github.com/google/generative-ai-go/genai"
)

func codeDeveloper(input string, stackToUse string, approachToUse string, databaseToUse string, projectStructure string, baseLogic string) (string, error) {
	// Create the prompt
	prompt := fmt.Sprintf(`Generate code for a project using the following parameters:
		- Input: %s
		- Stack: %s
		- Approach: %s
		- Database: %s
		- Project Structure: %s
		- Base Logic: %s

		Return the code in JSON format where each key is the file path and the value is the code for that file.`, input, stackToUse, approachToUse, databaseToUse, projectStructure, baseLogic)

	ctx := context.Background()
	clientGemini := gemini.GetGeminiCLient()

	// Generate content using the Gemini client
	resp, err := clientGemini.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("error generating content: %v", err)
	}

	// Debugging: Print the raw response
	marshalResponse, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println("Raw Response:", string(marshalResponse))

	var generateResponse ContentResponse
	if err := json.Unmarshal(marshalResponse, &generateResponse); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Extract the project structure from the response
	var codeDeveloped string
	for _, candidate := range *generateResponse.Candidates {
		if candidate.Content != nil && len(candidate.Content.Parts) > 0 {
			codeDeveloped = candidate.Content.Parts[0]
			break
		}
	}

	// Debugging: Print the raw codeDeveloped response
	fmt.Println("Raw codeDeveloped:", codeDeveloped)

	return codeDeveloped, nil
}
