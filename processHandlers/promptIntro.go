package processhandlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"geminiProject/gemini"
	"geminiProject/utils"
	"log"
	"os"
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

func ReadRequest() string {
	//Start the promptReader
	// Receive input from the user
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter input: ")
	userInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}
	userInput = strings.TrimSpace(userInput)

	approachToUse, errorApproach := fetchApproach(userInput)
	if errorApproach != nil {
		fmt.Println("Error fetching approach:", errorApproach)
		return ""
	}
	fmt.Print("Approach to use:", approachToUse)

	//Return the user input
	return userInput
}

func fetchApproach(input string) (string, error) {
	clientGemini := gemini.GetGeminiCLient()

	// Construct the prompt
	prompt := "Given the following approaches:\n"
	for key := range utils.PathApproach {
		prompt += "- " + key + "\n"
	}
	prompt += "\nWhat approach should be used for the following requirement:\n" + input + " \n please reply in one word only"

	ctx := context.Background()

	// Generate content using the Gemini client
	resp, err := clientGemini.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("error generating content: %v", err)
	}
	marshalResponse, _ := json.MarshalIndent(resp, "", "  ")

	var generateResponse ContentResponse
	if err := json.Unmarshal(marshalResponse, &generateResponse); err != nil {
		log.Fatal(err)
	}
	for _, cad := range *generateResponse.Candidates {
		if cad.Content != nil {
			for _, part := range cad.Content.Parts {
				return part, nil
			}
		}
	}

	return "No suitable approach found", nil
}
