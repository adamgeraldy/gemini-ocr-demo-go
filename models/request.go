package models

type OCRRequest struct {
	File     string `json:"file" binding:"required"`
	MimeType string `json:"mimeType" binding:"required"`
}
