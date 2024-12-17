package constant

type OrderStatus string

const (
	ORDER_STATUS_CREATED    OrderStatus = "CREATED"
	ORDER_STATUS_PROCESSING OrderStatus = "PROCESSING"
	ORDER_STATUS_COMPLETED  OrderStatus = "COMPLETED"
	ORDER_STATUS_CANCELLED  OrderStatus = "CANCELLED"
	ORDER_STATUS_FAILED     OrderStatus = "FAILED"
)
