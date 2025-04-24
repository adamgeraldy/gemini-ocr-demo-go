package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"cloud.google.com/go/vertexai/genai"
	"github.com/joho/godotenv"
)

type GeminiService struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiService() (*GeminiService, error) {
	_ = godotenv.Load()

	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("PROJECT_ID environment variable not set")
	}

	modelRegion := os.Getenv("MODEL_REGION")
	if modelRegion == "" {
		return nil, fmt.Errorf("MODEL_REGION environment variable not set")
	}

	modelName := os.Getenv("MODEL_NAME")
	if modelName == "" {
		return nil, fmt.Errorf("MODEL_NAME environment variable not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, modelRegion)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %v", err)
	}
	model := client.GenerativeModel(modelName)

	return &GeminiService{
		client: client,
		model:  model,
	}, nil
}

func (s *GeminiService) ProcessImage(ctx context.Context, base64Data string, mimeType string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 data: %v", err)
	}

	prompt := []genai.Part{
		genai.Text(`You will analyze the provided document (PDF or PNG format) which contains a document that contains a name, ID number, place and date of birth, along with their address. Please extract and structure the following personal and professional details into a precise JSON format:

{"name": "", "id_number": "", "place_of_birth": "", "date_of_birth": "", "address": ""}

Requirements:
- Extract only the information that matches the exact fields above
- Ensure all dates are in a consistent format
- Maintain the original spelling and capitalization of names and places
- Leave fields empty ("") if the information is not clearly stated in the document

Return the JSON response without any additional formatting, comments, or explanations. If the document is unclear or unreadable, respond with "unclear".`),
		&genai.Blob{
			MIMEType: mimeType,
			Data:     data,
		},
	}

	resp, err := s.model.GenerateContent(ctx, prompt...)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %v", err)
	}

	var result string
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			if text, ok := part.(genai.Text); ok {
				result += string(text)
			}
		}
	}

	result = result[8 : len(result)-4]

	return result, nil
}

func (s *GeminiService) Close() {
	if s.client != nil {
		s.client.Close()
	}
}
