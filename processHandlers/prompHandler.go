package processhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"geminiProject/gemini"
	"geminiProject/utils"
	"log"

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
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Println("Enter input: ")
	// userInput, err := reader.ReadString('\n')
	// if err != nil {
	// 	fmt.Println("Error reading input:", err)
	// 	return ""
	// }
	// userInput = strings.TrimSpace(userInput)
	userInput := "terminal based game"
	//Apporach to use
	approachToUse, errorApproach := fetchApproach(userInput)
	if errorApproach != nil {
		fmt.Println("Error fetching approach:", errorApproach)
		return ""
	}
	fmt.Print("Approach to use:", approachToUse)

	//Stack to use
	stackToUse, errorStack := fetchBestStack(userInput, approachToUse)
	if errorStack != nil {
		fmt.Println("Error fetching stack:", errorApproach)
		return ""
	}
	fmt.Print("Stack to use:", stackToUse)

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

	return "", nil
}

func fetchBestStack(input string, approachToUse string) (string, error) {
	clientGemini := gemini.GetGeminiCLient()

	// Construct the prompt
	// Construct the prompt
	prompt := fmt.Sprintf("Given the requirement:\n%s\n", input)
	prompt += fmt.Sprintf("\n tell me just the name of **ideal** languaguage or framework stack which we should use for %s?\n", approachToUse)

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
	bestStack := ""
	for _, cad := range *generateResponse.Candidates {
		if cad.Content != nil {
			for _, part := range cad.Content.Parts {
				bestStack = part
			}
		}
	}

	return bestStack, nil
}
