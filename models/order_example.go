package models

import (
	"encoding/json"
	"fmt"
)

// Example function to demonstrate the integer-based enum usage
func ExampleOrderStatusUsage() {
	// Create an order with integer-based status
	order := Order{
		UserID:       1,
		ExperienceID: 1,
		Price:        1000,
		Status:       OrderStatusCreated, // This is stored as integer 0 in database
		OrderNo:      "ORD20250624001",   // Internal order number (unique)
		OutOrderNo:   "PAY20250624001",   // External payment system order number
	}

	// Print the status as string for human readability
	fmt.Printf("Order status: %s\n", order.Status.String()) // Output: "created"

	// Check status using helper methods
	if order.IsCreated() {
		fmt.Println("Order is in created state")
	}

	// Update status
	order.SetStatus(OrderStatusPaid)
	fmt.Printf("Updated status: %s\n", order.Status.String()) // Output: "paid"

	// JSON marshaling will show the string representation
	jsonData, _ := json.Marshal(order)
	fmt.Printf("JSON representation: %s\n", string(jsonData))
	// The status field in JSON will be "paid" (string), but stored as 2 (int) in database

	// Database operations are more efficient with integers
	fmt.Printf("Database value: %d\n", int(order.Status)) // Output: 2
}
