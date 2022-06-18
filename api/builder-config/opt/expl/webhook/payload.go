package webhook

// https://raw.githubusercontent.com/mattermost/mattermost-server/v7.0.0/model/outgoing_webhook.go

type Request struct {
	Token       string `json:"token"`
	UserName    string `json:"user_name"`
	Text        string `json:"text"`
	TriggerWord string `json:"trigger_word"`
}

type Response struct {
	Text         string `json:"text"`
	ResponseType string `json:"response_type"` // "post" or "comment"
}

func NewResponse(text string) *Response {
	return &Response{
		Text:         text,
		ResponseType: "post",
	}
}
