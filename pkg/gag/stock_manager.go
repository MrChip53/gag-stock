package gag

import (
	"context"
	"log"
	"time"
)

type StockUpdateCallback func(sm *StockManager)

type StockManager struct {
	cache *Cache[GagStockResponse]

	lastStockId int64

	lastSeedRestockTime     int64
	lastGearRestockTime     int64
	lastEggRestockTime      int64
	lastMerchantRestockTime int64
	lastCosmeticRestockTime int64
	lastEventRestockTime    int64
}

func NewStockManager(ctx context.Context, cb StockUpdateCallback) *StockManager {
	sm := &StockManager{}

	// todo: this is dumb, i'll fix it later
	sm.cache = NewCache(func(stock GagStockResponse) {
		log.Println("Stock updated")
		cb(sm)
	})

	sm.refreshStock()

	sm.start(ctx)

	return sm
}

func (s *StockManager) start(ctx context.Context) {
	go s.refreshStockRoutine(ctx)
	log.Println("Stock manager started")
}

func (s *StockManager) refreshStockRoutine(ctx context.Context) {
	timer := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			err := s.refreshStock()
			if err != nil {
				log.Printf("Failed to refresh stock: %v", err)
			}
		}
	}
}

func (s *StockManager) refreshStock() error {
	stock, err := GetGagStock()
	if err != nil {
		return err
	}

	if stock.LastApiFetch == s.lastStockId {
		return nil
	}

	s.lastStockId = stock.LastApiFetch

	s.cache.Set(stock)

	return nil
}

func (s *StockManager) GetStockTimeString() string {
	return time.Unix(s.cache.Get().LastApiFetch/1000, 0).Format(time.RFC822Z)
}

func (s *StockManager) GetStock() GagStockResponse {
	return s.cache.Get()
}

func (s *StockManager) GetWantedStock(wantedItems []string) ([]string, error) {
	stock := s.cache.Get()

	foundItems, err := stock.CheckStock(wantedItems)

	if err != nil {
		return nil, err
	}

	s.lastSeedRestockTime = stock.CategoryRefreshStatus.Seeds.LastRefresh
	s.lastGearRestockTime = stock.CategoryRefreshStatus.Gears.LastRefresh
	s.lastEggRestockTime = stock.CategoryRefreshStatus.Eggs.LastRefresh
	s.lastMerchantRestockTime = stock.CategoryRefreshStatus.Merchants.LastRefresh
	s.lastCosmeticRestockTime = stock.CategoryRefreshStatus.Cosmetics.LastRefresh
	s.lastEventRestockTime = stock.CategoryRefreshStatus.Event.LastRefresh

	//notifyDesktop(foundItems, time.Unix(stock.LastApiFetch/1000, 0).Format(time.RFC822Z))

	return foundItems, nil
}
