package orders

type CreateOrderPizza struct {
	Pizza        string `json:"pizza"`
	Size         string `json:"size"`
	Instructions string `json:"instructions"`
}

type CreateOrderRequest struct {
	CustomerName string             `json:"customerName"`
	Phone        string             `json:"phone"`
	Address      string             `json:"address"`
	Pizzas       []CreateOrderPizza `json:"pizzas"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}
