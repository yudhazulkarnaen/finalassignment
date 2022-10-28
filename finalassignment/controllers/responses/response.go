package responses

type Message struct {
	Message string `json:"message"`
}

type ErrorMessage struct {
	ErrorMessage string `json:"error_message"`
}
