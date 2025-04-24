# Gemini OCR Demo - Go Backend

This project is a Go backend service that uses Google's Gemini model (defaults to Gemini 2.0 Flash) via Vertex AI to perform Optical Character Recognition (OCR) on uploaded images (PNG) or documents (PDF). It extracts specific personal details (name, ID number, place/date of birth, address) and returns them in a structured JSON format.

## Disclaimer

**This is a demonstration project.** It has not been extensively tested and is intended for illustrative purposes only. It should **not** be used in a production environment without further development, testing, and security hardening.

## Features

*   Accepts base64 encoded image (PNG) or PDF data.
*   Uses Gemini 2.0 Flash for OCR and data extraction.
*   Extracts predefined fields: `name`, `id_number`, `place_of_birth`, `date_of_birth`, `address`.
*   Returns extracted data in a clean JSON format.
*   Provides a simple REST API endpoint (`/ocr`).
*   Includes CORS middleware for frontend integration.
*   Configurable via environment variables.
*   Includes a Dockerfile for easy containerization.

## Prerequisites

*   Go (version 1.18 or later recommended)
*   Docker (optional, for containerized deployment)
*   Google Cloud Project with Vertex AI API enabled.
*   Service Account credentials for accessing Vertex AI (or configured Application Default Credentials).

## Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/adamgeraldy/gemini-ocr-demo-go
    cd gemini-ocr-demo-go
    ```

2.  **Set up Environment Variables:**
    Edit the `.env` file in the project root directory and add the following variables:

    ```dotenv
    # .env
    PROJECT_ID=YOUR_PROJECT_ID
    MODEL_REGION=us-central1 # Or your preferred region for the model
    MODEL_NAME=gemini-2.0-flash-001 # Or the specific Gemini model you want to use
    ```

    *   Replace `your-gcp-project-id` with your actual Google Cloud Project ID.
    *   Ensure the specified `MODEL_REGION` supports the `MODEL_NAME`.
    *   This demo utilizes Application Default Credentials (ADC) for authentication.

3.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

## Running the Application

### Locally

```bash
go run main.go
```

The server will start, on port 8080.

### Using Docker

1.  **Build the Docker image:**
    ```bash
    docker build -t gemini-ocr-demo-go .
    ```

2.  **Run the Docker container:**
    Make sure to pass the environment variables to the container. You can do this using the `--env-file` flag.
    ```bash
    docker run -p 8080:8080 --env-file .env gemini-ocr-backend
    ```
    *(Adjust the port mapping `-p 8080:8080` if your `PORT` environment variable is different)*

## Usage

Send a POST request to the `/ocr` endpoint with a JSON payload containing the base64 encoded file data and its MIME type.

**Example using `curl`:**

```bash
curl -X POST http://localhost:8080/ocr \
-H "Content-Type: application/json" \
-d '{
  "file": "data:image/png;base64,YOURIMAGEBASE64", # Replace with your actual base64 string
  "mimeType": "image/png" # Or "application/pdf"
}'
```

**Successful Response (200 OK):**

```json
{
  "result": "{\"name\": \"Jane Doe\", \"id_number\": \"123456789\", \"place_of_birth\": \"Anytown\", \"date_of_birth\": \"1990-01-15\", \"address\": \"123 Main St, Anytown\"}",
  "error": ""
}
```

**Error Response (e.g., 400 Bad Request or 500 Internal Server Error):**

```json
{
  "result": "",
  "error": "Invalid request format: <error details>"
}
```