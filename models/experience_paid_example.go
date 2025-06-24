package models

import (
	"fmt"

	"gorm.io/gorm"
)

// ExampleExperiencePaidUsage demonstrates how the virtual paid property works
func ExampleExperiencePaidUsage(db *gorm.DB) {
	// Create an experience
	experience := Experience{
		TopicID: 1,
		UserID:  1,
	}
	db.Create(&experience)

	// Initially, experience is not paid (no order)
	var exp1 Experience
	db.Preload("Order").First(&exp1, experience.ID)
	fmt.Printf("Experience paid status (no order): %t\n", exp1.Paid()) // false

	// Create an order with "created" status
	order := Order{
		UserID:       1,
		ExperienceID: experience.ID,
		Price:        1000,
		Status:       OrderStatusCreated,
		OrderNo:      "ORD001",
		OutOrderNo:   "PAY001",
	}
	db.Create(&order)

	// Experience is still not paid (order status is "created")
	var exp2 Experience
	db.Preload("Order").First(&exp2, experience.ID)
	fmt.Printf("Experience paid status (created order): %t\n", exp2.Paid()) // false

	// Update order status to "paid"
	order.Status = OrderStatusPaid
	db.Save(&order)

	// Now experience is paid
	var exp3 Experience
	db.Preload("Order").First(&exp3, experience.ID)
	fmt.Printf("Experience paid status (paid order): %t\n", exp3.Paid()) // true

	// Update order status to "confirmed"
	order.Status = OrderStatusConfirmed
	db.Save(&order)

	// Experience is still paid (order status is "confirmed")
	var exp4 Experience
	db.Preload("Order").First(&exp4, experience.ID)
	fmt.Printf("Experience paid status (confirmed order): %t\n", exp4.Paid()) // true

	// The virtual property works in JSON serialization too
	if exp4.Order != nil {
		fmt.Printf("Experience has order: %s with status %s\n", exp4.Order.OrderNo, exp4.Order.Status.String())
	} else {
		fmt.Printf("Experience has no order\n")
	}
}
