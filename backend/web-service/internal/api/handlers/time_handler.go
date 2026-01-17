package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/tonysanin/brobar/pkg/response"
)

type TimeHandler struct{}

func NewTimeHandler() *TimeHandler {
	return &TimeHandler{}
}

// Kiev timezone
var kievLocation *time.Location

func init() {
	var err error
	kievLocation, err = time.LoadLocation("Europe/Kiev")
	if err != nil {
		// Fallback to UTC+2 in winter, UTC+3 in summer (approximate)
		kievLocation = time.FixedZone("EET", 2*60*60)
	}
}

func (h *TimeHandler) GetServerTime(c fiber.Ctx) error {
	now := time.Now().In(kievLocation)
	return response.Success(c, fiber.Map{
		"timestamp":  now.Unix(),
		"datetime":   now.Format(time.RFC3339),
		"date":       now.Format("2006-01-02"),
		"time":       now.Format("15:04:05"),
		"day":        now.Weekday().String(),
		"day_number": int(now.Weekday()),
	})
}
