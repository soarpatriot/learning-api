package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderModel(t *testing.T) {
	db := InitTestDB()

	// Add Order to test database migration
	db.AutoMigrate(&Order{}, &Experience{}, &Topic{})

	// Create test user
	user := User{
		OpenID: "test_open_id",
		Name:   "Test User",
	}
	db.Create(&user)

	// Create test topic
	topic := Topic{
		Name: "Test Topic",
	}
	db.Create(&topic)

	// Create test experience
	experience := Experience{
		TopicID: topic.ID,
		UserID:  user.ID,
	}
	db.Create(&experience)

	// Create test order
	order := Order{
		UserID:       user.ID,
		ExperienceID: experience.ID,
		Price:        1000,
		Status:       OrderStatusCreated,
		OrderNo:      "ORD20250624001",
		OutOrderNo:   "OUT20250624001",
	}

	// Test creating order
	result := db.Create(&order)
	assert.NoError(t, result.Error)
	assert.NotZero(t, order.ID)
	assert.NotZero(t, order.CreatedAt)
	assert.NotZero(t, order.UpdatedAt)
	assert.Equal(t, OrderStatusCreated, order.Status)

	// Test table name
	assert.Equal(t, "orders", order.TableName())

	// Test retrieving order with associations
	var retrievedOrder Order
	db.Preload("User").Preload("Experience").First(&retrievedOrder, order.ID)

	assert.Equal(t, order.ID, retrievedOrder.ID)
	assert.Equal(t, user.ID, retrievedOrder.UserID)
	assert.Equal(t, experience.ID, retrievedOrder.ExperienceID)
	assert.Equal(t, 1000, retrievedOrder.Price)
	assert.Equal(t, OrderStatusCreated, retrievedOrder.Status)
	assert.Equal(t, "ORD20250624001", retrievedOrder.OrderNo)
	assert.Equal(t, "OUT20250624001", retrievedOrder.OutOrderNo)
	assert.Equal(t, "Test User", retrievedOrder.User.Name)
	assert.Equal(t, user.ID, retrievedOrder.Experience.UserID)

	// Test status updates
	retrievedOrder.Status = OrderStatusPending
	db.Save(&retrievedOrder)

	var updatedOrder Order
	db.First(&updatedOrder, order.ID)
	assert.Equal(t, OrderStatusPending, updatedOrder.Status)

	// Test all status values
	statuses := []OrderStatus{OrderStatusCreated, OrderStatusPending, OrderStatusPaid, OrderStatusConfirmed}
	orderNumbers := []string{"ORD20250624002", "ORD20250624003", "ORD20250624004", "ORD20250624005"}
	outOrderNumbers := []string{"OUT20250624002", "OUT20250624003", "OUT20250624004", "OUT20250624005"}

	for i, status := range statuses {
		testOrder := Order{
			UserID:       user.ID,
			ExperienceID: experience.ID,
			Price:        500,
			Status:       status,
			OrderNo:      orderNumbers[i], // Unique order numbers
			OutOrderNo:   outOrderNumbers[i],
		}
		result := db.Create(&testOrder)
		assert.NoError(t, result.Error)
		assert.Equal(t, status, testOrder.Status)
	}

	// Test helper methods
	createdOrder := Order{Status: OrderStatusCreated}
	assert.True(t, createdOrder.IsCreated())
	assert.False(t, createdOrder.IsPending())
	assert.False(t, createdOrder.IsPaid())
	assert.False(t, createdOrder.IsConfirmed())

	pendingOrder := Order{Status: OrderStatusPending}
	assert.False(t, pendingOrder.IsCreated())
	assert.True(t, pendingOrder.IsPending())
	assert.False(t, pendingOrder.IsPaid())
	assert.False(t, pendingOrder.IsConfirmed())

	paidOrder := Order{Status: OrderStatusPaid}
	assert.False(t, paidOrder.IsCreated())
	assert.False(t, paidOrder.IsPending())
	assert.True(t, paidOrder.IsPaid())
	assert.False(t, paidOrder.IsConfirmed())

	confirmedOrder := Order{Status: OrderStatusConfirmed}
	assert.False(t, confirmedOrder.IsCreated())
	assert.False(t, confirmedOrder.IsPending())
	assert.False(t, confirmedOrder.IsPaid())
	assert.True(t, confirmedOrder.IsConfirmed())

	// Test SetStatus method
	testOrder := Order{Status: OrderStatusCreated}
	testOrder.SetStatus(OrderStatusPaid)
	assert.Equal(t, OrderStatusPaid, testOrder.Status)

	// Test String() method
	assert.Equal(t, "created", OrderStatusCreated.String())
	assert.Equal(t, "pending", OrderStatusPending.String())
	assert.Equal(t, "paid", OrderStatusPaid.String())
	assert.Equal(t, "confirmed", OrderStatusConfirmed.String())
	assert.Equal(t, "unknown", OrderStatus(999).String())

	// Test integer values
	assert.Equal(t, 0, int(OrderStatusCreated))
	assert.Equal(t, 1, int(OrderStatusPending))
	assert.Equal(t, 2, int(OrderStatusPaid))
	assert.Equal(t, 3, int(OrderStatusConfirmed))

	// Test unique constraint on order_no
	duplicateOrder := Order{
		UserID:       user.ID,
		ExperienceID: experience.ID,
		Price:        500,
		Status:       OrderStatusCreated,
		OrderNo:      "ORD20250624001", // Same order_no as the first order
		OutOrderNo:   "OUT20250624002",
	}
	result = db.Create(&duplicateOrder)
	assert.Error(t, result.Error) // Should fail due to unique constraint

	// Test order_no and out_order_no fields
	uniqueOrder := Order{
		UserID:       user.ID,
		ExperienceID: experience.ID,
		Price:        750,
		Status:       OrderStatusCreated,
		OrderNo:      "ORD20250624006", // Unique order_no
		OutOrderNo:   "OUT20250624006",
	}
	result = db.Create(&uniqueOrder)
	assert.NoError(t, result.Error)
	assert.Equal(t, "ORD20250624006", uniqueOrder.OrderNo)
	assert.Equal(t, "OUT20250624006", uniqueOrder.OutOrderNo)
}
