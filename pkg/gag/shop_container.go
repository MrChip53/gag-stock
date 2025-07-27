package gag

import (
	"time"
)

type ShopContainer struct {
	container    map[string]Shop
	lastSeen     []LastSeenItem
	lastApiFetch int64
}

func NewShopContainer() ShopContainer {
	return ShopContainer{
		container: make(map[string]Shop),
		lastSeen:  make([]LastSeenItem, 0),
	}
}

func (sc ShopContainer) GetWantedStock(wantedItems []string) []Item {
	foundItems := []Item{}

	for _, item := range wantedItems {
		for _, shop := range sc.container {
			count := shop.GetItemCount(item)
			if count > 0 {
				foundItems = append(foundItems, Item{
					Name:        item,
					Count:       count,
					StockTime:   int(shop.LastRefresh),
					RestockTime: int(shop.RestockTime),
				})
			}
		}
	}

	return foundItems
}

func (sc ShopContainer) GetAllItems() []Item {
	foundItems := []Item{}

	for _, shop := range sc.container {
		foundItems = append(foundItems, shop.GetAllItems()...)
	}

	return foundItems
}

func (sc ShopContainer) GetLastSeenItems() []LastSeenItem {
	return sc.lastSeen
}

func (sc ShopContainer) GetTimeString() string {
	return time.Unix(sc.lastApiFetch/1000, 0).Format(time.RFC822Z)
}
