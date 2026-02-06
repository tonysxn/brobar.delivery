package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/tonysanin/brobar/pkg/response"
)

type TimeHandler struct {
	location *time.Location
}

func NewTimeHandler(timezone string) *TimeHandler {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.FixedZone("EET", 2*60*60)
	}
	return &TimeHandler{
		location: loc,
	}
}

func (h *TimeHandler) GetServerTime(c fiber.Ctx) error {
	now := time.Now().In(h.location)
	return response.Success(c, fiber.Map{
		"timestamp":  now.Unix(),
		"datetime":   now.Format(time.RFC3339),
		"date":       now.Format("2006-01-02"),
		"time":       now.Format("15:04:05"),
		"day":        now.Weekday().String(),
		"day_number": int(now.Weekday()),
	})
}
