# Experience Payment API

## POST /experiences/:id/paid

Mark an experience as paid by creating an associated order with "paid" status.

### Request

**URL Parameters:**
- `id` (integer, required): The ID of the experience to mark as paid

**Headers:**
- `Content-Type: application/json`
- `Authorization: Bearer <token>` (required for authentication)

**Request Body:**
```json
{
  "order_no": "ORD20250624001",
  "out_order_no": "PAY20250624001", 
  "price": 1000
}
```

**Fields:**
- `order_no` (string, required): Internal order number (must be unique)
- `out_order_no` (string, required): External payment system order number
- `price` (integer, required): Price in cents/smallest currency unit

### Response

**Success Response (200 OK):**
```json
{
  "id": 1,
  "topic_id": 1,
  "user_id": 1,
  "paid": true,
  "created_at": "2025-06-24T14:30:00Z",
  "updated_at": "2025-06-24T14:30:00Z",
  "order": {
    "id": 1,
    "order_no": "ORD20250624001",
    "out_order_no": "PAY20250624001",
    "price": 1000,
    "status": "paid",
    "created_at": "2025-06-24T14:30:00Z",
    "updated_at": "2025-06-24T14:30:00Z"
  }
}
```

**Error Responses:**

- **400 Bad Request**: Invalid experience ID or missing required fields
```json
{
  "error": "invalid experience id"
}
```

- **401 Unauthorized**: User not authenticated
```json
{
  "error": "User not authenticated"
}
```

- **403 Forbidden**: User doesn't own the experience
```json
{
  "error": "forbidden: not your experience"
}
```

- **404 Not Found**: Experience not found
```json
{
  "error": "experience not found"
}
```

- **409 Conflict**: Experience already has an order
```json
{
  "error": "experience already has an order"
}
```

### Business Logic

1. **Authentication**: User must be authenticated via JWT token
2. **Authorization**: User can only mark their own experiences as paid
3. **One-to-One Relationship**: Each experience can have at most one order
4. **Virtual Paid Property**: Experience's `paid` status is calculated from the associated order:
   - `paid: false` - No order, or order with status "created" or "pending"
   - `paid: true` - Order with status "paid" or "confirmed"
5. **Order Creation**: Creates order with status "paid" (integer value: 2)

### Example Usage

```bash
curl -X POST http://localhost:8000/experiences/1/paid \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "order_no": "ORD20250624001",
    "out_order_no": "PAY20250624001",
    "price": 1000
  }'
```

### Database Changes

The API creates:
1. **Order record** with:
   - `experience_id`: Links to the experience
   - `user_id`: Current authenticated user
   - `status`: Set to "paid" (integer value: 2)
   - `order_no`, `out_order_no`, `price`: From request body
   - Auto-generated timestamps

2. **Experience virtual property**:
   - `paid` field is now calculated dynamically
   - No longer stored in database, computed from order status

---

## PayOrder API Enhancement

The `POST /pay/order` endpoint has been enhanced to include the `out_order_no` in its response.

### PayOrder Response Enhancement

**Previous Behavior**: The PayOrder endpoint returned only the raw Douyin API response.

**New Behavior**: The PayOrder endpoint now merges the generated `out_order_no` into the response.

**Enhanced Response Example:**
```json
{
  "err_no": 0,
  "err_tips": "success",
  "order_id": "7123456789012345678",
  "order_token": "ChAKGG91dF9vcmRlcl9ub18xNjc...",
  "out_order_no": "1kz2x3c4v5b6n7m8Abc1"
}
```

**Key Changes:**
- The `out_order_no` field is now included in all PayOrder responses
- This allows clients to track the order using the generated order number
- The response maintains backward compatibility with existing Douyin API fields
- If JSON parsing fails, the original response is returned unchanged
