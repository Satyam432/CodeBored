package processhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"geminiProject/gemini"
	"io/ioutil"

	"github.com/google/generative-ai-go/genai"
)

func codeDeveloper(input string, stackToUse string, approachToUse string, databaseToUse string, projectStructure string, baseLogic string) (string, error) {
	//Prompt to iterate over projectStructure and projectStructure to generate code
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
	// Ensure the response is in JSON format
	marshalResponse, _ := json.MarshalIndent(resp, "", "  ")
	var generateResponse ContentResponse
	if err := json.Unmarshal(marshalResponse, &generateResponse); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}
	// Extract the project structure from the response
	var codeDeveloped string
	for _, candidate := range *generateResponse.Candidates {
		if candidate.Content != nil && len(candidate.Content.Parts) > 0 {
			part := candidate.Content.Parts[0]
			codeDeveloped = part
			break
		}
	}

	// Write the generated code to a Python file
	fileName := "generated_code.py"
	err = ioutil.WriteFile(fileName, []byte(codeDeveloped), 0644)
	if err != nil {
		return "", fmt.Errorf("error writing to file: %v", err)
	}

	return codeDeveloped, nil
}
