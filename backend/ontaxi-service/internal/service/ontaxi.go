package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
    "log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tonysanin/brobar/ontaxi-service/internal/config"
)

type OntaxiService struct {
	cfg        *config.Config
	client     *http.Client
	localCoords struct {
		Lat float64
		Lng float64
	}
}

func NewOntaxiService(cfg *config.Config) *OntaxiService {
	coordsParts := strings.Split(cfg.OntaxiLocalCoords, ",")
	lat, _ := strconv.ParseFloat(coordsParts[0], 64)
	lng := 0.0
	if len(coordsParts) > 1 {
		lng, _ = strconv.ParseFloat(coordsParts[1], 64)
	}

	return &OntaxiService{
		cfg:    cfg,
		client: &http.Client{},
		localCoords: struct {
			Lat float64
			Lng float64
		}{
			Lat: lat,
			Lng: lng,
		},
	}
}

func (s *OntaxiService) getHeaders() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + s.cfg.OntaxiToken,
		"Content-Type":  "application/json",
	}
}

func (s *OntaxiService) doRequest(method, url string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	for k, v := range s.getHeaders() {
		req.Header.Set(k, v)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ontaxi api error: %s, body: %s", resp.Status, string(respBody))
	}

	return io.ReadAll(resp.Body)
}

