package models

import (
	"context"
	"time"
)

// Yoyo is the type for all Yoyo
type Item struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	InventoryLevel int       `json:"inventory_level"`
	Description    string    `json:"description"`
	Price          int       `json:"price"`
	Image          string    `json:"image"`
	IsRecurring    bool      `json:"is_recurring"`
	PlanID         string    `json:"plan_id"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

// GetYoyo gets one yoyo by id
func (m *DBModel) GetItem(id int) (Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var item Item

	row := m.DB.QueryRowContext(ctx, `
		select 
			id, name, description, inventory_level, price, coalesce(image, ''),
			is_recurring, plan_id,
			created_at, updated_at
		from
			items
		where id = ?`, id)

	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.InventoryLevel,
		&item.Description,
		&item.Price,
		&item.Image,
		&item.IsRecurring,
		&item.PlanID,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return item, err
	}

	return item, nil
}
