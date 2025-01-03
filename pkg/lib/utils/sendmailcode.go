package utils

type EmailData struct {
	Sender      map[string]string   `json:"sender"`
	To          []map[string]string `json:"to"`
	Subject     string              `json:"subject"`
	HtmlContent string              `json:"htmlContent"`
}

func SendMailCode() {}