func (s *OntaxiService) GetPayloadByCoords(lat, lng float64) (string, error) {
	targetURL := fmt.Sprintf("%sbusiness/places/reverse/?lang=uk&lat=%f&lon=%f", s.cfg.OntaxiBaseURL, lat, lng)
	
	// Override base url logic for this specific specific call if needed, but looks like full URL in PHP
	
	resp, err := s.doRequest("GET", targetURL, nil)
	if err != nil {
		return "", err
	}

	var result struct {
		Data struct {
			Places []struct {
				Lat     float64 `json:"lat"`
				Lon     float64 `json:"lon"`
				Payload string  `json:"payload"`
			} `json:"places"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return "", err
	}

	var closestPayload string
	minDistance := math.MaxFloat64

	for _, place := range result.Data.Places {
		dist := s.haversineDistance(lat, lng, place.Lat, place.Lon)
		if dist < minDistance {
			minDistance = dist
			closestPayload = place.Payload
		}
	}

	if closestPayload == "" {
		return "", fmt.Errorf("no places found")
	}

	return closestPayload, nil
}

func (s *OntaxiService) GetPlaces(query string) ([]byte, error) {
    // Basic implementation if needed, mostly helper for other calls
    targetURL := fmt.Sprintf("%s/business/places?query=%s&cityId=1&lang=uk", strings.TrimSuffix(s.cfg.OntaxiBaseURL, "/"), query)
    return s.doRequest("GET", targetURL, nil)
}


func (s *OntaxiService) GetDestinationPlacePayload(address string, coords string) (string, float64, float64, error) {
    // If coords are provided, use them to find closest from search results
    // Logic from PHP: search by address, then find closest to coords among results
    
    var filterLat, filterLng float64
    var hasCoords bool

    if coords != "" {
        parts := strings.Split(coords, ";") 
        if len(parts) == 2 {
            if lat, err := strconv.ParseFloat(parts[0], 64); err == nil {
                if lng, err := strconv.ParseFloat(parts[1], 64); err == nil {
                    filterLat = lat
                    filterLng = lng
                    hasCoords = true
                }
            }
        }
    }
    
    targetURL := fmt.Sprintf("%sbusiness/places?query=%s&cityId=1&lang=uk", s.cfg.OntaxiBaseURL, url.QueryEscape(address))
    
    resp, err := s.doRequest("GET", targetURL, nil)
    if err != nil {
        return "", 0, 0, err
    }
    
    var result struct {
        Data struct {
            Items []struct {
                Lat     float64 `json:"lat"`
                Lon     float64 `json:"lon"`
                Payload string  `json:"payload"`
            } `json:"items"`
        } `json:"data"`
    }
    
    if err := json.Unmarshal(resp, &result); err != nil {
        return "", 0, 0, err
    }
    
    if len(result.Data.Items) == 0 {
		return "", 0, 0, fmt.Errorf("no places found")
    }

	var closestPayload string
    var closestLat, closestLng float64
	minDistance := math.MaxFloat64
    
    // If we have coords, try to find closest
    if hasCoords {
        for _, place := range result.Data.Items {
            dist := s.haversineDistance(filterLat, filterLng, place.Lat, place.Lon)
            if dist < minDistance {
                minDistance = dist
                closestPayload = place.Payload
                closestLat = place.Lat
                closestLng = place.Lon
            }
        }
    }
    
	if closestPayload == "" {
        // Fallback: take first item (or if no coords provided)
        item := result.Data.Items[0]
        return item.Payload, item.Lat, item.Lon, nil
	}

	return closestPayload, closestLat, closestLng, nil
}


func (s *OntaxiService) GetEstimate(payloadTo string) (float64, error) {
    // PHP: getEstimate uses LAT/LON, not payload?
    // PHP: getEstimate(array $coordsTo) ... place1Lat=... place2Lat=...
    // PHP: public function getDeliveryEstimate(array $coordsTo) call getEstimate then reduce tariffs to find EXPRESS
    
    // Wait, createOrder uses payloads. getEstimate uses coords.
    // So we need coords for estimate.
    
    return 0, fmt.Errorf("not implemented yet")
}

// Rewriting GetEstimate to match PHP logic which takes coords
func (s *OntaxiService) GetDeliveryEstimate(latTo, lngTo float64) (float64, error) {
    targetURL := fmt.Sprintf("%sbusiness/%s/estimate?lang=uk&place1Lat=%f&place1Lon=%f&place2Lat=%f&place2Lon=%f&clientId=%s",
        s.cfg.OntaxiBaseURL, s.cfg.OntaxiBusinessID, s.localCoords.Lat, s.localCoords.Lng, latTo, lngTo, s.cfg.OntaxiClientID)
        
    resp, err := s.doRequest("GET", targetURL, nil)
    if err != nil {
        return 0, err
    }
    
    // Log raw response for debugging
    log.Printf("Ontaxi estimate response: %s", string(resp))
    
    var result struct {
        Data struct {
            Tariffs []struct {
                ID       string `json:"id"`
                Estimate struct {
                    Cost json.Number `json:"cost"`
                } `json:"estimate"`
            } `json:"tariffs"`
        } `json:"data"`
    }
    
    if err := json.Unmarshal(resp, &result); err != nil {
        return 0, err
    }
    
    for _, tariff := range result.Data.Tariffs {
        if tariff.ID == "EXPRESS" {
            cost, _ := tariff.Estimate.Cost.Float64()
            return cost, nil
        }
    }
    
    return 0, fmt.Errorf("tariff EXPRESS not found")
}

func (s *OntaxiService) CreateOrder(payloadTo string, phone string, name string, comment string, entrance string, doorToDoor bool) (string, error) {
	// Local payload lookup
	localPayload, err := s.getLocalPlacePayload()
	if err != nil {
		return "", fmt.Errorf("failed to get local payload: %w", err)
	}

	var porchValue interface{} = nil
	if entrance != "" {
		if val, err := strconv.Atoi(entrance); err == nil {
			porchValue = val
		}
	}
	
	body := map[string]interface{}{
		"lang":     "uk",
		"clientId": s.cfg.OntaxiClientID,
		"route": map[string]string{
			"place1Payload": localPayload,
			"place2Payload": payloadTo,
		},
		"options": map[string]interface{}{
			"tariff":             "EXPRESS",
			"paymentMethod":      4,
			"paymentMethodId":    s.cfg.OntaxiPaymentMethodID,
			"porch":              porchValue,
			"comment":            comment,
			"tips":               0,
			"bigLuggage":         false,
			"airConditioner":     false,
			"soberDriver":        false,
			"childSeat":          false,
			"doNotCall":          false,
			"notEco":             false,
			"animals":            false,
			"doorToDoor":         doorToDoor,
			"trailer":            false,
			"driverMask":         false,
			"fridge":             false,
			"garbageRemoval":     false,
			"oneLoader":          false,
			"twoLoader":          false,
			"englishSpeak":       false,
			"volunteerTransfer":  false,
			"volunteerDelivery":  false,
			"sosRecharge":        false,
			"driverHelp":         false,
			"externalName":       name,
			"externalPhone":      phone,
		},
	}

	targetURL := fmt.Sprintf("%sbusiness/%s/orders", s.cfg.OntaxiBaseURL, s.cfg.OntaxiBusinessID)
	resp, err := s.doRequest("POST", targetURL, body)
	if err != nil {
		return "", err
	}
    
    // Parse response to get Order ID?
    // PHP just returns json.
    return string(resp), nil
}

func (s *OntaxiService) getLocalPlacePayload() (string, error) {
    // PHP: getPlaces($localPlace)
    resp, err := s.doRequest("GET", fmt.Sprintf("%sbusiness/places?query=%s&cityId=1&lang=uk", s.cfg.OntaxiBaseURL, url.QueryEscape(s.cfg.OntaxiLocalPlace)), nil)
    if err != nil {
        return "", err
    }
    
    var result struct {
        Data struct {
            Items []struct {
                Lat     float64 `json:"lat"`
                Lon     float64 `json:"lon"`
                Payload string  `json:"payload"`
            } `json:"items"`
        } `json:"data"`
    }
    
    if err := json.Unmarshal(resp, &result); err != nil {
        return "", err
    }
    
    // Exact match coords
    for _, place := range result.Data.Items {
        // rough float comparison
        if math.Abs(place.Lat - s.localCoords.Lat) < 0.0001 && math.Abs(place.Lon - s.localCoords.Lng) < 0.0001 {
            return place.Payload, nil
        }
    }
    
    if len(result.Data.Items) > 0 {
        return result.Data.Items[0].Payload, nil
    }
    return "", fmt.Errorf("local place not found")
}

func (s *OntaxiService) haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0

	latFrom := lat1 * math.Pi / 180
	lonFrom := lon1 * math.Pi / 180
	latTo := lat2 * math.Pi / 180
	lonTo := lon2 * math.Pi / 180

	latDelta := latTo - latFrom
	lonDelta := lonTo - lonFrom

	a := math.Sin(latDelta/2)*math.Sin(latDelta/2) +
		math.Cos(latFrom)*math.Cos(latTo)*
			math.Sin(lonDelta/2)*math.Sin(lonDelta/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
