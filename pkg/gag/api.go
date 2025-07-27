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

func (s StockItemList) ToShopStock() ShopStock {
	shopStock := make(ShopStock)
	for _, item := range s {
		shopStock[item.Name] = item.Value
	}
	return shopStock
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
	Easter                StockItemList         `json:"easterStock"`
	Gear                  StockItemList         `json:"gearStock"`
	Egg                   StockItemList         `json:"eggStock"`
	Night                 StockItemList         `json:"nightStock"`
	Event                 StockItemList         `json:"eventStock"`
	Cosmetics             StockItemList         `json:"cosmeticsStock"`
	Seeds                 StockItemList         `json:"seedsStock"`
	Merchants             StockItemList         `json:"merchantsStock"`
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

func roundToNearest5Minutes(ts int64) int64 {
	const intervalMs = 5 * 60 * 1000 // 5 minutes in milliseconds
	return ((ts + intervalMs/2) / intervalMs) * intervalMs
}

func (g GagStockResponse) ToShopContainer() ShopContainer {
	sc := NewShopContainer()

	sc.lastSeen = g.LastSeen

	sc.container["seeds"] = Shop{
		Items:       g.Seeds.ToShopStock(),
		LastRefresh: g.CategoryRefreshStatus.Seeds.LastRefresh,
		RestockTime: roundToNearest5Minutes(g.TimerCalculatedAt + g.RestockTimers.Seeds),
	}

	sc.container["gears"] = Shop{
		Items:       g.Gear.ToShopStock(),
		LastRefresh: g.CategoryRefreshStatus.Gears.LastRefresh,
		RestockTime: roundToNearest5Minutes(g.TimerCalculatedAt + g.RestockTimers.Gears),
	}

	sc.container["eggs"] = Shop{
		Items:       g.Egg.ToShopStock(),
		LastRefresh: g.CategoryRefreshStatus.Eggs.LastRefresh,
		RestockTime: roundToNearest5Minutes(g.TimerCalculatedAt + g.RestockTimers.Eggs),
	}

	sc.container["merchants"] = Shop{
		Items:       g.Merchants.ToShopStock(),
		LastRefresh: g.CategoryRefreshStatus.Merchants.LastRefresh,
		RestockTime: roundToNearest5Minutes(g.TimerCalculatedAt + g.RestockTimers.Merchants),
	}

	sc.container["cosmetics"] = Shop{
		Items:       g.Cosmetics.ToShopStock(),
		LastRefresh: g.CategoryRefreshStatus.Cosmetics.LastRefresh,
		RestockTime: roundToNearest5Minutes(g.TimerCalculatedAt + g.RestockTimers.Cosmetics),
	}

	sc.container["event"] = Shop{
		Items:       g.Event.ToShopStock(),
		LastRefresh: g.CategoryRefreshStatus.Event.LastRefresh,
		RestockTime: roundToNearest5Minutes(g.TimerCalculatedAt + g.RestockTimers.Event),
	}

	sc.container["night"] = Shop{
		Items: g.Night.ToShopStock(),
	}

	sc.container["easter"] = Shop{
		Items: g.Easter.ToShopStock(),
	}

	sc.lastApiFetch = g.LastApiFetch

	return sc
}

func (g GagStockResponse) GetAllStock() []StockItem {
	items := []StockItem{}

	items = append(items, g.Seeds...)
	items = append(items, g.Gear...)
	items = append(items, g.Egg...)
	items = append(items, g.Merchants...)
	items = append(items, g.Cosmetics...)
	items = append(items, g.Event...)
	items = append(items, g.Night...)
	items = append(items, g.Easter...)

	return items
}

func (g GagStockResponse) CheckStock(wantedItems []string) ([]string, error) {
	foundItems := []string{}

	for _, item := range wantedItems {
		stockItem := g.Seeds.GetItem(item)
		if stockItem != nil && stockItem.Value > 0 {
			foundItems = append(foundItems, item)
			continue
		}

		stockItem = g.Gear.GetItem(item)
		if stockItem != nil && stockItem.Value > 0 {
			foundItems = append(foundItems, item)
			continue
		}

		stockItem = g.Egg.GetItem(item)
		if stockItem != nil && stockItem.Value > 0 {
			foundItems = append(foundItems, item)
			continue
		}

		stockItem = g.Merchants.GetItem(item)
		if stockItem != nil && stockItem.Value > 0 {
			foundItems = append(foundItems, item)
			continue
		}

		stockItem = g.Cosmetics.GetItem(item)
		if stockItem != nil && stockItem.Value > 0 {
			foundItems = append(foundItems, item)
			continue
		}

		stockItem = g.Event.GetItem(item)
		if stockItem != nil && stockItem.Value > 0 {
			foundItems = append(foundItems, item)
			continue
		}

		stockItem = g.Night.GetItem(item)
		if stockItem != nil && stockItem.Value > 0 {
			foundItems = append(foundItems, item)
			continue
		}

		stockItem = g.Easter.GetItem(item)
		if stockItem != nil && stockItem.Value > 0 {
			foundItems = append(foundItems, item)
			continue
		}
	}

	return foundItems, nil
}

// GetGagStock fetches the stock data from the Grow a Garden API and decodes it into GagStockResponse.
func GetGagStock() (GagStockResponse, error) {
	const url = "https://growagarden.gg/api/stock"
	resp, err := http.Get(url)
	if err != nil {
		return GagStockResponse{}, fmt.Errorf("failed to fetch stock: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GagStockResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var stock GagStockResponse
	if err := json.NewDecoder(resp.Body).Decode(&stock); err != nil {
		return GagStockResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return stock, nil
}
