package services

import (
	"fmt"

	"github.com/francisco-alonso/stock-portfolio/portfolio-service/domain"
	"github.com/francisco-alonso/stock-portfolio/portfolio-service/infra"
)

type PortfolioService interface {
	AddPosition(userID, asset string, quantity int, price float64) error
	GetPositions(userID string) ([]domain.Position, error)
}

type PortfolioServiceImpl struct {
	repo infra.PortfolioRepository
}

func NewPortfolioService(repo infra.PortfolioRepository) PortfolioService{
	return &PortfolioServiceImpl{repo: repo}
}

func (s *PortfolioServiceImpl) AddPosition(userID, asset string, quantity int, price float64) error {
    if quantity <= 0 || price < 0 {
        return fmt.Errorf("invalid quantity or price")
    }
    return s.repo.AddPosition(userID, asset, quantity, price)
}

func (s *PortfolioServiceImpl) GetPositions(userID string) ([]domain.Position, error) {
    if userID == "" {
        return nil, fmt.Errorf("invalid user id - empty value")
    }
    return s.repo.GetPositions(userID)
}