package handlers

import (
	"net/http"
	"strings"

	"gemini-ocr-demo-go/models"
	"gemini-ocr-demo-go/services"

	"github.com/gin-gonic/gin"
)

type OCRHandler struct {
	geminiService *services.GeminiService
}

func NewOCRHandler(geminiService *services.GeminiService) *OCRHandler {
	return &OCRHandler{
		geminiService: geminiService,
	}
}

func (h *OCRHandler) ProcessOCR(c *gin.Context) {
	var request models.OCRRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.OCRResponse{
			Error: "Invalid request format: " + err.Error(),
		})
		return
	}

	base64Data := request.File
	if idx := strings.Index(base64Data, ";base64,"); idx > 0 {
		base64Data = base64Data[idx+8:]
	}

	mimeType := request.MimeType

	result, err := h.geminiService.ProcessImage(c.Request.Context(), base64Data, mimeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.OCRResponse{
			Error: "Failed to process image: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.OCRResponse{
		Result: result,
	})
}
