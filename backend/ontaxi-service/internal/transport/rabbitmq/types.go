package rabbitmq

type TaxiRequest struct {
	ChatID  int64  `json:"chat_id"`
	OrderID string `json:"order_id"`
}

type TaxiEstimateReady struct {
	ChatID    int64   `json:"chat_id"`
	OrderID   string  `json:"order_id"`
	Price     float64 `json:"price"`
	PayloadTo string  `json:"payload_to"`
}

type TaxiConfirm struct {
	ChatID  int64  `json:"chat_id"`
	OrderID string `json:"order_id"`
}

type TaxiOrdered struct {
	ChatID  int64  `json:"chat_id"`
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
