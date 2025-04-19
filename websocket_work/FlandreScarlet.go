package websocket_work

const ai_key = ""

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
	// This function should call the AI model and return the response.
	// For now, we will just return a placeholder response.
	return "This is a placeholder response for query: " + query
}
