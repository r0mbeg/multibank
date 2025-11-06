// internal/service/product/service.go
package product

import (
	"context"
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
	client *openbanking.ProductsClient
}

func New(log *slog.Logger, repo BanksRepo, tokens *bank.Service, httpClient *openbanking.ProductsClient) *Service {
	return &Service{log: log, repo: repo, tokens: tokens, client: httpClient}
}

func (s *Service) List(ctx context.Context, f domain.ProductFilter) ([]domain.Product, error) {
	// 1) банки
	banks, err := s.repo.ListEnabledBanks(ctx)
	if err != nil {
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

	// 2) параллельные запросы в банки
	var mu sync.Mutex
	out := make([]domain.Product, 0, len(banks)*8)

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(8) // ограничим параллелизм

	for _, b := range banks {
		b := b
		eg.Go(func() error {
			// токен (если каталог публичный — токен можно опустить; но оставим универсально)
			token, _, err := s.tokens.GetOrRefreshToken(ctx, b.ID)
			if err != nil {
				// мягко пропускаем банк, но возвращаем ошибку, если хочешь — меняй политику
				s.log.Warn("cannot get token for products", slog.Int64("bank_id", b.ID), slog.String("bank", b.Code), slog.Any("err", err))
				return nil
			}

			items, err := s.client.GetProducts(ctx, b.APIBaseURL, token, f.ProductType)
			if err != nil {
				s.log.Warn("products fetch failed", slog.Int64("bank_id", b.ID), slog.String("bank", b.Code), slog.Any("err", err))
				return nil
			}
			// map
			loc := make([]domain.Product, 0, len(items))
			for _, it := range items {
				loc = append(loc, it.ToDomain(b.ID, b.Code, b.Name))
			}
			mu.Lock()
			out = append(out, loc...)
			mu.Unlock()
			return nil
		})
	}
	_ = eg.Wait() // ошибки мы уже «проглотили» и залогировали

	return out, nil
}
