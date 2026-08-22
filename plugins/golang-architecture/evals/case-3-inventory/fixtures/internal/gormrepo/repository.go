package gormrepo

import "gorm.io/gorm"

// StockRepository persists and retrieves stock levels.
type StockRepository interface {
	AdjustQuantity(sku string, delta int) error
	GetQuantity(sku string) (int, error)
}

type Repository struct {
	db *gorm.DB
}

func Connect(dsn string) *gorm.DB {
	db, err := gorm.Open(nil, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AdjustQuantity(sku string, delta int) error {
	return r.db.Exec("UPDATE stock SET quantity = quantity + ? WHERE sku = ?", delta, sku).Error
}

func (r *Repository) GetQuantity(sku string) (int, error) {
	var qty int
	err := r.db.Raw("SELECT quantity FROM stock WHERE sku = ?", sku).Scan(&qty).Error
	return qty, err
}
