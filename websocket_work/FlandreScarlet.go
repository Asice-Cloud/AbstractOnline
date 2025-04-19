package websocket_work

import (
	"bytes"
	"fmt"
	"github.com/bytedance/sonic"
	"io"
	"net/http"
)

// this is a fake api
const api_key = "sk-1234abcd1234abcd1234abcd1234abcd1234abcd"
const openAIAPIURL = "https://api.openai.com/v1/chat/completions"

type ai_request struct {
	Model   string `json:"model"`
	Message []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type ai_reqponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func get_ai_response(query string) string {

	if api_key == "" {
		return "API key is not set"
	}

	request_body := ai_request{
		Model: "gpt-3.5-turbo",
		Message: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: query,
			},
		},
	}

	json_data, err := sonic.Marshal(request_body)
	if err != nil {
		return fmt.Sprintf("Error marshalling request body: %v", err)
	}

	req, err := http.NewRequest("POST", openAIAPIURL, bytes.NewBuffer(json_data))
	if err != nil {
		return fmt.Sprintf("Error: Failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api_key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("Error: Failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error: Failed to read response body: %v", err)
	}
	var aiResponse ai_reqponse
	err = sonic.Unmarshal(body, &aiResponse)
	if err != nil {
		return fmt.Sprintf("Error: Failed to unmarshal response: %v", err)
	}
	// Return the AI's response
	if len(aiResponse.Choices) > 0 {
		return aiResponse.Choices[0].Message.Content
	} else {
		return "This is a placeholder response for query: " + query
	}
}
