package gag

import (
	"context"
	"log"
	"sync"
	"time"
)

type StockUpdateCallback func(sc ShopContainer)

type ShopStock map[string]int

type Shop struct {
	Items       ShopStock
	LastRefresh int64
	RestockTime int64
}

func (s *Shop) GetItemCount(item string) int {
	count, ok := s.Items[item]
	if !ok {
		return 0
	}

	return count
}

func (s *Shop) GetAllItems() []Item {
	items := []Item{}

	for item, count := range s.Items {
		items = append(items, Item{
			Name:        item,
			Count:       count,
			RestockTime: int(s.RestockTime),
		})
	}

	return items
}

type StockManager struct {
	mu sync.RWMutex

	shopContainer ShopContainer
	cb            StockUpdateCallback

	lastApiFetch int64

	lastSeedRestockTime     int64
	lastGearRestockTime     int64
	lastEggRestockTime      int64
	lastMerchantRestockTime int64
	lastCosmeticRestockTime int64
	lastEventRestockTime    int64

	imageData map[string]string
}

func NewStockManager(ctx context.Context, cb StockUpdateCallback) *StockManager {
	sm := &StockManager{
		shopContainer: NewShopContainer(),
		cb:            cb,
	}

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

	if stock.LastApiFetch == s.lastApiFetch {
		return nil
	}

	s.lastApiFetch = stock.LastApiFetch

	s.lastSeedRestockTime = stock.CategoryRefreshStatus.Seeds.LastRefresh
	s.lastGearRestockTime = stock.CategoryRefreshStatus.Gears.LastRefresh
	s.lastEggRestockTime = stock.CategoryRefreshStatus.Eggs.LastRefresh
	s.lastMerchantRestockTime = stock.CategoryRefreshStatus.Merchants.LastRefresh
	s.lastCosmeticRestockTime = stock.CategoryRefreshStatus.Cosmetics.LastRefresh
	s.lastEventRestockTime = stock.CategoryRefreshStatus.Event.LastRefresh

	s.imageData = stock.ImageData

	sc := stock.ToShopContainer()

	s.mu.Lock()
	defer s.mu.Unlock()

	for shop, stock := range sc.container {
		s.shopContainer.container[shop] = stock
	}

	s.cb(sc)

	return nil
}

func (s *StockManager) GetStockTimeString() string {
	return time.Unix(s.lastApiFetch/1000, 0).Format(time.RFC822Z)
}

func (s *StockManager) GetShopStock(shop string) Shop {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.shopContainer.container[shop]
}

func (s *StockManager) GetShopContainer() ShopContainer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sc := NewShopContainer()
	sc.lastApiFetch = s.shopContainer.lastApiFetch

	for shop, stock := range s.shopContainer.container {
		sc.container[shop] = stock
	}

	return sc
}

func (s *StockManager) GetImageData() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.imageData
}
