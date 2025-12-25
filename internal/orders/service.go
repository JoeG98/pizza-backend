package orders

import (
	"encoding/json"
	"errors"

	"github.com/JoeG98/pizza-backend/internal/database"
	"github.com/JoeG98/pizza-backend/internal/models"
	"github.com/JoeG98/pizza-backend/internal/sse"
)

type Service struct {
	db  *database.Database
	hub *sse.Hub
}

func OrderService(db *database.Database, hub *sse.Hub) *Service {
	return &Service{
		db:  db,
		hub: hub,
	}
}

func (s *Service) CreateOrder(input CreateOrderRequest) (*models.Order, error) {
	if input.CustomerName == "" || input.Phone == "" || input.Address == "" {
		return nil, errors.New("Missing required customer fields")
	}

	if len(input.Pizzas) == 0 {
		return nil, errors.New("Order must include atleast one pizza.")
	}

	for _, p := range input.Pizzas {
		// validate size
		validSize := false

		for _, size := range models.PizzaSizes {
			if p.Size == size {
				validSize = true
				break
			}
		}

		if !validSize {
			return nil, errors.New("Invalid pizza size: " + p.Size)
		}

		validType := false

		for _, t := range models.PizzaTypes {
			if p.Pizza == t {
				validType = true
				break
			}
		}

		if !validType {
			return nil, errors.New("Invalid pizza type: " + p.Pizza)
		}
	}

	// map request to DB model
	order := models.Order{
		Status:       models.OrderStatuses[0],
		CustomerName: input.CustomerName,
		Phone:        input.Phone,
		Address:      input.Address,
	}

	for _, p := range input.Pizzas {
		order.Items = append(order.Items, models.OrderItem{
			Size:         p.Size,
			Pizza:        p.Pizza,
			Instructions: p.Instructions,
		})
	}

	err := s.db.DB.Create(&order).Error

	if err != nil {
		return nil, err
	}

	event := sse.Event{
		Type: "order_created",
		Data: order,
	}

	payload, err := json.Marshal(event)

	if err == nil {
		s.hub.Broadcast <- string(payload)
	}

	return &order, nil
}

func (s *Service) GetOrder(id string) (*models.Order, error) {
	var order models.Order

	// preload pizza items

	err := s.db.DB.Preload("Items").First(&order, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *Service) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order

	err := s.db.DB.Preload("Items").Order("created_at desc").Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Service) UpdateOrderStatus(id string, status string) (*models.Order, error) {
	// Validate against allowed statues

	valid := false

	for _, statusOption := range models.OrderStatuses {
		if statusOption == status {
			valid = true
			break
		}
	}

	if !valid {
		return nil, errors.New("Invalid status value")
	}

	// Get Order First

	order, err := s.GetOrder(id)

	if err != nil {
		return nil, err
	}

	// Update Status

	order.Status = status

	if err := s.db.DB.Save(order).Error; err != nil {
		return nil, err
	}

	event := sse.Event{
		Type: "order_status_updated",
		Data: order,
	}

	payload, err := json.Marshal(event)

	if err == nil {
		s.hub.Broadcast <- string(payload)
	}

	return order, nil
}

func (s *Service) DeleteOrder(id string) error {
	// check if exists

	_, err := s.GetOrder(id)

	if err != nil {
		return err
	}

	if err := s.db.DB.Select("Items").Delete(&models.Order{ID: id}).Error; err != nil {
		return err
	}

	return nil
}
