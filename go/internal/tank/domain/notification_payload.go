package domain

import "time"

type NotificationPayload struct {
	Type      string     `json:"type"`
	Message   *string    `json:"message"`
	Symbol    *string    `json:"symbol"`
	Price     *float64   `json:"price"`
	Currency  *string    `json:"currency"`
	Timestamp *time.Time `json:"timestamp"`
}
