package clients

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/tonysanin/brobar/pkg/helpers"
)

type WebClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewWebClient() *WebClient {
	return &WebClient{
		baseURL: helpers.GetEnv("WEB_SERVICE_URL", "http://web-service-dev:3006"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Settings response
type SettingsResponse struct {
	Success bool      `json:"success"`
	Data    []Setting `json:"data"`
}

type Setting struct {
	Key   string `json:"key"`
	Type  string `json:"setting_type"`
	Value string `json:"value"`
}

// Time response
type TimeResponse struct {
	Success bool     `json:"success"`
	Data    TimeData `json:"data"`
}

type TimeData struct {
	Timestamp int64  `json:"timestamp"`
	Datetime  string `json:"datetime"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	Day       string `json:"day"`
	DayNumber int    `json:"day_number"`
}

// Working hours structures
type DaySchedule struct {
	Start  string `json:"start"`
	End    string `json:"end"`
	Closed bool   `json:"closed"`
}

type WorkingHours struct {
	Delivery map[string]DaySchedule `json:"delivery"`
	Pickup   map[string]DaySchedule `json:"pickup"`
}

func (c *WebClient) GetSettings() ([]Setting, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/settings", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch settings: %w", err)
	}
	defer resp.Body.Close()

	var settingsResp SettingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&settingsResp); err != nil {
		return nil, fmt.Errorf("failed to decode settings response: %w", err)
	}

	return settingsResp.Data, nil
}

func (c *WebClient) GetServerTime() (*TimeData, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/time", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch server time: %w", err)
	}
	defer resp.Body.Close()

	var timeResp TimeResponse
	if err := json.NewDecoder(resp.Body).Decode(&timeResp); err != nil {
		return nil, fmt.Errorf("failed to decode time response: %w", err)
	}

	return &timeResp.Data, nil
}

func (c *WebClient) GetWorkingHours() (*WorkingHours, error) {
	settings, err := c.GetSettings()
	if err != nil {
		return nil, err
	}

	for _, s := range settings {
		if s.Key == "working_hours" {
			var wh WorkingHours
			if err := json.Unmarshal([]byte(s.Value), &wh); err != nil {
				return nil, fmt.Errorf("failed to parse working hours: %w", err)
			}
			return &wh, nil
		}
	}

	return nil, fmt.Errorf("working_hours setting not found")
}

func (c *WebClient) GetDeliveryDoorPrice() (float64, error) {
	settings, err := c.GetSettings()
	if err != nil {
		return 0, err
	}

	for _, s := range settings {
		if s.Key == "delivery_door_price" {
			var price float64
			if _, err := fmt.Sscanf(s.Value, "%f", &price); err != nil {
				return 0, fmt.Errorf("failed to parse delivery door price: %w", err)
			}
			return price, nil
		}
	}

	return 50.0, nil // Default
}

func (c *WebClient) GetDeliveryZones() ([]DeliveryZone, error) {
	settings, err := c.GetSettings()
	if err != nil {
		return nil, err
	}

	for _, s := range settings {
		if s.Key == "delivery_zones" {
			var zones []DeliveryZone
			if err := json.Unmarshal([]byte(s.Value), &zones); err != nil {
				return nil, fmt.Errorf("failed to parse delivery zones: %w", err)
			}
			return zones, nil
		}
	}

	return nil, fmt.Errorf("delivery_zones setting not found")
}

type DeliveryZone struct {
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	FreeOrderPrice float64 `json:"freeOrderPrice"`
	Radius         float64 `json:"radius"`
	InnerRadius    float64 `json:"innerRadius"`
}

// ZoneCenter represents the center point of delivery zones
type ZoneCenter struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Default center (fallback if setting not found)
var defaultCenter = ZoneCenter{Lat: 50.0014656, Lng: 36.245192}

// GetZoneCenter fetches the zone center from settings
func (c *WebClient) GetZoneCenter() (*ZoneCenter, error) {
	settings, err := c.GetSettings()
	if err != nil {
		return &defaultCenter, nil
	}

	for _, s := range settings {
		if s.Key == "zone_center" {
			var center ZoneCenter
			if err := json.Unmarshal([]byte(s.Value), &center); err != nil {
				return &defaultCenter, nil
			}
			return &center, nil
		}
	}

	return &defaultCenter, nil
}

// DetermineZoneByCoords finds the delivery zone based on coordinates
func (c *WebClient) DetermineZoneByCoords(lat, lng float64) (*DeliveryZone, error) {
	zones, err := c.GetDeliveryZones()
	if err != nil {
		return nil, err
	}

	center, _ := c.GetZoneCenter()

	// Calculate distance from center using Haversine formula
	distance := haversineDistance(center.Lat, center.Lng, lat, lng)

	// Find matching zone (zone with smallest radius that contains the point)
	for _, zone := range zones {
		if distance >= zone.InnerRadius && distance < zone.Radius {
			return &zone, nil
		}
	}

	return nil, fmt.Errorf("coordinates are outside delivery zones")
}

// haversineDistance calculates distance between two points in km
func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371 // Earth radius in km

	dLat := toRadians(lat2 - lat1)
	dLng := toRadians(lng2 - lng1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRadians(lat1))*math.Cos(toRadians(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func toRadians(deg float64) float64 {
	return deg * math.Pi / 180
}
