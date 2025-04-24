package models

type OCRResponse struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}
