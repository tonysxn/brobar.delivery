package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tonysanin/brobar/order-service/internal/clients"
)

var (
	ErrTimeNotAvailable = errors.New("обраний час недоступний")
	ErrPriceMismatch    = errors.New("ціни змінились, будь ласка, перевірте замовлення")
	ErrProductNotFound  = errors.New("товар не знайдено")
)

type ValidationService struct {
	productClient *clients.ProductClient
	webClient     *clients.WebClient
}

func NewValidationService(productClient *clients.ProductClient, webClient *clients.WebClient) *ValidationService {
	return &ValidationService{
		productClient: productClient,
		webClient:     webClient,
	}
}

var daysMap = map[int]string{
	0: "sunday",
	1: "monday",
	2: "tuesday",
	3: "wednesday",
	4: "thursday",
	5: "friday",
	6: "saturday",
}

// ValidateOrderTime validates that the requested order time is within working hours
func (s *ValidationService) ValidateOrderTime(timeStr string, deliveryType string) error {
	serverTime, err := s.webClient.GetServerTime()
	if err != nil {
		return fmt.Errorf("не вдалося перевірити час: %w", err)
	}

	workingHours, err := s.webClient.GetWorkingHours()
	if err != nil {
		return fmt.Errorf("не вдалося отримати розклад: %w", err)
	}

	// Handle ASAP
	if strings.ToUpper(timeStr) == "ASAP" {
		dayName := daysMap[serverTime.DayNumber]
		currentTime := serverTime.Time[:5] // "HH:MM"

		var schedule clients.DaySchedule
		var ok bool

		if deliveryType == "delivery" {
			schedule, ok = workingHours.Delivery[dayName]
		} else {
			schedule, ok = workingHours.Pickup[dayName]
		}

		if !ok || schedule.Closed {
			return ErrTimeNotAvailable
		}

		if currentTime < schedule.Start || currentTime >= schedule.End {
			return ErrTimeNotAvailable
		}

		return nil
	}

	// Parse scheduled time "2026-01-18 14:30"
	parts := strings.Split(timeStr, " ")
	if len(parts) != 2 {
		return fmt.Errorf("невірний формат часу")
	}

	datePart := parts[0]
	timePart := parts[1]

	orderDate, err := time.Parse("2006-01-02", datePart)
	if err != nil {
		return fmt.Errorf("невірний формат дати")
	}

	serverDate, err := time.Parse("2006-01-02", serverTime.Date)
	if err != nil {
		return fmt.Errorf("помилка серверного часу")
	}

	// Check if order date is in the past
	if orderDate.Before(serverDate) {
		return ErrTimeNotAvailable
	}

	// Get schedule for order day
	dayName := daysMap[int(orderDate.Weekday())]
	var schedule clients.DaySchedule
	var ok bool

	if deliveryType == "delivery" {
		schedule, ok = workingHours.Delivery[dayName]
	} else {
		schedule, ok = workingHours.Pickup[dayName]
	}

	if !ok || schedule.Closed {
		return ErrTimeNotAvailable
	}

	// Check if time is within schedule
	if timePart < schedule.Start || timePart > schedule.End {
		return ErrTimeNotAvailable
	}

	// If order is for today, check that time is not in the past
	if orderDate.Equal(serverDate) {
		currentTime := serverTime.Time[:5]
		if timePart <= currentTime {
			return ErrTimeNotAvailable
		}
	}

	return nil
}

// GetDeliveryCost calculates delivery cost based on coordinates, door delivery, and cart total
// returns: total_cost, door_part_cost, zone, error
func (s *ValidationService) GetDeliveryCost(coords string, deliveryDoor bool, cartTotal float64) (float64, float64, *clients.DeliveryZone, error) {
	var zonePrice float64 = 0
	var doorPartPrice float64 = 0
	var zone *clients.DeliveryZone

	if coords != "" {
		// Parse coordinates
		var lat, lng float64
		_, err := fmt.Sscanf(coords, "%f,%f", &lat, &lng)
		if err != nil {
			return 0, 0, nil, fmt.Errorf("невірний формат координат")
		}

		// Determine zone by coordinates
		zone, err = s.webClient.DetermineZoneByCoords(lat, lng)
		if err != nil {
			return 0, 0, nil, fmt.Errorf("адреса поза зоною доставки")
		}

		zonePrice = zone.Price

		// Free delivery if cart total >= freeOrderPrice
		if zone.FreeOrderPrice > 0 && cartTotal >= zone.FreeOrderPrice {
			zonePrice = 0
		}
	}

	if deliveryDoor {
		doorPrice, err := s.webClient.GetDeliveryDoorPrice()
		if err == nil {
			doorPartPrice = doorPrice
		}
	}

	return zonePrice, doorPartPrice, zone, nil
}

// ValidateAndNormalizePaymentMethod validates and normalizes payment method
// Returns: normalized_method, error
func (s *ValidationService) ValidateAndNormalizePaymentMethod(method string) (string, error) {
	method = strings.ToLower(strings.TrimSpace(method))

	switch method {
	case "cash", "terminal":
		return "cash", nil
	case "bank", "online", "card", "cashless":
		return "online", nil
	default:
		return "", fmt.Errorf("невідомий метод оплати: %s", method)
	}
}
