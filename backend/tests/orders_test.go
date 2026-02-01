package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonysanin/brobar/pkg/helpers"
)

func init() {
	// Try to load .env from project root
	_ = godotenv.Load("../../.env")
}

func getAPIURL() string {
	port := helpers.GetEnv("GATEWAY_PORT", "8000")
	return fmt.Sprintf("http://localhost:%s", port)
}

func getDBDSN() string {
	host := helpers.GetEnv("ORDER_DB_HOST", "localhost")
	if host == "order_db" {
		host = "localhost" // Support running from host
	}
	port := helpers.GetEnv("ORDER_DB_PORT", "5435")
	if port == "5432" {
		port = "5435" // Map internal port to host port for tests
	}
	user := helpers.GetEnv("ORDER_DB_USER", "sanin")
	password := helpers.GetEnv("ORDER_DB_PASSWORD", "saninsanin111444")
	dbname := helpers.GetEnv("ORDER_DB_NAME", "sanin")
	sslmode := helpers.GetEnv("DB_SSLMODE", "disable")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
}

type OrderItemInput struct {
	ProductID          string `json:"product_id"`
	ProductVariationID string `json:"product_variation_id,omitempty"`
	Quantity           int    `json:"quantity"`
}

type CreateOrderInput struct {
	Name           string           `json:"name"`
	Phone          string           `json:"phone"`
	Email          string           `json:"email"`
	DeliveryTypeID string           `json:"delivery_type_id"`
	Address        string           `json:"address"`
	Coords         string           `json:"coords,omitempty"`
	Entrance       string           `json:"entrance,omitempty"`
	DeliveryDoor   bool             `json:"delivery_door,omitempty"`
	Time           string           `json:"time"`
	PaymentMethod  string           `json:"payment_method"`
	Items          []OrderItemInput `json:"items"`
	ClientTotal    float64          `json:"client_total"`
	DeliveryCost   float64          `json:"delivery_cost"`
}

func sendOrder(t *testing.T, input CreateOrderInput) map[string]interface{} {
	body, err := json.Marshal(input)
	require.NoError(t, err)

	resp, err := http.Post(getAPIURL()+"/orders", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	return result
}

func TestOrders(t *testing.T) {
	today := time.Now().Format("2006-01-02")
	testTime := today + " 18:00"

	t.Run("PickupOrder", func(t *testing.T) {
		input := CreateOrderInput{
			Name:           "Go Pickup User",
			Phone:          "+380991111111",
			DeliveryTypeID: "pickup",
			Address:        "Pickup Store",
			Time:           testTime,
			PaymentMethod:  "cash",
			Items: []OrderItemInput{
				{ProductID: "9856c904-db64-48ed-866d-92c99b26bb6b", Quantity: 1},
			},
			ClientTotal: 1,
		}

		result := sendOrder(t, input)
		assert.True(t, result["success"].(bool), fmt.Sprintf("Order failed: %v", result["error"]))
	})

	t.Run("PaidDeliveryWithVariation", func(t *testing.T) {
		input := CreateOrderInput{
			Name:           "Go Paid User",
			Phone:          "+380992222222",
			DeliveryTypeID: "delivery",
			Address:        "Paid Addr",
			Coords:         "50.0101,36.2501",
			Time:           testTime,
			PaymentMethod:  "cash",
			Items: []OrderItemInput{
				{
					ProductID:          "8088872b-baa4-4033-a6db-4ee501876363",
					ProductVariationID: "ece126c3-4dac-455c-9ce5-34c4bea8caa4",
					Quantity:           1,
				},
			},
			ClientTotal:  200,
			DeliveryCost: 150,
		}

		result := sendOrder(t, input)
		assert.True(t, result["success"].(bool), fmt.Sprintf("Order failed: %v", result["error"]))
	})

	t.Run("FreeDelivery", func(t *testing.T) {
		input := CreateOrderInput{
			Name:           "Go Free User",
			Phone:          "+380993333333",
			DeliveryTypeID: "delivery",
			Address:        "Free Addr",
			Coords:         "50.0101,36.2501",
			Time:           testTime,
			PaymentMethod:  "cash",
			Items: []OrderItemInput{
				{ProductID: "df26d5b7-4766-4b71-9580-87d7fd4d1ab2", Quantity: 1},
			},
			ClientTotal:  840,
			DeliveryCost: 0,
		}

		result := sendOrder(t, input)
		assert.True(t, result["success"].(bool), fmt.Sprintf("Order failed: %v", result["error"]))
	})

	t.Run("DoorDelivery", func(t *testing.T) {
		input := CreateOrderInput{
			Name:           "Go Door User",
			Phone:          "+380994444444",
			DeliveryTypeID: "delivery",
			Address:        "Door Addr 123",
			Entrance:       "3",
			Coords:         "50.0101,36.2501",
			DeliveryDoor:   true,
			Time:           testTime,
			PaymentMethod:  "cash",
			Items: []OrderItemInput{
				{ProductID: "df26d5b7-4766-4b71-9580-87d7fd4d1ab2", Quantity: 1},
			},
			ClientTotal:  890,
			DeliveryCost: 50,
		}

		result := sendOrder(t, input)
		assert.True(t, result["success"].(bool), fmt.Sprintf("Order failed: %v", result["error"]))
	})
}
