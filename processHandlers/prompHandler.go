package processhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	eventhandlers "geminiProject/eventHandlers"
	"geminiProject/gemini"
	"geminiProject/utils"
	"log"
	"time"

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
	userInput := "terminal base game"
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
	time.Sleep(5 * time.Second)
	//database to use
	databaseToUse, errorDatabase := fetchDatabase(userInput, stackToUse, approachToUse)
	fmt.Println("Database to use:", databaseToUse)
	if errorDatabase != nil {
		fmt.Println("Error fetching Database:", errorDatabase)
		return ""
	}
	time.Sleep(2 * time.Second)
	projectStructure, errorStructure := eventhandlers.CodeDesigner(userInput, stackToUse, approachToUse, databaseToUse)
	fmt.Println("Project Structure:", projectStructure)
	if errorStructure != nil {
		fmt.Println("Error fetching Structure:", errorStructure)
		return ""
	}
	time.Sleep(1 * time.Second)
	baseLogic, errBaseLogic := BaseLogic(userInput, stackToUse, approachToUse, databaseToUse, projectStructure)
	if errBaseLogic != nil {
		fmt.Println("Error fetching Base Logic:", errBaseLogic)
		return ""
	}
	fmt.Println("baselogic:", baseLogic)

	//code, errorCode := codeDeveloper(userInput, stackToUse, approachToUse, databaseToUse, projectStructure, baseLogic)
	code, errorCode := codeDeveloper(userInput, stackToUse, approachToUse, databaseToUse, projectStructure, baseLogic)
	fmt.Println("Code:", code)
	if errorCode != nil {
		fmt.Println("Error fetching code:", errorCode)
		return ""
	}
	//Language to use
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
	a := string(marshalResponse)
	fmt.Println("RANDOM A:", a)
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
	prompt += fmt.Sprintf("\n tell me just the name of **ideal** languaguage or framework stack which we should use for %s?\n, return clean string", approachToUse)

	ctx := context.Background()

	// Generate content using the Gemini client
	resp, err := clientGemini.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("error generating content: %v", err)
	}
	marshalResponse, _ := json.MarshalIndent(resp, "", "  ")
	b := string(marshalResponse)
	fmt.Println("RANDON B:", b)
	var generateResponse ContentResponse
	if err := json.Unmarshal(marshalResponse, &generateResponse); err != nil {
		log.Fatal(err)
	}
	var bestStack string
	for _, cad := range *generateResponse.Candidates {
		if cad.Content != nil && len(cad.Content.Parts) > 0 {
			bestStack = cad.Content.Parts[0]
		}
	}

	return bestStack, nil
}

func fetchDatabase(input string, stackToUse string, approachToUse string) (string, error) {
	clientGemini := gemini.GetGeminiCLient()

	// Construct the prompt
	// Construct the prompt
	prompt := fmt.Sprintf("Given the requirement:\n%s\n\nConsidering the %s stack, for %s, what's the ideal database for this use case?\nIf no database is needed, please reply with '**NoNeed**'. Reply in 1 word\n", input, stackToUse, approachToUse)

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
	Database := ""
	for _, cad := range *generateResponse.Candidates {
		if cad.Content != nil {
			for _, part := range cad.Content.Parts {
				Database = part
			}
		}
	}

	return Database, nil
}
