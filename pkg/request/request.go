package request

import (
	// for marshalling
	"bytes"
	"log"

	// for json
	"encoding/json"
	"fmt"

	// for getting the api key
	"github.com/vibovenkat123/review-gpt/pkg/globals"
	// for reading the response
	"io"
	// for errors
	"errors"
	// for the requests
	"net/http"

	"github.com/charmbracelet/glamour"
)

// the struct to use for the body of the request
type Body struct {
	Model         string  `json:"model"`
	Prompt        string  `json:"prompt"`
	Temperature   float64 `json:"temperature"`
	Max_Tokens    int     `json:"max_tokens"`
	Top_P         float64 `json:"top_p"`
	Frequence_Pen float64 `json:"frequency_penalty"`
	Presence_Pen  float64 `json:"presence_penalty"`
	Best_Of       int     `json:"best_of"`
}

// the struct to use for the chat models
type ChatBody struct {
	Model         string    `json:"model"`
	Messages      []Message `json:"messages"`
	Temperature   float64   `json:"temperature"`
	Max_Tokens    int       `json:"max_tokens"`
	Top_P         float64   `json:"top_p"`
	Frequence_Pen float64   `json:"frequency_penalty"`
	Presence_Pen  float64   `json:"presence_penalty"`
}

// the text in the choices the response gives
type APIText struct {
	Text    string  `json:"text"`
	Message Message `json:"message"`
	Index   int     `json:"index"`
}

// the usage the response gives
type APIUsage struct {
	Prompt_Tokens     int `json:"prompt_tokens"`
	Completion_Tokens int `json:"completion_tokens"`
	Total_Tokens      int `json:"total_tokens"`
}
type ApiErr struct {
	Message string  `json:"message"`
	Type    string  `json:"type"`
	Param   *string `json:"param"`
	Code    string  `json:"code"`
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// the response the api gives
type APIResponse struct {
	Err     *ApiErr   `json:"error"`
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Choices []APIText `json:"choices"`
	Usage   APIUsage  `json:"usage"`
}

func LogVerbose(msg string) {
	if globals.Verbose {
		globals.Log.Info().
			Msg(msg)
	}
}

// request the api
func RequestApi(gitDiff string, model string, maxtokens int, temperature float64, top_p float64, frequence float64, presence float64, bestof int) {
	LogVerbose("Requesting for improvements")
	// get all the improvements
	improvements, err := RequestImprovements(globals.OpenaiKey, gitDiff, model, maxtokens, temperature, top_p, frequence, presence, bestof)
	LogVerbose("Got improvements")
	if err != nil {
		globals.Log.Error().
			Err(err).
			Msg("Error while getting improvements")
	}
	// print each improvement
	for _, improvement := range improvements {
		renderer, err := glamour.NewTermRenderer(
			glamour.WithStyles(glamour.DraculaStyleConfig), // ASCIIStyle is one of many available styles
		)
		if err != nil {
			log.Fatal(err)
		}

		out, err := renderer.Render(improvement)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(out)
	}
}

// checking the format
func CheckFormat(body Body, model bool) error {
	// model
	if !model {
		return globals.ErrWrongModel
	}
	// temperature
	if body.Temperature < globals.TempRangeMin || body.Temperature > globals.TempRangeMax {
		return globals.ErrWrongTempRange
	}
	// top_p
	if body.Top_P < globals.TopPMin || body.Top_P > globals.TopPMax {
		return globals.ErrWrongToppRange
	}
	// presense penalty
	if body.Presence_Pen < globals.PresenceMin || body.Presence_Pen > globals.PresenceMax {
		return globals.ErrWrongPresenceRange
	}
	// frequence penalty
	if body.Frequence_Pen < globals.FrequenceMin || body.Frequence_Pen > globals.FrequenceMax {
		return globals.ErrWrongFrequenceRange
	}
	// the best_of
	if body.Best_Of < globals.BestOfMin || body.Best_Of > globals.BestOfMax {
		return globals.ErrWrongBestOfRange
	}
	// if its all good
	return nil
}

const (
	API_BASE_URL = "https://api.openai.com/v1/"

	PROMPT_PREFIX = `From a code reviewer's perspective, Review the the git diff below and tell me what I can improve on in the code (the '+' in the git diff is an added line, the '-' is a removed line). Only review the changes that code that has been added i.e. the code denoted by the '+' icon all other codes i.e. codes denoted by '-' and with no indicator, are just for context dont comment on them. Do not suggest changes already made in the git diff. Do not explain the git diff. Only say what could be improved. Focus on what needs to be improved rather than what is already properly implemented. Also go into more detail, give me code snippets of how to enhance the code giving me code suggestions too. Give the response in Markdown`

	CHAT_PROMPT_INSTRUCTIONS = `You are a very intelligentX and professional senior engineer with over 10 years of experience. You have a deep understanding of software engineering principles and best practices. You are also proficient in a variety of programming languages and technologies. You are passionate about writing high-quality code and ensuring that our code is well-reviewed. You review only the added changed code in the while code review. You are also committed to continuous learning and improvement. When reviewing code, You  typically look for the following: Correctness: Does the code work as intended? Readability: Is the code easy to read and understand? Maintainability: Is the code easy to maintain and extend? Performance: Is the code efficient and performant? Security: Is the code secure and free from vulnerabilities? You provide code reviewers  with specific feedback and suggestions for improvement the cod You will take in a git diff, and review it for the user. You will provide user with detailed code review feedback, including the following: File name under 'File Name' section, Line number under 'Line Number' section, Comment under 'Comment' section, Sugegested Refactored code snippet for code that needs refactoring under 'Suggested Change' section. Please also try to provide the user with specific suggestions for improvement, such as: How to make the code more readable, How to improve the performance of the code, How to make the code more secure, How to improve the overall design of the code. The user appreciates your feedback and the user will use it to improve their code.`
)

// [Rest of your code as in the previous refactoring]

// Helper function to configure parameters
func configureParams(model globals.Model, maxtokens int, temperature, top_p, frequence, presence float64, bestof int) (Body, ChatBody) {
	// Normal GPT3 body struct
	params := Body{
		Model:         model.Name,
		Temperature:   temperature,
		Max_Tokens:    maxtokens,
		Top_P:         top_p,
		Frequence_Pen: frequence,
		Presence_Pen:  presence,
		Best_Of:       bestof,
	}

	// Chat models struct
	chatParams := ChatBody{
		Model:         model.Name,
		Temperature:   temperature,
		Max_Tokens:    maxtokens,
		Top_P:         top_p,
		Frequence_Pen: frequence,
		Presence_Pen:  presence,
	}

	return params, chatParams
}

// Helper function to create the request body
func createRequestBody(model globals.Model, params Body, chatParams ChatBody, gitDiff string) (*bytes.Buffer, error) {
	prompt := fmt.Sprintf("%s\n%s\n", PROMPT_PREFIX, gitDiff)
	if model.Chat {
		sysMessage := Message{
			Role:    "system",
			Content: CHAT_PROMPT_INSTRUCTIONS,
		}
		usrMessage := Message{
			Role:    "user",
			Content: gitDiff,
		}
		chatParams.Messages = []Message{sysMessage, usrMessage}
	} else {
		params.Prompt = prompt
	}

	var jsonParams []byte
	var err error
	if model.Chat {
		jsonParams, err = json.Marshal(chatParams)
	} else {
		jsonParams, err = json.Marshal(params)
	}
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonParams), nil
}

