package models

import "time"

type Review struct {
	ID            string    `db:"id" json:"id"`
	FoodRating    int       `db:"food_rating" json:"food_rating"`
	ServiceRating int       `db:"service_rating" json:"service_rating"`
	Comment       string    `db:"comment" json:"comment"`
	Phone         *string   `db:"phone" json:"phone,omitempty"`
	Email         *string   `db:"email" json:"email,omitempty"`
	Name          *string   `db:"name" json:"name,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
