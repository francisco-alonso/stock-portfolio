package application

import (
	"context"

	"github.com/francisco-alonso/stock-portfolio/trade-engine-service/internal/domain"
)

type Publisher interface {
	Publish(ctx context.Context, order domain.Order) error
}

type TradeService struct {
	publisher Publisher
}

func NewTradeService(publisher Publisher) *TradeService {
	return &TradeService{
		publisher: publisher,
	}
}

func (s *TradeService) CreateOrder(ctx context.Context, order domain.Order)error {
	return s.publisher.Publish(ctx, order)
}