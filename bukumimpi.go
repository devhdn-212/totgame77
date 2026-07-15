package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

const bukuMimpiUpstreamURL = "https://api.lobaapps.com/api/bukumimpi"

type bukuMimpiRequest struct {
	Tipe string `json:"tipe"`
	Nama string `json:"nama"`
}

type bukuMimpiHandler struct {
	httpClient *http.Client
}

func newBukuMimpiHandler() *bukuMimpiHandler {
	return &bukuMimpiHandler{httpClient: &http.Client{Timeout: 15 * time.Second}}
}

// Search proxies to api.lobaapps.com server-side, since that upstream API
// doesn't send CORS headers and browsers refuse to call it directly.
func (h *bukuMimpiHandler) Search(c *fiber.Ctx) error {
	var req bukuMimpiRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(apiResponse{
			Status: http.StatusBadRequest, Message: "invalid request body",
		})
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(apiResponse{
			Status: http.StatusInternalServerError, Message: "internal server error",
		})
	}

	upstreamReq, err := http.NewRequestWithContext(c.Context(), http.MethodPost, bukuMimpiUpstreamURL, bytes.NewReader(payload))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(apiResponse{
			Status: http.StatusInternalServerError, Message: "internal server error",
		})
	}
	upstreamReq.Header.Set("content-type", "application/json")

	res, err := h.httpClient.Do(upstreamReq)
	if err != nil {
		log.Println("bukumimpi proxy failed:", err)
		return c.Status(http.StatusBadGateway).JSON(apiResponse{
			Status: http.StatusBadGateway, Message: "failed to reach bukumimpi service",
		})
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(apiResponse{
			Status: http.StatusInternalServerError, Message: "internal server error",
		})
	}

	c.Status(res.StatusCode)
	c.Set("content-type", "application/json")
	return c.Send(body)
}
