package gag

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type LastRestockTimes struct {
	Seeds     int64 `json:"seeds"`
	Gears     int64 `json:"gears"`
	Eggs      int64 `json:"eggs"`
	Merchants int64 `json:"merchants"`
	Cosmetics int64 `json:"cosmetics"`
	Event     int64 `json:"event"`
}

type StockItemList []StockItem

func (s StockItemList) GetItem(name string) *StockItem {
	name = strings.ToLower(name)
	for _, item := range s {
		if strings.ToLower(item.Name) == name {
			return &item
		}
	}
	return nil
}

// StockItem represents an item with a name and value
type StockItem struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type LastSeenItem struct {
	Name string  `json:"name"`
	Seen *string `json:"seen"`
}

type RestockTimers struct {
	Seeds     int64 `json:"seeds"`
	Gears     int64 `json:"gears"`
	Eggs      int64 `json:"eggs"`
	Merchants int64 `json:"merchants"`
	Cosmetics int64 `json:"cosmetics"`
	Event     int64 `json:"event"`
}

type CategoryRefresh struct {
	WasRefreshed     bool  `json:"wasRefreshed"`
	LastRefresh      int64 `json:"lastRefresh"`
	TimeSinceRefresh int64 `json:"timeSinceRefresh"`
	ExpectingUpdate  bool  `json:"expectingUpdate"`
}

type CategoryRefreshStatus struct {
	Seeds     CategoryRefresh `json:"seeds"`
	Gears     CategoryRefresh `json:"gears"`
	Eggs      CategoryRefresh `json:"eggs"`
	Merchants CategoryRefresh `json:"merchants"`
	Cosmetics CategoryRefresh `json:"cosmetics"`
	Event     CategoryRefresh `json:"event"`
}

type GagStockResponse struct {
	EasterStock           StockItemList         `json:"easterStock"`
	GearStock             StockItemList         `json:"gearStock"`
	EggStock              StockItemList         `json:"eggStock"`
	NightStock            StockItemList         `json:"nightStock"`
	EventStock            StockItemList         `json:"eventStock"`
	CosmeticsStock        StockItemList         `json:"cosmeticsStock"`
	SeedsStock            StockItemList         `json:"seedsStock"`
	MerchantsStock        StockItemList         `json:"merchantsStock"`
	LastSeen              []LastSeenItem        `json:"lastSeen"`
	RestockTimers         RestockTimers         `json:"restockTimers"`
	CategoryRefreshStatus CategoryRefreshStatus `json:"categoryRefreshStatus"`
	TimerCalculatedAt     int64                 `json:"timerCalculatedAt"`
	ServerStartTime       int64                 `json:"serverStartTime"`
	LastApiFetch          int64                 `json:"lastApiFetch"`
	NextScheduledFetch    int64                 `json:"nextScheduledFetch"`
	ImageData             map[string]string     `json:"imageData"`
}

func (g *GagStockResponse) GetLastRestockTimes() LastRestockTimes {
	return LastRestockTimes{
		Seeds:     g.CategoryRefreshStatus.Seeds.LastRefresh,
		Gears:     g.CategoryRefreshStatus.Gears.LastRefresh,
		Eggs:      g.CategoryRefreshStatus.Eggs.LastRefresh,
		Merchants: g.CategoryRefreshStatus.Merchants.LastRefresh,
		Cosmetics: g.CategoryRefreshStatus.Cosmetics.LastRefresh,
		Event:     g.CategoryRefreshStatus.Event.LastRefresh,
	}
}

func (g *GagStockResponse) CheckStock(wantedItems []string, lastRestockTimes LastRestockTimes, force bool) ([]string, error) {
	foundItems := []string{}

	times := g.GetLastRestockTimes()

	for _, item := range wantedItems {
		stockItem := g.SeedsStock.GetItem(item)
		if stockItem != nil && stockItem.Value > 0 && (lastRestockTimes.Seeds < times.Seeds || force) {
			foundItems = append(foundItems, item)
			continue
		}

		stockItem = g.GearStock.GetItem(item)
		if stockItem != nil && stockItem.Value > 0 && (lastRestockTimes.Gears < times.Gears || force) {
			foundItems = append(foundItems, item)
			continue
		}

		stockItem = g.EggStock.GetItem(item)
		if stockItem != nil && stockItem.Value > 0 && (lastRestockTimes.Eggs < times.Eggs || force) {
			foundItems = append(foundItems, item)
			continue
		}

		stockItem = g.MerchantsStock.GetItem(item)
		if stockItem != nil && stockItem.Value > 0 && (lastRestockTimes.Merchants < times.Merchants || force) {
			foundItems = append(foundItems, item)
			continue
		}

		stockItem = g.CosmeticsStock.GetItem(item)
		if stockItem != nil && stockItem.Value > 0 && (lastRestockTimes.Cosmetics < times.Cosmetics || force) {
			foundItems = append(foundItems, item)
			continue
		}

		stockItem = g.EventStock.GetItem(item)
		if stockItem != nil && stockItem.Value > 0 && (lastRestockTimes.Event < times.Event || force) {
			foundItems = append(foundItems, item)
			continue
		}

		stockItem = g.NightStock.GetItem(item)
		if stockItem != nil && stockItem.Value > 0 {
			foundItems = append(foundItems, item)
			continue
		}

		stockItem = g.EasterStock.GetItem(item)
		if stockItem != nil && stockItem.Value > 0 {
			foundItems = append(foundItems, item)
			continue
		}
	}

	return foundItems, nil
}

// GetGagStock fetches the stock data from the Grow a Garden API and decodes it into GagStockResponse.
func GetGagStock() (*GagStockResponse, error) {
	const url = "https://growagarden.gg/api/stock"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var stock GagStockResponse
	if err := json.NewDecoder(resp.Body).Decode(&stock); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &stock, nil
}
