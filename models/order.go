package models

import "time"

// OrderStatus represents the status of an order
type OrderStatus int

const (
	OrderStatusCreated   OrderStatus = 0 // created
	OrderStatusPending   OrderStatus = 1 // pending
	OrderStatusPaid      OrderStatus = 2 // paid
	OrderStatusConfirmed OrderStatus = 3 // confirmed
)

// String returns the string representation of OrderStatus
func (s OrderStatus) String() string {
	switch s {
	case OrderStatusCreated:
		return "created"
	case OrderStatusPending:
		return "pending"
	case OrderStatusPaid:
		return "paid"
	case OrderStatusConfirmed:
		return "confirmed"
	default:
		return "unknown"
	}
}

// MarshalJSON implements json.Marshaler interface
func (s OrderStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

// Order represents an order entity
type Order struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	UserID       uint        `gorm:"not null" json:"user_id"`
	ExperienceID uint        `gorm:"not null" json:"experience_id"`
	Price        int         `json:"price"`
	Status       OrderStatus `gorm:"type:int;default:0" json:"status"`
	OrderNo      string      `gorm:"type:varchar(100);uniqueIndex" json:"order_no"`
	OutOrderNo   string      `gorm:"type:varchar(100)" json:"out_order_no"`
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	User         User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Experience   Experience  `gorm:"foreignKey:ExperienceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"experience"`
}

// TableName specifies the table name for the Order model
func (Order) TableName() string {
	return "orders"
}

// IsCreated checks if the order status is created
func (o *Order) IsCreated() bool {
	return o.Status == OrderStatusCreated
}

// IsPending checks if the order status is pending
func (o *Order) IsPending() bool {
	return o.Status == OrderStatusPending
}

// IsPaid checks if the order status is paid
func (o *Order) IsPaid() bool {
	return o.Status == OrderStatusPaid
}

// IsConfirmed checks if the order status is confirmed
func (o *Order) IsConfirmed() bool {
	return o.Status == OrderStatusConfirmed
}

// SetStatus updates the order status
func (o *Order) SetStatus(status OrderStatus) {
	o.Status = status
}
