package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const anthropicMessagesURL = "https://api.anthropic.com/v1/messages"

const ocrPrompt = `Kamu membaca foto kertas taruhan togel/lotere tulisan tangan. Tiap baris berisi label (4D/3D/2D), nomor, dan nilai taruhan (bet), dalam urutan dan pemisah apapun (titik dua, spasi, dsb).

Baca setiap baris dan keluarkan HANYA JSON array (tanpa markdown, tanpa penjelasan), setiap item berbentuk:
{"number": "<nomor sebagai digit saja>", "bet": "<nilai bet sebagai digit saja>"}

Abaikan label (4D/3D/2D) itu sendiri. Kalau tidak ada baris yang terbaca, keluarkan [].`

type ocrRow struct {
	Number string `json:"number"`
	Bet    string `json:"bet"`
}

type apiResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Record  any    `json:"record"`
}

type claudeContentBlock struct {
	Type   string `json:"type"`
	Text   string `json:"text,omitempty"`
	Source *struct {
		Type      string `json:"type"`
		MediaType string `json:"media_type"`
		Data      string `json:"data"`
	} `json:"source,omitempty"`
}

type claudeMessage struct {
	Role    string                `json:"role"`
	Content []claudeContentBlock `json:"content"`
}

type claudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []claudeMessage `json:"messages"`
}

type claudeResponse struct {
	Content []claudeContentBlock `json:"content"`
	Error   *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type ocrHandler struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

func newOcrHandler(apiKey, model string) *ocrHandler {
	return &ocrHandler{
		apiKey:     apiKey,
		model:      model,
		httpClient: &http.Client{Timeout: 25 * time.Second},
	}
}

func (h *ocrHandler) Scan(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(apiResponse{
			Status: http.StatusBadRequest, Message: "image file is required",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(apiResponse{
			Status: http.StatusInternalServerError, Message: "internal server error",
		})
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(apiResponse{
			Status: http.StatusInternalServerError, Message: "internal server error",
		})
	}

	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	rows, err := h.scanBetSlip(ctx, imageBytes, mimeType)
	if err != nil {
		log.Println("ocr scan failed:", err)
		return c.Status(http.StatusInternalServerError).JSON(apiResponse{
			Status: http.StatusInternalServerError, Message: "failed to read image",
		})
	}

	return c.JSON(apiResponse{Status: http.StatusOK, Message: "success", Record: rows})
}

func (h *ocrHandler) scanBetSlip(ctx context.Context, imageBytes []byte, mimeType string) ([]ocrRow, error) {
	if h.apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is not configured")
	}

	reqBody := claudeRequest{
		Model:     h.model,
		MaxTokens: 1024,
		Messages: []claudeMessage{
			{
				Role: "user",
				Content: []claudeContentBlock{
					{
						Type: "image",
						Source: &struct {
							Type      string `json:"type"`
							MediaType string `json:"media_type"`
							Data      string `json:"data"`
						}{
							Type:      "base64",
							MediaType: mimeType,
							Data:      base64.StdEncoding.EncodeToString(imageBytes),
						},
					},
					{Type: "text", Text: ocrPrompt},
				},
			},
		},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, anthropicMessagesURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", h.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	res, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call anthropic api: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var parsed claudeResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		if parsed.Error != nil {
			return nil, fmt.Errorf("anthropic api error: %s", parsed.Error.Message)
		}
		return nil, fmt.Errorf("anthropic api returned status %d", res.StatusCode)
	}

	if len(parsed.Content) == 0 {
		return []ocrRow{}, nil
	}

	return extractRows(parsed.Content[0].Text)
}

// Claude is asked to return raw JSON but may still wrap it in prose or a
// markdown fence, so pull out the first top-level array before decoding.
func extractRows(text string) ([]ocrRow, error) {
	start := strings.IndexByte(text, '[')
	end := strings.LastIndexByte(text, ']')
	if start == -1 || end == -1 || end < start {
		return []ocrRow{}, nil
	}

	var rows []ocrRow
	if err := json.Unmarshal([]byte(text[start:end+1]), &rows); err != nil {
		return nil, fmt.Errorf("unmarshal ocr rows: %w", err)
	}
	return rows, nil
}
