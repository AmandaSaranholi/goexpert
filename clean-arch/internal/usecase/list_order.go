package usecase

import (
	"github.com/AmandaSaranholi/goexpert/clean-arch/internal/entity"
	"github.com/AmandaSaranholi/goexpert/clean-arch/pkg/events"
)

type ListOrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderListed     events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewListOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderListed events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
		OrderListed:     OrderListed,
		EventDispatcher: EventDispatcher,
	}
}

func (c *ListOrderUseCase) Execute() ([]ListOrderOutputDTO, error) {
	orders, err := c.OrderRepository.List()
	if err != nil {
		return nil, err
	}

	var dtoOrders []ListOrderOutputDTO
	for _, order := range orders {
		dto := ListOrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
		dtoOrders = append(dtoOrders, dto)
	}

	c.OrderListed.SetPayload(dtoOrders)
	c.EventDispatcher.Dispatch(c.OrderListed)

	return dtoOrders, nil
}
