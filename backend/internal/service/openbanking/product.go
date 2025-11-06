// internal/service/openbanking/product.go
package openbanking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"multibank/backend/internal/domain"
	httputils "multibank/backend/internal/http-server/utils"
	"net/http"
	"net/url"
	"time"
)

type ProductsClient struct {
	log  log.Logger
	HTTP *http.Client
}

type bankProduct struct {
	ProductID    string  `json:"productId"`
	ProductType  string  `json:"productType"`
	ProductName  string  `json:"productName"`
	Description  string  `json:"description"`
	InterestRate float64 `json:"interestRate"`
	MinAmount    float64 `json:"minAmount"`
	MaxAmount    float64 `json:"maxAmount"`
	TermMonths   int     `json:"termMonths"`
}

func (c *ProductsClient) GetProducts(ctx context.Context, apiBaseURL, bearer string, productType string) ([]bankProduct, error) {

	const op = "service.openbanking.GetProducts"

	base, err := httputils.NormalizeURL(apiBaseURL)
	if err != nil {
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
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("products %d", resp.StatusCode)
	}

	var out []bankProduct
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// Map в доменную модель уже агрегатор сделает.
func (bp bankProduct) ToDomain(bid int64, bcode, bname string) domain.Product {
	return domain.Product{
		ProductID:    bp.ProductID,
		ProductType:  domain.ProductType(bp.ProductType),
		ProductName:  bp.ProductName,
		Description:  bp.Description,
		InterestRate: bp.InterestRate,
		MinAmount:    bp.MinAmount,
		MaxAmount:    bp.MaxAmount,
		TermMonths:   bp.TermMonths,
		BankID:       bid, BankCode: bcode, BankName: bname,
		FetchedAt: time.Now().UTC(),
	}
}
