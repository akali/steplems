package deepinfra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/vincent-petithory/dataurl"
	"mime/multipart"
	"net/http"
)

type ImageGenerationModel string

const (
	FluxSchnell   ImageGenerationModel = "black-forest-labs/FLUX-1-schnell"
	FluxDev       ImageGenerationModel = "black-forest-labs/FLUX-1-dev"
	StabilitySDXL ImageGenerationModel = "stability-ai/sdxl"
)

type config struct {
	prompt string
	model  ImageGenerationModel
}

func (c *config) String() string {
	return fmt.Sprintf("prompt: %q, model: %q", c.prompt, c.model)
}

type ImageGenerationOption func(*config) error

func WithPrompt(prompt string) ImageGenerationOption {
	return func(c *config) error {
		c.prompt = prompt
		return nil
	}
}

func WithModel(model ImageGenerationModel) ImageGenerationOption {
	return func(c *config) error {
		c.model = model
		return nil
	}
}

type Image struct {
	Data []byte
	URL  string
}

type ImageResponse struct {
	Images []*Image
}

func (r ImageResponse) Urls() []string {
	var urls []string
	for _, image := range r.Images {
		urls = append(urls, image.URL)
	}
	return urls
}

func (r ImageResponse) ImageBytes() [][]byte {
	var result [][]byte
	for _, image := range r.Images {
		result = append(result, image.Data)
	}
	return result
}

type ImageGenerator interface {
	GenerateImage(ctx context.Context, c *Client, igc *config) (*ImageResponse, error)
}

func fromFluxResp(response *FluxResponse) *ImageResponse {
	var images []*Image
	for _, str := range response.Images {
		dataURL, _ := dataurl.DecodeString(str)
		images = append(images, &Image{
			Data: dataURL.Data,
		})
	}
	return &ImageResponse{Images: images}
}

func fromApiResponse(response *APIResponse) *ImageResponse {
	var images []*Image
	for _, url := range response.Output {
		images = append(images, &Image{
			URL: url,
		})
	}
	return &ImageResponse{
		Images: images,
	}
}

// GenerateImageFlux sends a request to the DeepInfra API to generate an image with Flux models.
func (c *Client) generateImageFlux(igc *config) (*ImageResponse, error) {
	var url string
	switch igc.model {
	case FluxSchnell:
		url = fmt.Sprintf("%s/inference/black-forest-labs/FLUX-1-schnell", c.BaseURL)
	case FluxDev:
		url = fmt.Sprintf("%s/inference/black-forest-labs/FLUX-1-dev", c.BaseURL)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	fw, err := writer.CreateFormField("prompt")
	if err != nil {
		return nil, err
	}
	_, err = fw.Write([]byte(igc.prompt))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to generate image: status code %d", resp.StatusCode)
	}

	var fluxResp FluxResponse
	if err := json.NewDecoder(resp.Body).Decode(&fluxResp); err != nil {
		return nil, err
	}
	return fromFluxResp(&fluxResp), nil
}

// generateImageStableDiffusion sends a request to the DeepInfra API to generate an image
func (c *Client) generateImageStableDiffusion(igc *config) (*ImageResponse, error) {
	url := fmt.Sprintf("%s/inference/stability-ai/sdxl", c.BaseURL) // Replace with actual endpoint

	body, err := json.Marshal(&Request{
		Input: &Input{
			Prompt: &igc.prompt,
		},
	})
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

	return fromApiResponse(&apiResponse), nil
}

func (c *Client) createImage(ctx context.Context, igc *config) (*ImageResponse, error) {
	switch igc.model {
	case FluxSchnell, FluxDev:
		return c.generateImageFlux(igc)
	default:
		return c.generateImageStableDiffusion(igc)
	}
}

func (c *Client) CreateImage(ctx context.Context, opts ...ImageGenerationOption) (*ImageResponse, error) {
	igc := &config{}
	for _, opt := range opts {
		if err := opt(igc); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	images, err := c.createImage(ctx, igc)

	if len(images.Images) == 0 {
		err = multierror.Append(err, fmt.Errorf("empty response"))
	}

	return images, err
}