// Main function with refactored content
func RequestImprovements(key string, gitDiff string, rawModel string, maxtokens int, temperature float64, top_p float64, frequence float64, presence float64, bestof int) ([]string, error) {
	answers := []string{}
	model := globals.Models[rawModel]
	_, hasModel := globals.Models[rawModel]

	params, chatParams := configureParams(model, maxtokens, temperature, top_p, frequence, presence, bestof)

	// Move this check before createRequestBody()
	if err := CheckFormat(params, hasModel); err != nil {
		return answers, err
	}

	reqBody, err := createRequestBody(model, params, chatParams, gitDiff)
	if err != nil {
		globals.Log.Error().
			Msg(fmt.Sprintf("Error creating request body: %s", err))
		return answers, err
	}

	endUrl := "completions"
	if model.Chat {
		endUrl = "chat/completions"
	}
	url := fmt.Sprintf("%s%s", API_BASE_URL, endUrl)

	LogVerbose("Creating new request")
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		globals.Log.Error().
			Msg(fmt.Sprintf("Error sending request: %s", err))
		return answers, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", key))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	LogVerbose("Requesting GPT")
	resp, err := client.Do(req)
	if err != nil {
		return answers, err
	}
	defer resp.Body.Close()

	LogVerbose("Got back the request information")
	body, _ := io.ReadAll(resp.Body)
	apiReq := APIResponse{}
	err = json.Unmarshal([]byte(string(body)), &apiReq)
	if err != nil {
		globals.Log.Panic().
			Msg(err.Error())
	}
	if apiReq.Err != nil {
		err := apiReq.Err
		return answers, errors.New(err.Message)
	}

	choices := apiReq.Choices
	for _, c := range choices {
		if model.Chat {
			answers = append(answers, c.Message.Content)
			continue
		}
		if len(c.Text) != 0 {
			answers = append(answers, c.Text)
		}
	}
	return answers, nil
}
