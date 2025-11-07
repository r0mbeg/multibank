// internal/service/openbanking/product.go
package openbanking

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"multibank/backend/internal/domain"
	httputils "multibank/backend/internal/http-server/utils"
	"multibank/backend/internal/logger"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ProductClient struct {
	log  *slog.Logger
	HTTP *http.Client
}

func NewProductClient(log *slog.Logger, httpClient *http.Client) *ProductClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &ProductClient{log: log, HTTP: httpClient}
}

type bankResponse struct {
	Data struct {
		Product []bankProduct `json:"product"`
	} `json:"data"`
}

// unified product sturcture
type bankProduct struct {
	ProductID    string  `json:"productId"`
	ProductType  string  `json:"productType"`
	ProductName  string  `json:"productName"`
	Description  *string `json:"description"`
	InterestRate *string `json:"interestRate"` // строка или null
	MinAmount    *string `json:"minAmount"`
	MaxAmount    *string `json:"maxAmount"`
	TermMonths   *int    `json:"termMonths"`
}

// helpers
func parseFloatPtr(s *string) float64 {
	if s == nil {
		return 0
	}
	str := strings.TrimSpace(*s)
	if str == "" {
		return 0
	}
	str = strings.ReplaceAll(str, ",", ".")
	f, err := strconv.ParseFloat(str, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return 0
	}
	return f
}

func strOrEmpty(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func intOrZero(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func (p bankProduct) ToDomain(bid int64, bcode, bname string) domain.Product {
	return domain.Product{
		ProductID:    p.ProductID,
		ProductType:  domain.ProductType(p.ProductType),
		ProductName:  p.ProductName,
		Description:  strOrEmpty(p.Description),
		InterestRate: parseFloatPtr(p.InterestRate),
		MinAmount:    parseFloatPtr(p.MinAmount),
		MaxAmount:    parseFloatPtr(p.MaxAmount),
		TermMonths:   intOrZero(p.TermMonths),
		BankID:       bid,
		BankCode:     bcode,
		BankName:     bname,
		FetchedAt:    time.Now().UTC(),
	}
}

func (c *ProductClient) GetProducts(ctx context.Context, apiBaseURL, bearer, productType string) ([]domain.Product, error) {
	const op = "service.openbanking.GetProducts"

	log := c.log.With(
		slog.String("op", op),
		slog.String("apiBaseURL", apiBaseURL),
		slog.String("productType", productType),
	)

	log.Info("getting products")

	base, err := httputils.NormalizeURL(apiBaseURL)
	if err != nil {
		log.Warn("error normalizing base url", logger.Err(err))
		return nil, err
	}

	u := base.ResolveReference(&url.URL{Path: "/products"})
	if productType != "" {
		q := u.Query()
		q.Set("product_type", productType)
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		log.Warn("error creating request", logger.Err(err))
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if bearer == "" {
		log.Warn("no auth header provided")
	} else {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		log.Warn("error getting products", logger.Err(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Warn("got non-200 status code",
			slog.Int("status_code", resp.StatusCode),
			slog.String("body", string(body)),
		)
		return nil, fmt.Errorf("products: http %d", resp.StatusCode)
	}

	var br bankResponse
	if err := json.NewDecoder(resp.Body).Decode(&br); err != nil {
		log.Warn("error decoding products", logger.Err(err))
		return nil, err
	}

	out := make([]domain.Product, 0, len(br.Data.Product))
	for _, p := range br.Data.Product {
		out = append(out, p.ToDomain(0, "", "")) // подставь bankID/Code/Name при необходимости
	}

	log.Info("successfully got products")

	return out, nil
}
