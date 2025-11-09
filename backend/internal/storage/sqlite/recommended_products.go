package sqlite

import (
	"context"
	"database/sql"
)

type RecommendedProductsRepo struct {
	db *sql.DB
}

func NewRecommendedProductsRepo(db *sql.DB) *RecommendedProductsRepo {
	return &RecommendedProductsRepo{db: db}
}

func (r *RecommendedProductsRepo) Upsert(ctx context.Context, productID, bankCode, productType string) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO recommended_products (product_id, bank_code, product_type)
VALUES (?, ?, ?)
ON CONFLICT(product_id, bank_code, product_type) DO NOTHING
`, productID, bankCode, productType)
	return err
}

func (r *RecommendedProductsRepo) Delete(ctx context.Context, productID, bankCode, productType string) error {
	_, err := r.db.ExecContext(ctx, `
DELETE FROM recommended_products
WHERE product_id=? AND bank_code=? AND product_type=?`, productID, bankCode, productType)
	return err
}

// Snapshot returns the map for a quick product match
func (r *RecommendedProductsRepo) Snapshot(ctx context.Context) (map[string]struct{}, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT product_id, bank_code, product_type FROM recommended_products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]struct{}, 256)
	for rows.Next() {
		var pid, bcode, ptype string
		if err := rows.Scan(&pid, &bcode, &ptype); err != nil {
			return nil, err
		}
		out[makeKey(pid, bcode, ptype)] = struct{}{}
	}
	return out, rows.Err()
}

func (r *RecommendedProductsRepo) List(ctx context.Context) ([]struct {
	ProductID, BankCode, ProductType string
	CreatedAt                        string
}, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT product_id, bank_code, product_type, created_at
FROM recommended_products
ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]struct {
		ProductID, BankCode, ProductType string
		CreatedAt                        string
	}, 0, 128)

	for rows.Next() {
		var pid, bcode, ptype, created string
		if err := rows.Scan(&pid, &bcode, &ptype, &created); err != nil {
			return nil, err
		}
		out = append(out, struct {
			ProductID, BankCode, ProductType string
			CreatedAt                        string
		}{pid, bcode, ptype, created})
	}
	return out, rows.Err()
}

func makeKey(pid, bcode, ptype string) string {
	// без разделителей, чтобы не аллоцировать много строк
	return pid + "\x00" + bcode + "\x00" + ptype
}
