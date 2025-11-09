// internal/service/product/recommended_service.go
package product

import (
	"context"
	"time"
)

type RecommendedService struct {
	repo RecommendedRepo
}

func NewRecommendedService(repo RecommendedRepo) *RecommendedService {
	return &RecommendedService{repo: repo}
}

type Rule struct {
	ProductID, BankCode, ProductType string
	CreatedAt                        time.Time
}

func (s *RecommendedService) List(ctx context.Context) ([]Rule, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]Rule, 0, len(rows))
	for _, r := range rows {
		t, _ := time.Parse("2006-01-02 15:04:05", r.CreatedAt)
		if t.IsZero() {
			t, _ = time.Parse("2006-01-02 15:04:05.99999", r.CreatedAt)
		}
		out = append(out, Rule{
			ProductID:   r.ProductID,
			BankCode:    r.BankCode,
			ProductType: r.ProductType,
			CreatedAt:   t,
		})
	}
	return out, nil
}

func (s *RecommendedService) Upsert(ctx context.Context, productID, bankCode, productType string) error {
	return s.repo.Upsert(ctx, productID, bankCode, productType)
}

func (s *RecommendedService) Delete(ctx context.Context, productID, bankCode, productType string) error {
	return s.repo.Delete(ctx, productID, bankCode, productType)
}
