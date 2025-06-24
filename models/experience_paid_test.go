package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExperiencePaidVirtualProperty(t *testing.T) {
	db := InitTestDB()

	// Add models to test database migration
	db.AutoMigrate(&Order{}, &Experience{}, &Topic{}, &User{})

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

	// Test 1: Experience without order should not be paid
	var exp1 Experience
	db.Preload("Order").First(&exp1, experience.ID)
	assert.False(t, exp1.IsPaid())
	assert.False(t, exp1.Paid())

	// Test 2: Experience with created order should not be paid
	order := Order{
		UserID:       user.ID,
		ExperienceID: experience.ID,
		Price:        1000,
		Status:       OrderStatusCreated,
		OrderNo:      "ORD20250624001",
		OutOrderNo:   "OUT20250624001",
	}
	db.Create(&order)

	var exp2 Experience
	db.Preload("Order").First(&exp2, experience.ID)
	assert.False(t, exp2.IsPaid())
	assert.False(t, exp2.Paid())

	// Test 3: Experience with pending order should not be paid
	order.Status = OrderStatusPending
	db.Save(&order)

	var exp3 Experience
	db.Preload("Order").First(&exp3, experience.ID)
	assert.False(t, exp3.IsPaid())
	assert.False(t, exp3.Paid())

	// Test 4: Experience with paid order should be paid
	order.Status = OrderStatusPaid
	db.Save(&order)

	var exp4 Experience
	db.Preload("Order").First(&exp4, experience.ID)
	assert.True(t, exp4.IsPaid())
	assert.True(t, exp4.Paid())

	// Test 5: Experience with confirmed order should be paid
	order.Status = OrderStatusConfirmed
	db.Save(&order)

	var exp5 Experience
	db.Preload("Order").First(&exp5, experience.ID)
	assert.True(t, exp5.IsPaid())
	assert.True(t, exp5.Paid())

	// Test 6: Create another experience with confirmed order
	experience2 := Experience{
		TopicID: topic.ID,
		UserID:  user.ID,
	}
	db.Create(&experience2)

	confirmedOrder := Order{
		UserID:       user.ID,
		ExperienceID: experience2.ID,
		Price:        1200,
		Status:       OrderStatusConfirmed,
		OrderNo:      "ORD20250624002",
		OutOrderNo:   "OUT20250624002",
	}
	db.Create(&confirmedOrder)

	var exp6 Experience
	db.Preload("Order").First(&exp6, experience2.ID)
	assert.True(t, exp6.IsPaid())
	assert.True(t, exp6.Paid())
}
