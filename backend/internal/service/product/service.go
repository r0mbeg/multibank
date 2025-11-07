// internal/service/product/service.go
package product

import (
	"context"
	"multibank/backend/internal/logger"
	"sync"

	"golang.org/x/sync/errgroup"

	"multibank/backend/internal/domain"
	"multibank/backend/internal/service/bank"
	"multibank/backend/internal/service/openbanking"

	"log/slog"
	"time"
)

type BanksRepo interface {
	ListEnabledBanks(ctx context.Context) ([]domain.Bank, error)
}

type BankTokens interface {
	GetOrRefreshToken(ctx context.Context, bankID int64) (string, time.Time, error)
}

type Service struct {
	log    *slog.Logger
	repo   BanksRepo  // чтобы получить список банков
	tokens BankTokens // чтобы достать токен
	client *openbanking.ProductClient
}

func New(log *slog.Logger, repo BanksRepo, tokens *bank.Service, httpClient *openbanking.ProductClient) *Service {
	return &Service{log: log, repo: repo, tokens: tokens, client: httpClient}
}

func (s *Service) List(ctx context.Context, f domain.ProductFilter) ([]domain.Product, error) {

	const op = "service.product.List"

	log := s.log.With(slog.String("op", op))

	// 1) банки
	banks, err := s.repo.ListEnabledBanks(ctx)
	if err != nil {
		log.Warn("failed to list enabled banks")
		return nil, err
	}

	// если фильтр по банкам — отфильтруем
	if len(f.BankIDs) > 0 {
		set := make(map[int64]struct{}, len(f.BankIDs))
		for _, id := range f.BankIDs {
			set[id] = struct{}{}
		}
		tmp := banks[:0]
		for _, b := range banks {
			if _, ok := set[b.ID]; ok {
				tmp = append(tmp, b)
			}
		}
		banks = tmp
	}

	// 2) parallel requests to banks
	var mu sync.Mutex
	out := make([]domain.Product, 0, len(banks)*8)

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(8) // limit parallelism (sth like worker pool)

	for _, b := range banks {
		b := b
		eg.Go(func() error {
			// token (if the catalog is public, token can be omited)
			token, _, err := s.tokens.GetOrRefreshToken(ctx, b.ID)
			if err != nil {
				// мягко пропускаем банк, но возвращаем ошибку
				log.Warn("cannot get token for products",
					slog.Int64("bank_id", b.ID),
					slog.String("bank", b.Code),
					logger.Err(err),
				)
				return nil
			}

			items, err := s.client.GetProducts(ctx, b.APIBaseURL, token, f.ProductType)
			if err != nil {
				log.Warn("products fetch failed",
					slog.Int64("bank_id", b.ID),
					slog.String("bank", b.Code),
					logger.Err(err),
				)
				return nil
			}
			// map
			loc := make([]domain.Product, 0, len(items))
			for _, it := range items {

				it.BankID = b.ID
				it.BankCode = b.Code
				it.BankName = b.Name
				
				loc = append(loc, it)
			}
			mu.Lock()
			out = append(out, loc...)
			mu.Unlock()
			return nil
		})
	}
	_ = eg.Wait() // errors are already passed and logged

	return out, nil
}
