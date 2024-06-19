package processhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"geminiProject/gemini"
	"log"

	"github.com/google/generative-ai-go/genai"
)

func Rester(projectStructure string, baseLogic string) (string, error) {
	baseLogicv1 := checkConnectivityInternally(projectStructure, baseLogic)
	retestedCode, errorRetestedCode := testCode(projectStructure, baseLogicv1)
	if errorRetestedCode != nil {
		return "", errorRetestedCode
	}
	return retestedCode, nil
}

func checkConnectivityInternally(projectStructure string, baseLogic string) string {
	clientGemini := gemini.GetGeminiCLient()
	ctx := context.Background()
	//Prompt to check if base logic is accurate and consistent for the project structure
	prompt := fmt.Sprintf(`
		Given the following project structure and base logic:

		- Project Structure:
		%s

		- Base Logic:
		%s

		Please verify if all files are correctly connected, routed, and named according to the project requirements. Ensure all dependencies and file relationships are correctly specified.
		`, projectStructure, baseLogic)

	resp, err := clientGemini.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return ""
	}

	if len(resp.Candidates) == 0 {
		return ""
	}

	marshalResponse, _ := json.MarshalIndent(resp, "", "  ")

	var generateResponse ContentResponse
	if err := json.Unmarshal(marshalResponse, &generateResponse); err != nil {
		log.Fatal(err)
	}
	for _, cad := range *generateResponse.Candidates {
		if cad.Content != nil {
			for _, part := range cad.Content.Parts {
				return part
			}
		}
	}
	return projectStructure
}

func testCode(projectStructure string, baseLogic string) (string, error) {
	return projectStructure, nil
}
