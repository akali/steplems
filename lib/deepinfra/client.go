package deepinfra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type Client struct {
	BaseURL    string
	APIToken   string
	HTTPClient *http.Client
}

// NewClient creates a new DeepInfra API client
func NewClient(apiToken string) *Client {
	return &Client{
		BaseURL:    "https://api.deepinfra.com/v1",
		APIToken:   apiToken,
		HTTPClient: &http.Client{},
	}
}

type Request struct {
	Input *Input `json:"input,omitempty"`
}

// Input represents the input schema for the API
type Input struct {
	Prompt *string `json:"prompt,omitempty"`
}

// VoiceToTextResponse represents the response structure for the voice-to-text API.
type VoiceToTextResponse struct {
	Text            string          `json:"text"`
	Segments        []Segment       `json:"segments"`
	Language        string          `json:"language"`
	InputLengthMs   int             `json:"input_length_ms"`
	RequestID       *string         `json:"request_id,omitempty"`
	InferenceStatus InferenceStatus `json:"inference_status"`
}

// Segment represents each segment in the voice-to-text response.
type Segment struct {
	ID    int     `json:"id"`
	Text  string  `json:"text"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

// InferenceStatus contains details about the inference process
type InferenceStatus struct {
	Status          string  `json:"status"`
	RuntimeMS       int     `json:"runtime_ms"`
	Cost            float64 `json:"cost"`
	TokensGenerated *int    `json:"tokens_generated,omitempty"`
	TokensInput     *int    `json:"tokens_input,omitempty"`
}

// InputDetails contains the details of the input that was used for the inference
type InputDetails struct {
	Prompt            string  `json:"prompt"`
	NegativePrompt    string  `json:"negative_prompt,omitempty"`
	Image             *string `json:"image,omitempty"`
	Mask              *string `json:"mask,omitempty"`
	Width             int     `json:"width"`
	Height            int     `json:"height"`
	NumOutputs        int     `json:"num_outputs"`
	Scheduler         string  `json:"scheduler"`
	NumInferenceSteps int     `json:"num_inference_steps"`
	GuidanceScale     float64 `json:"guidance_scale"`
	PromptStrength    float64 `json:"prompt_strength"`
	Seed              *int    `json:"seed,omitempty"`
	Refine            string  `json:"refine"`
	HighNoiseFrac     float64 `json:"high_noise_frac"`
	RefineSteps       *int    `json:"refine_steps,omitempty"`
	ApplyWatermark    bool    `json:"apply_watermark"`
}

// APIResponse represents the full response from the API
type APIResponse struct {
	RequestID           string                 `json:"request_id"`
	InferenceStatus     InferenceStatus        `json:"inference_status"`
	Input               InputDetails           `json:"input"`
	Output              []string               `json:"output"`
	ID                  string                 `json:"id"`
	Version             *string                `json:"version,omitempty"`
	CreatedAt           *time.Time             `json:"created_at,omitempty"`
	StartedAt           time.Time              `json:"started_at"`
	CompletedAt         time.Time              `json:"completed_at"`
	Logs                string                 `json:"logs,omitempty"`
	Error               *string                `json:"error,omitempty"`
	Status              string                 `json:"status"`
	Metrics             map[string]interface{} `json:"metrics"`
	WebhookEventsFilter []string               `json:"webhook_events_filter,omitempty"`
	Webhook             *string                `json:"webhook,omitempty"`
	OutputFilePrefix    *string                `json:"output_file_prefix,omitempty"`
}

// GenerateImage sends a request to the DeepInfra API to generate an image
func (c *Client) GenerateImage(input *Request) (*APIResponse, error) {
	url := fmt.Sprintf("%s/inference/stability-ai/sdxl", c.BaseURL) // Replace with actual endpoint

	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to generate image: status code %d", resp.StatusCode)
	}

	var apiResponse APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse, nil
}

// VoiceToText sends the audio file to the DeepInfra API and returns the transcription result.
func (c *Client) VoiceToText(audioFilePath string) (*VoiceToTextResponse, error) {
	file, err := os.Open(audioFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("audio", audioFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}
	url := fmt.Sprintf("%s/inference/openai/whisper-large-v3", c.BaseURL)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "bearer "+c.APIToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response VoiceToTextResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
