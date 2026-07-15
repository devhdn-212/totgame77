package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/genai"
)

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

type ocrHandler struct {
	apiKey string
	model  string
}

func newOcrHandler(apiKey, model string) *ocrHandler {
	return &ocrHandler{apiKey: apiKey, model: model}
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
		return nil, fmt.Errorf("GEMINI_API_KEY is not configured")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  h.apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("create genai client: %w", err)
	}

	parts := []*genai.Part{
		{InlineData: &genai.Blob{MIMEType: mimeType, Data: imageBytes}},
		{Text: ocrPrompt},
	}
	contents := []*genai.Content{{Role: "user", Parts: parts}}

	res, err := client.Models.GenerateContent(ctx, h.model, contents, nil)
	if err != nil {
		return nil, fmt.Errorf("call gemini api: %w", err)
	}

	return extractRows(res.Text())
}

// Gemini is asked to return raw JSON but may still wrap it in prose or a
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
