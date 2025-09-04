package dto

type EmailDTO struct {
	To      []string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}
