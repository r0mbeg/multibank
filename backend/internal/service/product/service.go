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

type RecommendedRepo interface {
	Snapshot(ctx context.Context) (map[string]struct{}, error)
	List(ctx context.Context) ([]struct {
		ProductID, BankCode, ProductType string
		CreatedAt                        string
	}, error)
	Upsert(ctx context.Context, productID, bankCode, productType string) error
	Delete(ctx context.Context, productID, bankCode, productType string) error
}

type Service struct {
	log         *slog.Logger
	repo        BanksRepo       // to get a list of banks
	recommended RecommendedRepo // to set IsRecommended
	tokens      BankTokens      // to get the bank token
	client      *openbanking.ProductClient
}

func New(
	log *slog.Logger,
	repo BanksRepo,
	tokens *bank.Service,
	recommended RecommendedRepo,
	httpClient *openbanking.ProductClient,
) *Service {
	return &Service{
		log:         log,
		repo:        repo,
		tokens:      tokens,
		recommended: recommended,
		client:      httpClient,
	}
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

	eg, egCtx := errgroup.WithContext(ctx) // <-- НЕ затеняем исходный ctx
	eg.SetLimit(8)

	for _, b := range banks {
		b := b
		eg.Go(func() error {
			token, _, err := s.tokens.GetOrRefreshToken(egCtx, b.ID)
			if err != nil {
				log.Warn("cannot get token for products",
					slog.Int64("bank_id", b.ID),
					slog.String("bank", b.Code),
					logger.Err(err),
				)
				return nil // мягкий пропуск
			}

			items, err := s.client.GetProducts(egCtx, b.APIBaseURL, token, f.ProductType)
			if err != nil {
				log.Warn("products fetch failed",
					slog.Int64("bank_id", b.ID),
					slog.String("bank", b.Code),
					logger.Err(err),
				)
				return nil // мягкий пропуск
			}

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

	// wait for all requests
	if err := eg.Wait(); err != nil {
		log.Warn("product fetch group finished with error", logger.Err(err))
	}

	// 3) set IsRecommended from snapshot
	snapCtx := context.Background()

	set, err := s.recommended.Snapshot(snapCtx)
	if err != nil {
		log.Warn("recommended snapshot failed", logger.Err(err))
		set = map[string]struct{}{}
	} else {
		log.Info("got recommend map", slog.Int("count", len(set)))
	}

	for i := range out {
		key := recKey(out[i].ProductID, out[i].BankCode, out[i].ProductType)
		if _, ok := set[key]; ok {
			out[i].IsRecommended = true
			// debug
			log.Info("product is in recommended list",
				slog.String("productId", out[i].ProductID),
				slog.String("productType", out[i].ProductType),
				slog.String("bank_code", out[i].BankCode),
			)
		}
	}

	return out, nil
}

// key compatible with the repository
func recKey(pid, bankCode, ptype string) string {
	return pid + "\x00" + bankCode + "\x00" + ptype
}
