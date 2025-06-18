package models

import "time"

// Invoice describes the JSON payload sent to the microservice
type Invoice struct {
	ID        int       `json:"id"`
	WidgetID  int       `json:"widget_id"`
	Amount    int       `json:"amount"`
	Product   string    `json:"product"`
	Quantity  int       `json:"quantity"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
